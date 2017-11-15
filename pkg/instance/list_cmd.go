package instance

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
