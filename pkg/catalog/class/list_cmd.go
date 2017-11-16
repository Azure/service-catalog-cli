package class

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newListCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "list",
		Short: "Get a list of ClusterServiceClasses",
		RunE: func(cmd *cobra.Command, args []string) error {
			classes, err := cl.Servicecatalog().ClusterServiceClasses().List(v1.ListOptions{})
			if err != nil {
				return fmt.Errorf("Listing classes (%s)", err)
			}
			t := output.NewTable()
			output.ClusterServiceClassHeaders(t)
			for _, class := range classes.Items {
				output.AppendClusterServiceClass(t, &class)
			}
			t.Render()
			return nil
		},
	}
	return rootCmd
}
