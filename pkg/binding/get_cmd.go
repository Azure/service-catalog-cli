package binding

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type bindingGetCmd struct {
	cl *clientset.Clientset
	ns string
}

func (b *bindingGetCmd) run(name string) error {
	binding, err := b.cl.Servicecatalog().ServiceBindings(b.ns).Get(name, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Error getting binding (%s)", err)
	}
	t := output.NewTable()
	output.BindingHeaders(t)
	output.AppendBinding(t, binding)
	t.Render()
	logger.Printf("\n%s binds to the following instance:\n\n", binding.Name)
	instance, err := traverseBindingToInstance(b.cl, binding)
	if err != nil {
		return fmt.Errorf("Error traversing binding to its instance (%s)", err)
	}
	t = output.NewTable()
	output.InstanceHeaders(t)
	output.AppendInstance(t, instance)
	t.Render()
	return nil
}

func newBindingGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := bindingGetCmd{cl: cl}
	rootCmd := &cobra.Command{
		Use:   "get",
		Short: "svc-cat binding get -n <namespace> <binding name>",
		Long:  "Get a specific binding along with the instance that it points to",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Missing binding name")
			}
			bindingName := args[0]
			return getCmd.run(bindingName)
		},
	}

	rootCmd.Flags().StringVarP(
		&getCmd.ns,
		"namespace",
		"n",
		"default",
		"The namespace from which to get the binding",
	)
	return rootCmd
}

// traverseBindingToInstance traverses from b to the ServiceInstance that it refers to
func traverseBindingToInstance(
	cl *clientset.Clientset,
	b *v1beta1.ServiceBinding,
) (*v1beta1.ServiceInstance, error) {
	ns := b.Namespace
	instName := b.Spec.ServiceInstanceRef.Name
	inst, err := cl.Servicecatalog().ServiceInstances(ns).Get(instName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return inst, nil
}
