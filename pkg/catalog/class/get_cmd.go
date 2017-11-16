package class

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	// "github.com/Azure/service-catalog-cli/pkg/traverse"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type getCmd struct {
	cl       *clientset.Clientset
	traverse bool
}

func (c *getCmd) run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Usage: class get <plan name>")
	}
	className := args[0]
	plan, err := c.cl.Servicecatalog().ClusterServiceClasses().Get(className, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Getting class %s (%s)", className, err)
	}
	t := output.NewTable()
	output.ClusterServiceClassHeaders(t)
	output.AppendClusterServiceClass(t, plan)
	t.Render()
	if !c.traverse {
		return nil
	}
	// TODO: traversal
	return nil
}

func newGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := &getCmd{cl: cl}
	rootCmd := &cobra.Command{
		Use:   "get",
		Short: "Get detailed information about a ClusterServiceClass",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmd.run(args)
		},
	}
	rootCmd.Flags().BoolVarP(
		&getCmd.traverse,
		"traverse",
		"t",
		false,
		"Whether or not to traverse from class -> broker",
	)
	return rootCmd
}
