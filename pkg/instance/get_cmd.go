package instance

import (
	"errors"
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/Azure/service-catalog-cli/pkg/traverse"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type instanceGetCmd struct {
	cl       *clientset.Clientset
	ns       string
	traverse bool
}

func (i *instanceGetCmd) run(args []string) error {
	if len(args) != 1 {
		return errors.New("Usage: instance get <instance name>")
	}
	instanceName := args[0]
	instance, err := i.cl.Servicecatalog().ServiceInstances(i.ns).Get(instanceName, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Error getting instance %s (%s)", instanceName, err)
	}
	t := output.NewTable(0)
	output.InstanceHeaders(t)
	output.AppendInstance(t, instance)
	t.Render()
	if !i.traverse {
		return nil
	}

	// Traverse from instance to service class and plan
	class, plan, err := traverse.InstanceToServiceClassAndPlan(i.cl, instance)
	if err != nil {
		return fmt.Errorf("Error traversing instance to its service class and plan (%s)", err)
	}
	logger.Printf("\n\nSERVICE CLASS")
	t = output.NewTable(1)
	output.ClusterServiceClassHeaders(t)
	output.AppendClusterServiceClass(t, class)
	t.Render()

	logger.Printf("\n\nSERVICE PLAN")
	t = output.NewTable(1)
	output.ClusterServicePlanHeaders(t)
	output.AppendClusterServicePlan(t, plan)
	t.Render()

	// traverse from service class to broker
	broker, err := traverse.ServiceClassToBroker(i.cl, class)
	if err != nil {
		return fmt.Errorf("Error traversing service class to broker (%s)", err)
	}
	logger.Printf("\n\nBROKER")
	t = output.NewTable(2)
	output.ClusterServiceBrokerHeaders(t)
	output.AppendClusterServiceBroker(t, broker)
	t.Render()

	return nil
}

func newInstanceGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := &instanceGetCmd{cl: cl}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get detailed information of the given instance, along with the service class/plan that instance references",
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
		"Whether or not to traverse from binding -> instance -> service class/service plan -> broker",
	)
	return cmd
}
