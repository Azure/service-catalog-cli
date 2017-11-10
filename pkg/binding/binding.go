package binding

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

func NewRootCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "binding",
	}
	rootCmd.AddCommand(newBindingListCmd(cl))
	return rootCmd
}

func newBindingListCmd(cl *clientset.Clientset) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all bindings in the given namespace, along with the instances that they belong to",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Printf("binding list command")
			return nil
		},
	}
}
