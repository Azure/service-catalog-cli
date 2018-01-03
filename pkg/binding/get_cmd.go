package binding

import (
	"github.com/Azure/service-catalog-cli/pkg/command"
	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/Azure/service-catalog-cli/pkg/service-catalog/client"
	"github.com/spf13/cobra"
)

type getCmd struct {
	*command.Context
	ns string
}

// NewGetCmd builds a "svcat get bindings" command
func NewGetCmd(cxt *command.Context) *cobra.Command {
	getCmd := getCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:     "bindings [name]",
		Aliases: []string{"binding", "bnd"},
		Short:   "List bindings, optionally filtered by name",
		Example: `
  svcat get bindings
  svcat get binding wordpress-mysql-binding
  svcat get binding -n ci concourse-postgres-binding
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmd.run(args)
		},
	}

	cmd.Flags().StringVarP(
		&getCmd.ns,
		"namespace",
		"n",
		"default",
		"The namespace from which to get the bindings",
	)
	return cmd
}

func (c *getCmd) run(args []string) error {
	if len(args) == 0 {
		return c.getAll()
	} else {
		name := args[0]
		return c.get(name)
	}
}

func (c *getCmd) getAll() error {
	bindings, err := client.RetrieveBindings(c.Client, c.ns)
	if err != nil {
		return err
	}

	output.WriteBindingList(c.Output, bindings.Items...)
	return nil
}

func (c *getCmd) get(name string) error {
	binding, err := client.RetrieveBinding(c.Client, c.ns, name)
	if err != nil {
		return err
	}

	output.WriteBindingList(c.Output, *binding)
	return nil
}
