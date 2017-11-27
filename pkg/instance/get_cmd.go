package instance

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/Azure/service-catalog-cli/pkg/traverse"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type getCmd struct {
	cl       *clientset.Clientset
	ns       string
	traverse bool
}

// NewGetCmd builds a "svc-cat get instances" command
func NewGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := &getCmd{cl: cl}
	cmd := &cobra.Command{
		Use:     "instances [name]",
		Aliases: []string{"instance", "inst"},
		Short:   "List instances, optionally filtered by name",
		Example: `
  svc-cat get instances
  svc-cat get instance wordpress-mysql-instance
  svc-cat get instance -n ci concourse-postgres-instance
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmd.run(args)
		},
	}
	cmd.Flags().StringVarP(
		&getCmd.ns,
		"namespace",
		"n",
		"default",
		"The namespace in which to get the ServiceInstance",
	)
	cmd.Flags().BoolVarP(
		&getCmd.traverse,
		"traverse",
		"t",
		false,
		"Whether or not to traverse from instance -> service class/service plan -> broker",
	)
	return cmd
}

func (c *getCmd) run(args []string) error {
	if len(args) == 0 {
		return c.getAll()
	} else {
		name := args[0]
		return c.get(name)
	}
}

func (c *getCmd) getAll() error {
	instances, err := c.cl.Servicecatalog().ServiceInstances(c.ns).List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Error listing instances (%s)", err)
	}
	t := output.NewTable()
	output.InstanceHeaders(t)
	for _, instance := range instances.Items {
		output.AppendInstance(t, &instance)
	}
	t.Render()
	return nil
}

func (c *getCmd) get(name string) error {
	instance, err := c.cl.Servicecatalog().ServiceInstances(c.ns).Get(name, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Error getting instance %s (%s)", name, err)
	}
	t := output.NewTable()
	output.InstanceHeaders(t)
	output.AppendInstance(t, instance)
	t.Render()
	if !c.traverse {
		return nil
	}

	// Traverse from instance to service class and plan
	class, plan, err := traverse.InstanceToServiceClassAndPlan(c.cl, instance)
	if err != nil {
		return fmt.Errorf("Error traversing instance to its service class and plan (%s)", err)
	}
	logger.Printf("\n\nSERVICE CLASS")
	t = output.NewTable()
	output.ClusterServiceClassHeaders(t)
	output.AppendClusterServiceClass(t, class)
	t.Render()

	logger.Printf("\n\nSERVICE PLAN")
	t = output.NewTable()
	output.ClusterServicePlanHeaders(t)
	output.AppendClusterServicePlan(t, plan)
	t.Render()

	// traverse from service class to broker
	broker, err := traverse.ServiceClassToBroker(c.cl, class)
	if err != nil {
		return fmt.Errorf("Error traversing service class to broker (%s)", err)
	}
	logger.Printf("\n\nBROKER")
	t = output.NewTable()
	output.ClusterServiceBrokerHeaders(t)
	output.AppendClusterServiceBroker(t, broker)
	t.Render()

	return nil
}
