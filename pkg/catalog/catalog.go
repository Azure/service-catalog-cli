package catalog

import (
	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewRootCmd creates a new cobra command that represents the root of the
// catalog command tree
func NewRootCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "catalog",
	}
	rootCmd.AddCommand(newCatalogListCmd(cl))
	return rootCmd
}
func newCatalogListCmd(cl *clientset.Clientset) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Return a list of items in the catalog",
		RunE: func(cmd *cobra.Command, args []string) error {
			classes, err := cl.Servicecatalog().ClusterServiceClasses().List(v1.ListOptions{})
			if err != nil {
				logger.Fatalf("Error fetching ClusterServiceClasses")
			}
			if len(classes.Items) == 0 {
				logger.Printf("The catalog is empty!")
				return nil
			}
			table := output.NewTable()
			table.SetHeader([]string{
				"Name",
				"Number of plans",
				"UUID",
			})
			for _, class := range classes.Items {
				table.Append([]string{
					class.Spec.ClusterServiceBrokerName,
					"TODO",
					class.Name,
				})
			}
			return nil
		},
	}
}
