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
		Use: "instance",
	}
	rootCmd.AddCommand(newInstanceListCmd(cl))
	return rootCmd
}

type instanceListCmd struct {
	cl *clientset.Clientset
	ns string
}

func (i *instanceListCmd) run() error {
	logger.Printf("instance list command")
	return nil
}

func newInstanceListCmd(cl *clientset.Clientset) *cobra.Command {
	listCmd := &instanceListCmd{cl: cl}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all instances in the given namespace, along with the service class/plan that they reference",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listCmd.run()
		},
	}
	cmd.Flags().StringVarP(&listCmd.ns, "namespace", "n", "default", "The namespace in which to list ServiceInstances")
	return cmd
}

func getServiceClassAndPlanForInstance(
	cl *clientset.Clientset,
	instance *v1beta1.ServiceInstance,
) (*v1beta1.ClusterServiceClass, *v1beta1.ClusterServicePlan, error) {
	return nil, nil, nil
}
