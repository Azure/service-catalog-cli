package broker

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type getCmd struct {
	cl *clientset.Clientset
}

// NewGetCmd builds a "svc-cat get brokers" command
func NewGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := getCmd{cl: cl}
	cmd := &cobra.Command{
		Use:     "brokers [name]",
		Aliases: []string{"broker", "brk"},
		Short:   "List brokers, optionally filtered by name",
		Example: `
  svc-cat get brokers
  svc-cat get broker asb
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmd.run(args)
		},
	}

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
	lst, err := c.cl.Servicecatalog().ClusterServiceBrokers().List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Error listing brokers (%s)", err)
	}
	if len(lst.Items) == 0 {
		logger.Printf("No brokers are installed!")
		return nil
	}
	table := output.NewTable()
	table.SetHeader([]string{"Name", "URL", "Status"})
	for _, broker := range lst.Items {
		latestCond := "None"
		if len(broker.Status.Conditions) >= 1 {
			latestCond = broker.Status.Conditions[len(broker.Status.Conditions)-1].Reason
		}
		table.Append([]string{
			broker.Name,
			broker.Spec.URL,
			latestCond,
		})
	}
	table.Render()

	return nil
}

func (c *getCmd) get(name string) error {
	broker, err := retrieveByName(c.cl, name)
	if err != nil {
		return err
	}

	logger.Printf("Broker URL: %s", broker.Spec.URL)
	t := output.NewTable()
	t.SetCaption(true, fmt.Sprintf("%d status condition(s)", len(broker.Status.Conditions)))
	output.ClusterServiceBrokerHeaders(t)
	output.AppendClusterServiceBroker(t, broker)
	t.Render()
	return nil
}
