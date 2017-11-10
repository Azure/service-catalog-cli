package binding

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
	output.BindingHeaders(t)
	for _, binding := range bindings.Items {
		output.AppendBinding(t, &binding)
	}
	t.Render()
	return nil
}

func newBindingListCmd(cl *clientset.Clientset) *cobra.Command {
	listCmd := &bindingListCmd{cl: cl}
	rootCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all bindings in the given namespace, along with the instances that they belong to",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listCmd.run()
		},
	}
	rootCmd.Flags().StringVarP(
		&listCmd.ns,
		"namespace",
		"n",
		"default",
		"The namespace from which to get the binding",
	)
	return rootCmd
}
