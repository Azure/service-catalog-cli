package broker

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/command"
	"github.com/Azure/service-catalog-cli/pkg/service-catalog/client"
	"github.com/spf13/cobra"
)

type syncCmd struct {
	*command.Context
}

// NewSyncCmd builds a "svcat sync broker" command
func NewSyncCmd(cxt *command.Context) *cobra.Command {
	syncCmd := syncCmd{Context: cxt}
	rootCmd := &cobra.Command{
		Use:   "broker [name]",
		Short: "Syncs service catalog for a service broker",
		RunE: func(cmd *cobra.Command, args []string) error {
			return syncCmd.run(args)
		},
	}
	return rootCmd
}

func (c *syncCmd) run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("name is required")
	}
	name := args[0]
	return c.sync(name)
}

func (c *syncCmd) sync(name string) error {
	const retries = 3
	err := client.Sync(c.Client, name, retries)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Output, "Successfully fetched catalog entries from %s broker", name)
	return nil
}
