package plan

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type getCmd struct {
	cl           *clientset.Clientset
	lookupByUUID bool
}

// NewGetCmd builds a "svc-cat get plans" command
func NewGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := &getCmd{cl: cl}
	cmd := &cobra.Command{
		Use:     "plans [name]",
		Aliases: []string{"plan", "pl"},
		Short:   "List plans, optionally filtered by name",
		Example: `
  svc-cat get plans
  svc-cat get plan standard800
  svc-cat get plan --uuid 08e4b43a-36bc-447e-a81f-8202b13e339c
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmd.run(args)
		},
	}
	cmd.Flags().BoolVarP(
		&getCmd.lookupByUUID,
		"uuid",
		"u",
		false,
		"Whether or not to get the plan by UUID (the default is by name)",
	)
	return cmd
}

func (c *getCmd) run(args []string) error {
	if len(args) == 0 {
		return c.getAll()
	} else {
		key := args[0]
		return c.get(key)
	}
}

func (c *getCmd) getAll() error {
	plans, err := c.cl.ServicecatalogV1beta1().ClusterServicePlans().List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("unable to list plans (%s)", err)
	}

	// Retrieve the classes as well because plans don't have the external class name
	classes, err := c.cl.ServicecatalogV1beta1().ClusterServiceClasses().List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("unable to list classes (%s)", err)
	}

	output.WritePlanList(plans.Items, classes.Items)
	return nil
}

func (c *getCmd) get(key string) error {
	var plan *v1beta1.ClusterServicePlan
	var err error
	if c.lookupByUUID {
		plan, err = retrieveByUUID(c.cl, key)
	} else {
		plan, err = retrieveByName(c.cl, key)
	}

	// Retrieve the class as well because plans don't have the external class name
	class, err := c.cl.ServicecatalogV1beta1().ClusterServiceClasses().Get(plan.Spec.ClusterServiceClassRef.Name, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to get class '%s' (%s)", plan.Spec.ClusterServiceClassRef.Name, err)
	}

	output.WritePlanList([]v1beta1.ClusterServicePlan{*plan}, []v1beta1.ClusterServiceClass{*class})

	return nil
}
