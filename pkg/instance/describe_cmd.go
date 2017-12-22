package instance

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

// NewDescribeCmd builds a "svcat describe instance" command
func NewDescribeCmd(cxt *command.Context) *cobra.Command {
	describeCmd := &describeCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:     "instance NAME",
		Aliases: []string{"instances", "inst"},
		Short:   "Show details of a specific instance",
		Example: `
  svcat describe instance wordpress-mysql-instance
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
		"The namespace in which to get the instance",
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

	name := args[0]
	return c.describe(name)
}

func (c *describeCmd) describe(name string) error {
	instance, err := retrieveByName(c.Client, c.ns, name)
	if err != nil {
		return err
	}

	output.WriteInstanceDetails(c.Output, instance)

	bindings, err := traverse.InstanceToBindings(c.Client, instance)
	if err != nil {
		return err
	}
	output.WriteAssociatedBindings(c.Output, bindings)

	if c.traverse {
		class, plan, broker, err := traverse.InstanceParentHierarchy(c.Client, instance)
		if err != nil {
			return fmt.Errorf("unable to traverse up the instance hierarchy (%s)", err)
		}
		output.WriteParentClass(c.Output, class)
		output.WriteParentPlan(c.Output, plan)
		output.WriteParentBroker(c.Output, broker)
	}

	return nil
}
