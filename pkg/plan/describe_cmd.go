package plan

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/Azure/service-catalog-cli/pkg/traverse"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type describeCmd struct {
	cl           *clientset.Clientset
	traverse     bool
	lookupByUUID bool
}

// NewDescribeCmd builds a "svcat describe plan" command
func NewDescribeCmd(cl *clientset.Clientset) *cobra.Command {
	describeCmd := &describeCmd{cl: cl}
	cmd := &cobra.Command{
		Use:     "plan NAME",
		Aliases: []string{"plans", "pl"},
		Short:   "Show details of a specific plan",
		Example: `
  svcat describe plan standard800
  svcat describe plan --uuid 08e4b43a-36bc-447e-a81f-8202b13e339c
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return describeCmd.run(args)
		},
	}
	cmd.Flags().BoolVarP(
		&describeCmd.traverse,
		"traverse",
		"t",
		false,
		"Whether or not to traverse from plan -> class -> broker",
	)
	cmd.Flags().BoolVarP(
		&describeCmd.lookupByUUID,
		"uuid",
		"u",
		false,
		"Whether or not to get the class by UUID (the default is by name)",
	)
	return cmd
}

func (c *describeCmd) run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("name or uuid is required")
	}

	key := args[0]
	return c.describe(key)
}

func (c *describeCmd) describe(key string) error {
	var plan *v1beta1.ClusterServicePlan
	var err error
	if c.lookupByUUID {
		plan, err = retrieveByUUID(c.cl, key)
	} else {
		plan, err = retrieveByName(c.cl, key)
	}
	if err != nil {
		return err
	}

	// Retrieve the class as well because plans don't have the external class name
	class, err := c.cl.ServicecatalogV1beta1().ClusterServiceClasses().Get(plan.Spec.ClusterServiceClassRef.Name, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to get class (%s)", err)
	}

	output.WritePlanDetails(plan, class)

	instances, err := traverse.PlanToInstances(c.cl, plan)
	if err != nil {
		return err
	}
	output.WriteAssociatedInstances(instances)

	return nil
}
