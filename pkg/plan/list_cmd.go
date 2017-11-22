package plan

import (
	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type planListCmd struct {
	cl *clientset.Clientset
}

func (p *planListCmd) run(args []string) error {
	plans, err := p.cl.ServicecatalogV1beta1().ClusterServicePlans().List(v1.ListOptions{})
	if err != nil {
		logger.Fatalf("Error fetching ClusterServicePlans (%s)", err)
	}

	// Retrieve the classes as well because plans don't have the external class name
	classes, err := p.cl.ServicecatalogV1beta1().ClusterServiceClasses().List(v1.ListOptions{})
	if err != nil {
		logger.Fatalf("Error fetching ClusterServiceClasses (%s)", err)
	}

	output.WritePlanList(plans, classes)
	return nil
}

func newPlanListCmd(cl *clientset.Clientset) *cobra.Command {
	listCmd := &planListCmd{cl: cl}
	return &cobra.Command{
		Use:   "list",
		Short: "Return a list of service plans in the catalog",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listCmd.run(args)
		},
	}
}
