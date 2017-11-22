package plan

import (
	"errors"
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

type planGetCmd struct {
	cl           *clientset.Clientset
	traverse     bool
	lookupByUUID bool
}

const (
	fieldExternalName    = "spec.externalName"
	fieldServiceClassRef = "spec.clusterServicePlanRef.name"
)

func (c *planGetCmd) run(args []string) error {
	if len(args) != 1 {
		return errors.New("Usage: catalog plan get <NAME> or catalog plan get --uuid <UUID>")
	}

	var plan *v1beta1.ClusterServicePlan
	if c.lookupByUUID {
		var err error
		uuid := args[0]
		plan, err = c.cl.ServicecatalogV1beta1().ClusterServicePlans().Get(uuid, v1.GetOptions{})
		if err != nil {
			return fmt.Errorf("unable to get ClusterServicePlan (%s)", err)
		}
	} else {
		name := args[0]
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

func newPlanGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := &planGetCmd{cl: cl}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Return detailed information for a given plan",
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
