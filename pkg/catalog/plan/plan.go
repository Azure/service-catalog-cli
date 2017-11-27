package plan

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

// NewRootCmd creates a new cobra command that represents the root of the
// catalog plan command tree
func NewRootCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "plan",
		Aliases: []string{"plans"},
	}

	rootCmd.AddCommand(newPlanGetCmd(cl))
	rootCmd.AddCommand(newPlanListCmd(cl))
	return rootCmd
}
