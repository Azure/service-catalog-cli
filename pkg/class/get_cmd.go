package class

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type getCmd struct {
	cl           *clientset.Clientset
	lookupByUUID bool
}

// NewGetCmd builds a "svcat get classes" command
func NewGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := &getCmd{cl: cl}
	cmd := &cobra.Command{
		Use:     "classes [name]",
		Aliases: []string{"class", "cl"},
		Short:   "List classes, optionally filtered by name",
		Example: `
  svcat get classes
  svcat get class azure-mysqldb
  svcat get class --uuid 997b8372-8dac-40ac-ae65-758b4a5075a5
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmd.run(args)
		},
	}
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
	classes, err := c.cl.ServicecatalogV1beta1().ClusterServiceClasses().List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("unable to list classes (%s)", err)
	}

	output.WriteClassList(classes.Items...)
	return nil
}

func (c *getCmd) get(key string) error {
	var class *v1beta1.ClusterServiceClass
	var err error
	if c.lookupByUUID {
		class, err = retrieveByUUID(c.cl, key)
	} else {
		class, err = retrieveByName(c.cl, key)
	}
	if err != nil {
		return err
	}

	output.WriteClassList(*class)
	return nil
}
