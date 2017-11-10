package binding

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewRootCmd returns a cobra command that represents the root of the binding
// command tree
func NewRootCmd(cl *clientset.Clientset) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "binding",
	}
	rootCmd.AddCommand(newBindingListCmd(cl))
	return rootCmd
}

type bindingListCmd struct {
	cl *clientset.Clientset
	ns string
}

func (b *bindingListCmd) run() error {
	bindings, err := b.cl.Servicecatalog().ServiceBindings(b.ns).List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Error listing bindings (%s)", err)
	}

	t := output.NewTable()
	t.SetHeader([]string{
		"Name",
		"Instance Name",
		"Status",
	})
	for _, binding := range bindings.Items {
		latestCond := "None"
		if len(binding.Status.Conditions) >= 1 {
			latestCond = binding.Status.Conditions[len(binding.Status.Conditions)-1].Reason
		}
		t.Append([]string{
			binding.Name,
			binding.Spec.ServiceInstanceRef.Name,
			latestCond,
		})
	}
	t.Render()
	return nil
}

func newBindingListCmd(cl *clientset.Clientset) *cobra.Command {
	listCmd := &bindingListCmd{cl: cl}
	return &cobra.Command{
		Use:   "list",
		Short: "List all bindings in the given namespace, along with the instances that they belong to",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listCmd.run()
		},
	}
}
