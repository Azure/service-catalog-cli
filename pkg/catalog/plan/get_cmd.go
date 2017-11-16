package plan

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	// "github.com/Azure/service-catalog-cli/pkg/traverse"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type getCmd struct {
	cl           *clientset.Clientset
	traverse     bool
	lookupByUUID bool
}

func getByUUID(
	cl *clientset.Clientset,
	uuid string,
) ([]*v1beta1.ClusterServicePlan, error) {
	plan, err := cl.Servicecatalog().ClusterServicePlans().Get(uuid, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("Getting the plan with UUID %s (%s)", uuid, err)
	}
	return []*v1beta1.ClusterServicePlan{plan}, nil
}

func getByName(
	cl *clientset.Clientset,
	planName string,
) ([]*v1beta1.ClusterServicePlan, error) {
	plans, err := cl.Servicecatalog().ClusterServicePlans().List(v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Getting a list of classes (%s)", err)
	}
	planLookup := plansByName(plans)
	retPlans := planLookup.get(planName)
	if len(retPlans) == 0 {
		return nil, fmt.Errorf("no plan with name %s was found", planName)
	}
	if len(retPlans) > 1 {
		logger.Printf("%d plans found for name %s", len(retPlans), planName)
	}
	return retPlans, nil
}

func (c *getCmd) run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Usage: plan get <plan name>")
	}
	planName := args[0]
	plans := []*v1beta1.ClusterServicePlan{}
	if c.lookupByUUID {
		pl, err := getByUUID(c.cl, planName)
		if err != nil {
			return err
		}
		plans = pl
	} else {
		pl, err := getByName(c.cl, planName)
		if err != nil {
			return err
		}
		plans = pl
	}
	t := output.NewTable()
	output.ClusterServicePlanHeaders(t)
	for _, plan := range plans {
		output.AppendClusterServicePlan(t, plan)
		// TODO: traversal
		// if !c.traverse {
		// 	return nil
		// }
	}
	t.Render()
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
	rootCmd.Flags().BoolVarP(
		&getCmd.lookupByUUID,
		"uuid",
		"u",
		false,
		"Whether or not to get the ServiceClass by UUID (the default is by name)",
	)
	return rootCmd
}
