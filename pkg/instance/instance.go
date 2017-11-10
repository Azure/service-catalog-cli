package instance

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewRootCmd returns a cobra command that represents the root of the instance
// command tree
func NewRootCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "instance",
		Aliases: []string{"instances", "inst"},
	}
	rootCmd.AddCommand(newInstanceListCmd(cl))
	return rootCmd
}

type instanceListCmd struct {
	cl *clientset.Clientset
	ns string
}

func (i *instanceListCmd) run() error {
	instances, err := i.cl.Servicecatalog().ServiceInstances(i.ns).List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Error listing instances (%s)", err)
	}
	t := output.NewTable()
	output.InstanceHeaders(t)
	for _, instance := range instances.Items {
		output.AppendInstance(t, &instance)
	}
	t.Render()
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
