package class

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

// NewRootCmd creates a new cobra command that represents the root of the
// 'catalog class' command tree
func NewRootCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "class",
		Aliases: []string{"cl", "classes"},
		Short:   "Access ClusterServiceClasses",
	}
	rootCmd.AddCommand(newGetCmd(cl))
	rootCmd.AddCommand(newListCmd(cl))
	return rootCmd
}
