package broker

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

// NewRootCmd returns a new cobra command that represents the root
// of the broker command tree
func NewRootCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "broker",
		Aliases: []string{"brokers", "brk"},
		Short:   "Access ClusterServiceBrokers",
	}

	rootCmd.AddCommand(newBrokerListCmd(cl))
	rootCmd.AddCommand(newBrokerGetCmd(cl))
	return rootCmd
}
