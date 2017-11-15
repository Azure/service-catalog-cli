package instance

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

// NewRootCmd returns a cobra command that represents the root of the instance
// command tree
func NewRootCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "instance",
		Aliases: []string{"instances", "inst"},
	}
	rootCmd.AddCommand(newInstanceListCmd(cl))
	rootCmd.AddCommand(newInstanceGetCmd(cl))
	return rootCmd
}

func getServiceClassAndPlanForInstance(
	cl *clientset.Clientset,
	instance *v1beta1.ServiceInstance,
) (*v1beta1.ClusterServiceClass, *v1beta1.ClusterServicePlan, error) {
	return nil, nil, nil
}
