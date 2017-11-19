package broker

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type syncCmd struct {
	cl *clientset.Clientset
}

func (c *syncCmd) run(name string) error {
	catalog, err := c.cl.Servicecatalog().ClusterServiceBrokers().Get(name, v1.GetOptions{})
	if err != nil {
		logger.Fatalf("Error fetching ClusterServiceBrokers (%s)", err)
	}

	catalog.Spec.RelistRequests = catalog.Spec.RelistRequests + 1

	c.cl.Servicecatalog().ClusterServiceBrokers().Update(catalog)

	logger.Printf("Successfully fetched catalog entries from %s broker", name)

	return nil
}

// NewSyncCmd builds a "svc-cat sync broker" command
func NewSyncCmd(cl *clientset.Clientset) *cobra.Command {
	syncCmd := syncCmd{cl: cl}
	rootCmd := &cobra.Command{
		Use:   "broker <broker name>",
		Short: "Syncs service catalog for a service broker",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Missing service broker name")
			}
			brokerName := args[0]

			return syncCmd.run(brokerName)
		},
	}
	return rootCmd
}
