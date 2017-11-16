package class

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
) ([]*v1beta1.ClusterServiceClass, error) {
	class, err := cl.Servicecatalog().ClusterServiceClasses().Get(uuid, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("Getting the class with UUID %s (%s)", uuid, err)
	}
	return []*v1beta1.ClusterServiceClass{class}, nil
}

func getByName(
	cl *clientset.Clientset,
	className string,
) ([]*v1beta1.ClusterServiceClass, error) {
	classes, err := cl.Servicecatalog().ClusterServiceClasses().List(v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Getting a list of classes (%s)", err)
	}
	classesByName := classesByName(classes)
	retClasses := classesByName.get(className)
	if len(retClasses) == 0 {
		return nil, fmt.Errorf("no class with name %s was found", className)
	}
	if len(retClasses) > 1 {
		logger.Printf("%d classes found for name %s", len(retClasses), className)
	}
	return retClasses, nil
}

func (c *getCmd) run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Usage: class get <class name>")
	}
	className := args[0]
	classes := []*v1beta1.ClusterServiceClass{}
	if c.lookupByUUID {
		cls, err := getByUUID(c.cl, className)
		if err != nil {
			return err
		}
		classes = cls
	} else {
		cls, err := getByName(c.cl, className)
		if err != nil {
			return err
		}
		classes = cls
	}
	t := output.NewTable()
	output.ClusterServiceClassHeaders(t)
	for _, class := range classes {
		output.AppendClusterServiceClass(t, class)
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
		Short: "get <class name>",
		Long:  "Get detailed information about a ClusterServiceClass",
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
	rootCmd.Flags().BoolVarP(
		&getCmd.lookupByUUID,
		"uuid",
		"u",
		false,
		"Whether or not to get the ServiceClass by UUID (the default is by name)",
	)
	return rootCmd
}
