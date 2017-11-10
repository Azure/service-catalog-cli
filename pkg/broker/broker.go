package broker

import (
	"errors"
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewRootCmd returns a new cobra command that represents the root
// of the broker command tree
func NewRootCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "broker",
		Aliases: []string{"brokers", "brk"},
	}

	rootCmd.AddCommand(newBrokerListCmd(cl))
	rootCmd.AddCommand(newBrokerStatusCmd(cl))
	return rootCmd
}

func newBrokerListCmd(cl *clientset.Clientset) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Get a list of ClusterServiceBrokers",
		RunE: func(cmd *cobra.Command, args []string) error {
			lst, err := cl.Servicecatalog().ClusterServiceBrokers().List(v1.ListOptions{})
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
		},
	}
}
func newBrokerStatusCmd(cl *clientset.Clientset) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Get the detailed status of an individual ClusterServiceBroker",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Usage: brokers status <broker name>")
			}
			brokerName := args[0]
			broker, err := cl.Servicecatalog().ClusterServiceBrokers().Get(brokerName, v1.GetOptions{})
			if err != nil {
				return fmt.Errorf("Error getting broker (%s)", err)
			}
			logger.Printf("Broker URL: %s", broker.Spec.URL)
			t := output.NewTable()
			t.SetCaption(true, fmt.Sprintf("%d status condition(s)", len(broker.Status.Conditions)))
			t.SetHeader([]string{
				"Type",
				"Status",
				"Reason",
				"Message",
			})
			for _, cond := range broker.Status.Conditions {
				t.Append([]string{string(cond.Type),
					string(cond.Status),
					cond.Reason,
					cond.Message,
				})
			}
			t.Render()
			return nil
		},
	}
}
