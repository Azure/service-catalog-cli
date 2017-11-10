package binding

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

// NewRootCmd returns a cobra command that represents the root of the binding
// command tree
func NewRootCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "binding",
		Aliases: []string{"bindings", "bnd"},
	}
	rootCmd.AddCommand(newBindingListCmd(cl))
	rootCmd.AddCommand(newBindingGetCmd(cl))
	return rootCmd
}
