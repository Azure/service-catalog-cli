package class

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

type classGetCmd struct {
	cl           *clientset.Clientset
	traverse     bool
	lookupByUUID bool
}

const (
	fieldExternalName    = "spec.externalName"
	fieldServiceClassRef = "spec.clusterServiceClassRef.name"
)

func (c *classGetCmd) run(args []string) error {
	if len(args) != 1 {
		return errors.New("Usage: catalog class get <NAME> or catalog class get --uuid <UUID>")
	}

	var svc *v1beta1.ClusterServiceClass
	if c.lookupByUUID {
		var err error
		uuid := args[0]
		svc, err = c.cl.ServicecatalogV1beta1().ClusterServiceClasses().Get(uuid, v1.GetOptions{})
		if err != nil {
			return fmt.Errorf("unable to get ClusterServiceClass (%s)", err)
		}
	} else {
		name := args[0]
		svcOpts := v1.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(fieldExternalName, name).String(),
		}
		searchResults, err := c.cl.ServicecatalogV1beta1().ClusterServiceClasses().List(svcOpts)
		if err != nil {
			return fmt.Errorf("unable to search ClusterServiceClass (%s)", err)
		}
		if len(searchResults.Items) == 0 {
			logger.Fatalf(`"%s" not found`, name)
		}
		if len(searchResults.Items) > 1 {
			logger.Fatalf(`%s = "%s" matches more than one item`, fieldExternalName, name)
		}
		svc = &searchResults.Items[0]
	}

	output.WriteClusterServiceClass(svc)

	if c.traverse {
		planOpts := v1.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(fieldServiceClassRef, svc.Name).String(),
		}
		plans, err := c.cl.ServicecatalogV1beta1().ClusterServicePlans().List(planOpts)
		if err != nil {
			return fmt.Errorf("unable to list ClusterServicePlans (%s)", err)
		}

		fmt.Println("\nPlans:")
		if len(plans.Items) == 0 {
			fmt.Println("No plans defined")
		} else {
			output.WriteClusterServicePlanList(plans, true)
		}
	}

	return nil
}

func newClassGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := &classGetCmd{cl: cl}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Return detailed information for a given service class",
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
		"Whether or not to get the class by UUID (the default is by name)",
	)
	return cmd
}
