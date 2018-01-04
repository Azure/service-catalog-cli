package client

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

const (
	FieldServicePlanRef = "spec.clusterServicePlanRef.name"
)

func (cl *Client) RetrieveInstances(ns string) (*v1beta1.ServiceInstanceList, error) {
	instances, err := cl.ServicecatalogV1beta1().ServiceInstances(ns).List(v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list instances in %s (%s)", ns, err)
	}

	return instances, nil
}

func (cl *Client) RetrieveInstance(ns, name string) (*v1beta1.ServiceInstance, error) {
	instance, err := cl.ServicecatalogV1beta1().ServiceInstances(ns).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get instance '%s.%s' (%s)", ns, name, err)
	}
	return instance, nil
}

// RetrieveInstanceByBinding retrieves the parent instance for a binding.
func (cl *Client) RetrieveInstanceByBinding(b *v1beta1.ServiceBinding,
) (*v1beta1.ServiceInstance, error) {
	ns := b.Namespace
	instName := b.Spec.ServiceInstanceRef.Name
	inst, err := cl.ServicecatalogV1beta1().ServiceInstances(ns).Get(instName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return inst, nil
}

// RetrieveInstancesByPlan retrieves all instances of a plan.
func (cl *Client) RetrieveInstancesByPlan(plan *v1beta1.ClusterServicePlan,
) ([]v1beta1.ServiceInstance, error) {
	planOpts := v1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(FieldServicePlanRef, plan.Name).String(),
	}
	instances, err := cl.ServicecatalogV1beta1().ServiceInstances("").List(planOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to list instances (%s)", err)
	}

	return instances.Items, nil
}

// InstanceParentHierarchy retrieves all ancestor resources of an instance.
func (cl *Client) InstanceParentHierarchy(instance *v1beta1.ServiceInstance,
) (*v1beta1.ClusterServiceClass, *v1beta1.ClusterServicePlan, *v1beta1.ClusterServiceBroker, error) {
	class, plan, err := cl.InstanceToServiceClassAndPlan(instance)
	if err != nil {
		return nil, nil, nil, err
	}

	broker, err := cl.RetrieveBrokerByClass(class)
	if err != nil {
		return nil, nil, nil, err
	}

	return class, plan, broker, nil
}

// InstanceToServiceClassAndPlan retrieves the parent class and plan for an instance.
func (cl *Client) InstanceToServiceClassAndPlan(instance *v1beta1.ServiceInstance,
) (*v1beta1.ClusterServiceClass, *v1beta1.ClusterServicePlan, error) {
	classID := instance.Spec.ClusterServiceClassRef.Name
	classCh := make(chan *v1beta1.ClusterServiceClass)
	classErrCh := make(chan error)
	go func() {
		class, err := cl.ServicecatalogV1beta1().ClusterServiceClasses().Get(classID, v1.GetOptions{})
		if err != nil {
			classErrCh <- err
			return
		}
		classCh <- class
	}()

	planID := instance.Spec.ClusterServicePlanRef.Name
	planCh := make(chan *v1beta1.ClusterServicePlan)
	planErrCh := make(chan error)
	go func() {
		plan, err := cl.ServicecatalogV1beta1().ClusterServicePlans().Get(planID, v1.GetOptions{})
		if err != nil {
			planErrCh <- err
			return
		}
		planCh <- plan
	}()

	var class *v1beta1.ClusterServiceClass
	var plan *v1beta1.ClusterServicePlan
	for {
		select {
		case cl := <-classCh:
			class = cl
			if class != nil && plan != nil {
				return class, plan, nil
			}
		case err := <-classErrCh:
			return nil, nil, err
		case pl := <-planCh:
			plan = pl
			if class != nil && plan != nil {
				return class, plan, nil
			}
		case err := <-planErrCh:
			return nil, nil, err

		}
	}
}

func (cl *Client) Provision(namespace, instanceName, className, planName string,
	params map[string]string, secrets map[string]string) (*v1beta1.ServiceInstance, error) {

	request := &v1beta1.ServiceInstance{
		ObjectMeta: v1.ObjectMeta{
			Name:      instanceName,
			Namespace: namespace,
		},
		Spec: v1beta1.ServiceInstanceSpec{
			PlanReference: v1beta1.PlanReference{
				ClusterServiceClassExternalName: className,
				ClusterServicePlanExternalName:  planName,
			},
			Parameters:     BuildParameters(params),
			ParametersFrom: BuildParametersFrom(secrets),
		},
	}

	result, err := cl.ServicecatalogV1beta1().ServiceInstances(namespace).Create(request)
	if err != nil {
		return nil, fmt.Errorf("provision request failed (%s)", err)
	}
	return result, nil
}

func (cl *Client) Deprovision(namespace, instanceName string) error {
	err := cl.ServicecatalogV1beta1().ServiceInstances(namespace).Delete(instanceName, &v1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("deprovision request failed (%s)", err)
	}
	return nil
}
