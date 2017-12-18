package instance

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/command"
	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/Azure/service-catalog-cli/pkg/params"
	"github.com/spf13/cobra"
)

type provisonCmd struct {
	*command.Context
	ns        string
	className string
	planName  string
	params    []string
}

// NewProvisionCmd builds a "svcat provision" command
func NewProvisionCmd(cxt *command.Context) *cobra.Command {
	provisionCmd := &provisonCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:   "provision NAME",
		Short: "Create a new instance of a service",
		Example: `
  svcat provision wordpress-mysql-instance --class azure-mysqldb --plan standard800 -p location=eastus -p sslEnforcement=disabled
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return provisionCmd.run(args)
		},
	}
	cmd.Flags().StringVarP(&provisionCmd.ns, "namespace", "n", "default",
		"The namespace in which to create the instance")
	cmd.Flags().StringVar(&provisionCmd.className, "class", "",
		"The class name")
	cmd.MarkFlagRequired("class")
	cmd.Flags().StringVar(&provisionCmd.planName, "plan", "",
		"The plan name")
	cmd.MarkFlagRequired("plan")
	cmd.Flags().StringArrayVarP(&provisionCmd.params, "params", "p", nil,
		"Additional parameters to use when provisioning the service")
	return cmd
}

func (c *provisonCmd) run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("name is required")
	}

	key := args[0]
	return c.provision(key)
}

func (c *provisonCmd) provision(instanceName string) error {
	params, err := params.ParseVariableAssignments(c.params)
	if err != nil {
		return fmt.Errorf("invalid --param value (%s)", err)
	}

	instance, err := provision(c.Client, c.ns, instanceName, c.className, c.planName, params)
	if err != nil {
		return err
	}

	instance, err = waitForStatus(c.Client, instance)
	if err != nil {
		return err
	}

	output.WriteInstanceDetails(c.Output, instance)

	return nil
}
