package plan

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/command"
	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type getCmd struct {
	*command.Context
	lookupByUUID bool
}

// NewGetCmd builds a "svcat get plans" command
func NewGetCmd(cxt *command.Context) *cobra.Command {
	getCmd := &getCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:     "plans [name]",
		Aliases: []string{"plan", "pl"},
		Short:   "List plans, optionally filtered by name",
		Example: `
  svcat get plans
  svcat get plan standard800
  svcat get plan --uuid 08e4b43a-36bc-447e-a81f-8202b13e339c
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
	plans, err := c.Client.ServicecatalogV1beta1().ClusterServicePlans().List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("unable to list plans (%s)", err)
	}

	// Retrieve the classes as well because plans don't have the external class name
	classes, err := c.Client.ServicecatalogV1beta1().ClusterServiceClasses().List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("unable to list classes (%s)", err)
	}

	output.WritePlanList(c.Output, plans.Items, classes.Items)
	return nil
}

func (c *getCmd) get(key string) error {
	var plan *v1beta1.ClusterServicePlan
	var err error
	if c.lookupByUUID {
		plan, err = retrieveByUUID(c.Client, key)
	} else {
		plan, err = retrieveByName(c.Client, key)
	}

	// Retrieve the class as well because plans don't have the external class name
	class, err := c.Client.ServicecatalogV1beta1().ClusterServiceClasses().Get(plan.Spec.ClusterServiceClassRef.Name, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to get class '%s' (%s)", plan.Spec.ClusterServiceClassRef.Name, err)
	}

	output.WritePlanList(c.Output, []v1beta1.ClusterServicePlan{*plan}, []v1beta1.ClusterServiceClass{*class})

	return nil
}
