package instance

import (
	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

type getCmd struct {
	cl *clientset.Clientset
	ns string
}

// NewGetCmd builds a "svc-cat get instances" command
func NewGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := &getCmd{cl: cl}
	cmd := &cobra.Command{
		Use:     "instances [name]",
		Aliases: []string{"instance", "inst"},
		Short:   "List instances, optionally filtered by name",
		Example: `
  svc-cat get instances
  svc-cat get instance wordpress-mysql-instance
  svc-cat get instance -n ci concourse-postgres-instance
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
		"The namespace in which to get the ServiceInstance",
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
	instances, err := retrieveAll(c.cl, c.ns)
	if err != nil {
		return err
	}

	output.WriteInstanceList(instances.Items...)
	return nil
}

func (c *getCmd) get(name string) error {
	instance, err := retrieveByName(c.cl, c.ns, name)
	if err != nil {
		return err
	}

	output.WriteInstanceList(*instance)
	return nil
}
