package catalog

import (
	"github.com/Azure/service-catalog-cli/pkg/catalog/class"
	"github.com/Azure/service-catalog-cli/pkg/catalog/plan"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

// NewRootCmd creates a new cobra command that represents the root of the
// catalog command tree
func NewRootCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "catalog",
		Aliases: []string{"cat"},
		Short:   "Access ClusterServiceClasses and ClusterServicePlans",
	}
	rootCmd.AddCommand(newCatalogListCmd(cl))
	rootCmd.AddCommand(class.NewRootCmd(cl))
	rootCmd.AddCommand(plan.NewRootCmd(cl))
	return rootCmd
}
