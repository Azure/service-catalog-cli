package plan

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	// "github.com/Azure/service-catalog-cli/pkg/traverse"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newListCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "list",
		Short: "List all the ClusterServicePlans",
		RunE: func(cmd *cobra.Command, args []string) error {
			plans, err := cl.Servicecatalog().ClusterServicePlans().List(v1.ListOptions{})
			if err != nil {
				return fmt.Errorf("Listing plans (%s)", err)
			}
			t := output.NewTable()
			output.ClusterServicePlanHeaders(t)
			for _, plan := range plans.Items {
				output.AppendClusterServicePlan(t, &plan)
			}
			t.Render()
			return nil
		},
	}
	return rootCmd
}
