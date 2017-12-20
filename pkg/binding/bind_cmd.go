package binding

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/command"
	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/Azure/service-catalog-cli/pkg/parameters"
	"github.com/spf13/cobra"
)

type bindCmd struct {
	*command.Context
	ns           string
	instanceName string
	bindingName  string
	secretName   string
	params       []string
	secrets      []string
}

// NewBindCmd builds a "svcat bind" command
func NewBindCmd(cxt *command.Context) *cobra.Command {
	bindCmd := &bindCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:   "bind INSTANCE_NAME",
		Short: "Binds an instance's metadata to a secret, which can then be used by an application to connect to the instance",
		Example: `
  svcat bind wordpress
  svcat bind wordpress-mysql-instance --name wordpress-mysql-binding --secret-name wordpress-mysql-secret
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return bindCmd.run(args)
		},
	}
	cmd.Flags().StringVarP(
		&bindCmd.ns,
		"namespace",
		"n",
		"default",
		"The instance namespace",
	)
	cmd.Flags().StringVarP(
		&bindCmd.bindingName,
		"name",
		"",
		"",
		"The name of the binding. Defaults to the name of the instance.",
	)
	cmd.Flags().StringVarP(
		&bindCmd.secretName,
		"secret-name",
		"",
		"",
		"The name of the secret. Defaults to the name of the instance.",
	)
	cmd.Flags().StringArrayVarP(&bindCmd.params, "param", "p", nil,
		"Additional parameter to use when binding the instance, format: NAME=VALUE")
	cmd.Flags().StringArrayVarP(&bindCmd.secrets, "secret", "s", nil,
		"Additional parameter, whose value is stored in a secret, to use when binding the instance, format: SECRET[KEY]")

	return cmd
}

func (c *bindCmd) run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("instance is required")
	}
	c.instanceName = args[0]

	params, err := parameters.ParseVariableAssignments(c.params)
	if err != nil {
		return fmt.Errorf("invalid --param value (%s)", err)
	}

	secrets, err := parameters.ParseKeyMaps(c.secrets)
	if err != nil {
		return fmt.Errorf("invalid --secret value (%s)", err)
	}

	return c.bind(params, secrets)
}

func (c *bindCmd) bind(params map[string]string, secrets map[string]string) error {
	binding, err := bind(c.Client, c.ns, c.bindingName, c.instanceName, c.secretName, params, secrets)
	if err != nil {
		return err
	}

	output.WriteBindingDetails(c.Output, binding)

	return nil
}
