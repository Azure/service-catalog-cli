package broker

import (
	"errors"
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newBrokerGetCmd(cl *clientset.Clientset) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get the detailed status of an individual ClusterServiceBroker",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Usage: brokers get <broker name>")
			}
			brokerName := args[0]
			broker, err := cl.Servicecatalog().ClusterServiceBrokers().Get(brokerName, v1.GetOptions{})
			if err != nil {
				return fmt.Errorf("Error getting broker (%s)", err)
			}
			logger.Printf("Broker URL: %s", broker.Spec.URL)
			t := output.NewTable()
			t.SetCaption(true, fmt.Sprintf("%d status condition(s)", len(broker.Status.Conditions)))
			output.ClusterServiceBrokerHeaders(t)
			output.AppendClusterServiceBroker(t, broker)
			t.Render()
			return nil
		},
	}
}
