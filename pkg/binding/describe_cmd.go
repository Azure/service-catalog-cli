package binding

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/Azure/service-catalog-cli/pkg/traverse"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

type describeCmd struct {
	cl       *clientset.Clientset
	ns       string
	traverse bool
}

// NewDescribeCmd builds a "svc-cat describe binding" command
func NewDescribeCmd(cl *clientset.Clientset) *cobra.Command {
	describeCmd := &describeCmd{cl: cl}
	cmd := &cobra.Command{
		Use:     "binding NAME",
		Aliases: []string{"bindings", "bnd"},
		Short:   "Show details of a specific binding",
		Example: `
  svc-cat describe binding wordpress-mysql-binding
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return describeCmd.run(args)
		},
	}
	cmd.Flags().StringVarP(
		&describeCmd.ns,
		"namespace",
		"n",
		"default",
		"The namespace in which to get the binding",
	)
	cmd.Flags().BoolVarP(
		&describeCmd.traverse,
		"traverse",
		"t",
		false,
		"Whether or not to traverse from binding -> instance -> class/plan -> broker",
	)
	return cmd
}

func (c *describeCmd) run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("name is required")
	}

	key := args[0]
	return c.describe(key)
}

func (c *describeCmd) describe(name string) error {
	binding, err := retrieveByName(c.cl, c.ns, name)
	if err != nil {
		return err
	}

	output.WriteBindingDetails(binding)

	if c.traverse {
		instance, class, plan, broker, err := traverse.BindingParentHierarchy(c.cl, binding)
		if err != nil {
			return fmt.Errorf("unable to traverse up the binding hierarchy (%s)", err)
		}
		output.WriteParentInstance(instance)
		output.WriteParentClass(class)
		output.WriteParentPlan(plan)
		output.WriteParentBroker(broker)
	}

	return nil
}
