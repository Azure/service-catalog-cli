package plan

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

const (
	fieldExternalName    = "spec.externalName"
	fieldServiceClassRef = "spec.clusterServicePlanRef.name"
)

type getCmd struct {
	cl           *clientset.Clientset
	traverse     bool
	lookupByUUID bool
}

// NewGetCmd builds a "svc-cat get plans" command
func NewGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := &getCmd{cl: cl}
	cmd := &cobra.Command{
		Use:     "plans [name]",
		Aliases: []string{"plan", "pl"},
		Short:   "List plans, optionally filtered by name",
		Example: `
  svc-cat get plans
  svc-cat get plan standard800
  svc-cat get plan --uuid 08e4b43a-36bc-447e-a81f-8202b13e339c
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmd.run(args)
		},
	}
	cmd.Flags().BoolVarP(
		&getCmd.traverse,
		"traverse",
		"t",
		false,
		"Whether or not to traverse from plan -> class -> broker",
	)
	cmd.Flags().BoolVarP(
		&getCmd.lookupByUUID,
		"uuid",
		"u",
		false,
		"Whether or not to get the plan by UUID (the default is by name)",
	)
	return cmd
}

func (c *getCmd) run(args []string) error {
	if len(args) == 0 {
		return c.getAll()
	} else {
		key := args[0]
		return c.get(key)
	}
}

func (c *getCmd) getAll() error {
	plans, err := c.cl.ServicecatalogV1beta1().ClusterServicePlans().List(v1.ListOptions{})
	if err != nil {
		logger.Fatalf("Error fetching ClusterServicePlans (%s)", err)
	}

	// Retrieve the classes as well because plans don't have the external class name
	classes, err := c.cl.ServicecatalogV1beta1().ClusterServiceClasses().List(v1.ListOptions{})
	if err != nil {
		logger.Fatalf("Error fetching ClusterServiceClasses (%s)", err)
	}

	output.WritePlanList(plans, classes)
	return nil
}

func (c *getCmd) get(key string) error {
	var plan *v1beta1.ClusterServicePlan
	if c.lookupByUUID {
		var err error
		uuid := key
		plan, err = c.cl.ServicecatalogV1beta1().ClusterServicePlans().Get(uuid, v1.GetOptions{})
		if err != nil {
			return fmt.Errorf("unable to get ClusterServicePlan (%s)", err)
		}
	} else {
		name := key
		svcOpts := v1.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(fieldExternalName, name).String(),
		}
		searchResults, err := c.cl.ServicecatalogV1beta1().ClusterServicePlans().List(svcOpts)
		if err != nil {
			return fmt.Errorf("unable to search ClusterServicePlans (%s)", err)
		}
		if len(searchResults.Items) == 0 {
			logger.Fatalf(`ClusterServicePlan "%s" not found`, name)
		}
		if len(searchResults.Items) > 1 {
			logger.Fatalf(`ClusterServicePlan %s = "%s" matches more than one item`, fieldExternalName, name)
		}
		plan = &searchResults.Items[0]
	}

	// Retrieve the class as well because plans don't have the external class name
	class, err := c.cl.ServicecatalogV1beta1().ClusterServiceClasses().Get(plan.Spec.ClusterServiceClassRef.Name, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to get ClusterServiceClass (%s)", err)
	}

	output.WritePlanDetails(plan, class)

	if c.traverse {
		planOpts := v1.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(fieldServiceClassRef, plan.Name).String(),
		}
		instances, err := c.cl.ServicecatalogV1beta1().ServiceInstances("").List(planOpts)
		if err != nil {
			return fmt.Errorf("unable to list ServiceInstances (%s)", err)
		}
		output.WriteAssociatedInstances(instances)
	}

	return nil
}
