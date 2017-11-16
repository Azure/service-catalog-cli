package catalog

import (
	// "github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	// "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newCatalogPlanGetCmd(cl *clientset.Clientset) *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Get detailed information about a ClusterServicePlan",
	}
}
