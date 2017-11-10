package broker

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
