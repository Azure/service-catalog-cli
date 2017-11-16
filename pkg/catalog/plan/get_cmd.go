package plan

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
		return fmt.Errorf("Usage: plan get <plan name>")
	}
	planName := args[0]
	plan, err := c.cl.Servicecatalog().ClusterServicePlans().Get(planName, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Getting plan %s (%s)", planName, err)
	}
	t := output.NewTable()
	output.ClusterServicePlanHeaders(t)
	output.AppendClusterServicePlan(t, plan)
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
		Short: "Get detailed information about a ClusterServicePlan",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmd.run(args)
		},
	}
	rootCmd.Flags().BoolVarP(
		&getCmd.traverse,
		"traverse",
		"t",
		false,
		"Whether or not to traverse from plan -> class -> broker",
	)
	return rootCmd
}
