package broker

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type syncCmd struct {
	cl *clientset.Clientset
}

// NewSyncCmd builds a "svc-cat sync broker" command
func NewSyncCmd(cl *clientset.Clientset) *cobra.Command {
	syncCmd := syncCmd{cl: cl}
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
	for j := 0; j < retries; j++ {
		catalog, err := c.cl.ServicecatalogV1beta1().ClusterServiceBrokers().Get(name, v1.GetOptions{})
		if err != nil {
			logger.Fatalf("Error fetching ClusterServiceBrokers (%s)", err)
		}

		catalog.Spec.RelistRequests = catalog.Spec.RelistRequests + 1

		_, err = c.cl.ServicecatalogV1beta1().ClusterServiceBrokers().Update(catalog)
		if err == nil {
			logger.Printf("Successfully fetched catalog entries from %s broker", name)
			return nil
		}
		if !errors.IsConflict(err) {
			return err
		}
		logger.Printf("Conflict when syncing service broker, retries left: %v", retries-j)
	}
	return fmt.Errorf("Failed to sync service broker")
}
