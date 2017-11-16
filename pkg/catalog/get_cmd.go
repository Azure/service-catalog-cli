package catalog

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

func newCatalogGetCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "get",
		Short: "Get detailed information about either a ClusterServiceClass or ClusterServicePlan",
	}
	rootCmd.AddCommand(newCatalogClassGetCmd(cl))
	rootCmd.AddCommand(newCatalogPlanGetCmd(cl))
	return rootCmd
}
