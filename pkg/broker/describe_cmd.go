package broker

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

type describeCmd struct {
	cl       *clientset.Clientset
	ns       string
	traverse bool
}

// NewDescribeCmd builds a "svcat describe broker" command
func NewDescribeCmd(cl *clientset.Clientset) *cobra.Command {
	describeCmd := &describeCmd{cl: cl}
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

	key := args[0]
	return c.describe(key)
}

func (c *describeCmd) describe(name string) error {
	broker, err := retrieveByName(c.cl, name)
	if err != nil {
		return err
	}

	output.WriteBrokerDetails(broker)
	return nil
}
