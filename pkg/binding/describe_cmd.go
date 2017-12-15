package binding

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/command"
	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/Azure/service-catalog-cli/pkg/traverse"
	"github.com/spf13/cobra"
)

type describeCmd struct {
	*command.Context
	ns       string
	traverse bool
}

// NewDescribeCmd builds a "svcat describe binding" command
func NewDescribeCmd(cxt *command.Context) *cobra.Command {
	describeCmd := &describeCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:     "binding NAME",
		Aliases: []string{"bindings", "bnd"},
		Short:   "Show details of a specific binding",
		Example: `
  svcat describe binding wordpress-mysql-binding
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
	binding, err := retrieveByName(c.Client, c.ns, name)
	if err != nil {
		return err
	}

	output.WriteBindingDetails(binding)

	if c.traverse {
		instance, class, plan, broker, err := traverse.BindingParentHierarchy(c.Client, binding)
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
