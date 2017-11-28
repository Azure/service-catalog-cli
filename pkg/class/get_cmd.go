package class

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
	fieldServiceClassRef = "spec.clusterServiceClassRef.name"
)

type getCmd struct {
	cl           *clientset.Clientset
	traverse     bool
	lookupByUUID bool
}

// NewGetCmd builds a "svc-cat get classes" command
func NewGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := &getCmd{cl: cl}
	cmd := &cobra.Command{
		Use:     "classes [name]",
		Aliases: []string{"class", "cl"},
		Short:   "List classes, optionally filtered by name",
		Example: `
  svc-cat get classes
  svc-cat get class azure-mysqldb
  svc-cat get class --uuid 997b8372-8dac-40ac-ae65-758b4a5075a5
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
		"Whether or not to get the class by UUID (the default is by name)",
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
	classes, err := c.cl.Servicecatalog().ClusterServiceClasses().List(v1.ListOptions{})
	if err != nil {
		logger.Fatalf("Error fetching ClusterServiceClasses (%s)", err)
	}
	if len(classes.Items) == 0 {
		logger.Printf("The catalog is empty!")
		return nil
	}
	output.WriteClusterServiceClassList(classes.Items)
	return nil
}

func (c *getCmd) get(key string) error {
	var svc *v1beta1.ClusterServiceClass
	if c.lookupByUUID {
		var err error
		uuid := key
		svc, err = c.cl.ServicecatalogV1beta1().ClusterServiceClasses().Get(uuid, v1.GetOptions{})
		if err != nil {
			return fmt.Errorf("unable to get ClusterServiceClass (%s)", err)
		}
	} else {
		name := key
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

		output.WriteAssociatedPlans(plans)
	}

	return nil
}
