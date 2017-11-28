package binding

import (
	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

type getCmd struct {
	cl *clientset.Clientset
	ns string
}

// NewGetCmd builds a "svc-cat get bindings" command
func NewGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := getCmd{cl: cl}
	cmd := &cobra.Command{
		Use:     "bindings [name]",
		Aliases: []string{"binding", "bnd"},
		Short:   "List bindings, optionally filtered by name",
		Example: `
  svc-cat get bindings
  svc-cat get binding wordpress-mysql-binding
  svc-cat get binding -n ci concourse-postgres-binding
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
	bindings, err := retrieveAll(c.cl, c.ns)
	if err != nil {
		return err
	}

	output.WriteBindingList(bindings.Items...)
	return nil
}

func (c *getCmd) get(name string) error {
	binding, err := retrieveByName(c.cl, c.ns, name)
	if err != nil {
		return err
	}

	output.WriteBindingList(*binding)
	return nil
}
