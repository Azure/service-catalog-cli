package traverse

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstanceToServiceClassAndPlan fetches the ClusterServiceClass and
// ClusterServicePlan for instance, using cl to do the fetches
func InstanceToServiceClassAndPlan(
	cl *clientset.Clientset,
	instance *v1beta1.ServiceInstance,
) (*v1beta1.ClusterServiceClass, *v1beta1.ClusterServicePlan, error) {
	classID := instance.Spec.ClusterServiceClassRef.Name
	classCh := make(chan *v1beta1.ClusterServiceClass)
	classErrCh := make(chan error)
	go func() {
		class, err := cl.Servicecatalog().ClusterServiceClasses().Get(classID, v1.GetOptions{})
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
		plan, err := cl.Servicecatalog().ClusterServicePlans().Get(planID, v1.GetOptions{})
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

// InstanceParentHierarchy retrieves all ancestor resources of an instance.
func InstanceParentHierarchy(cl *clientset.Clientset, instance *v1beta1.ServiceInstance,
) (*v1beta1.ClusterServiceClass, *v1beta1.ClusterServicePlan, *v1beta1.ClusterServiceBroker, error) {
	class, plan, err := InstanceToServiceClassAndPlan(cl, instance)
	if err != nil {
		return nil, nil, nil, err
	}

	broker, err := ServiceClassToBroker(cl, class)
	if err != nil {
		return nil, nil, nil, err
	}

	return class, plan, broker, nil
}
