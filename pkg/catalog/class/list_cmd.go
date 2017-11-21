package class

import (
	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type classListCmd struct {
	cl       *clientset.Clientset
	traverse bool
}

func (c *classListCmd) run(args []string) error {
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

func newClassListCmd(cl *clientset.Clientset) *cobra.Command {
	listCmd := &classListCmd{cl: cl}
	return &cobra.Command{
		Use:   "list",
		Short: "Return a list of classes in the catalog",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listCmd.run(args)
		},
	}
}
