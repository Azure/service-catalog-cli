package instance

import (
	"errors"
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type instanceGetCmd struct {
	cl *clientset.Clientset
	ns string
}

func (i *instanceGetCmd) run(args []string) error {
	if len(args) != 1 {
		return errors.New("Usage: instance get <instance name>")
	}
	instanceName := args[0]
	instance, err := i.cl.Servicecatalog().ServiceInstances(i.ns).Get(instanceName, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Error getting instance %s (%s)", instanceName, err)
	}
	t := output.NewTable()
	output.InstanceHeaders(t)
	output.AppendInstance(t, instance)
	t.Render()
	return nil
}

func newInstanceGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := &instanceGetCmd{cl: cl}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get detailed information of the given instance, along with the service class/plan that instance references",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmd.run(args)
		},
	}
	cmd.Flags().StringVarP(&getCmd.ns, "namespace", "n", "default", "The namespace in which to get the ServiceInstance")
	return cmd
}
