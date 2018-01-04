package broker

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/cmd/svcat/command"
	"github.com/Azure/service-catalog-cli/cmd/svcat/output"
	"github.com/Azure/service-catalog-cli/pkg/service-catalog/client"
	"github.com/spf13/cobra"
)

type describeCmd struct {
	*command.Context
	ns       string
	name     string
	traverse bool
}

// NewDescribeCmd builds a "svcat describe broker" command
func NewDescribeCmd(cxt *command.Context) *cobra.Command {
	describeCmd := &describeCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:     "broker NAME",
		Aliases: []string{"brokers", "brk"},
		Short:   "Show details of a specific broker",
		Example: `
  svcat describe broker asb
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return describeCmd.run(args)
		},
	}
	return cmd
}

func (c *describeCmd) run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("name is required")
	}
	c.name = args[0]

	return c.describe()
}

func (c *describeCmd) describe() error {
	broker, err := client.RetrieveBroker(c.Client, c.name)
	if err != nil {
		return err
	}

	output.WriteBrokerDetails(c.Output, broker)
	return nil
}
