package catalog

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

// NewRootCmd creates a new cobra command that represents the root of the
// catalog command tree
func NewRootCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "catalog",
		Aliases: []string{"cat"},
	}
	rootCmd.AddCommand(newCatalogListCmd(cl))
	rootCmd.AddCommand(newCatalogGetCmd(cl))
	return rootCmd
}
