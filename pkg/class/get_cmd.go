package class

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/command"
	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/Azure/service-catalog-cli/pkg/service-catalog/client"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/spf13/cobra"
)

type getCmd struct {
	*command.Context
	lookupByUUID bool
}

// NewGetCmd builds a "svcat get classes" command
func NewGetCmd(cxt *command.Context) *cobra.Command {
	getCmd := &getCmd{Context: cxt}
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
	classes, err := client.RetrieveClasses(c.Client)
	if err != nil {
		return fmt.Errorf("unable to list classes (%s)", err)
	}

	output.WriteClassList(c.Output, classes...)
	return nil
}

func (c *getCmd) get(key string) error {
	var class *v1beta1.ClusterServiceClass
	var err error
	if c.lookupByUUID {
		class, err = client.RetrieveClassByID(c.Client, key)
	} else {
		class, err = client.RetrieveClassByName(c.Client, key)
	}
	if err != nil {
		return err
	}

	output.WriteClassList(c.Output, *class)
	return nil
}
