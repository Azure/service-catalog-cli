package instance

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/command"
	"github.com/Azure/service-catalog-cli/pkg/service-catalog/client"
	"github.com/spf13/cobra"
)

type deprovisonCmd struct {
	*command.Context
	ns string
}

// NewDeprovisionCmd builds a "svcat deprovision" command
func NewDeprovisionCmd(cxt *command.Context) *cobra.Command {
	deprovisonCmd := &deprovisonCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:   "deprovision NAME",
		Short: "Deletes an instance of a service",
		Example: `
  svcat deprovision wordpress-mysql-instance
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deprovisonCmd.run(args)
		},
	}
	cmd.Flags().StringVarP(&deprovisonCmd.ns, "namespace", "n", "default",
		"The namespace of the instance")
	return cmd
}

func (c *deprovisonCmd) run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("name is required")
	}

	key := args[0]
	return c.deprovision(key)
}

func (c *deprovisonCmd) deprovision(instanceName string) error {
	return client.Deprovision(c.Client, c.ns, instanceName)
}
