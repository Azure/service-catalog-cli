package traverse

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BindingToInstance traverses from b to the ServiceInstance that it refers to
func BindingToInstance(
	cl *clientset.Clientset,
	b *v1beta1.ServiceBinding,
) (*v1beta1.ServiceInstance, error) {
	ns := b.Namespace
	instName := b.Spec.ServiceInstanceRef.Name
	inst, err := cl.ServicecatalogV1beta1().ServiceInstances(ns).Get(instName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return inst, nil
}

// BindingParentHierarchy retrieves all ancestor resources of a binding.
func BindingParentHierarchy(cl *clientset.Clientset, binding *v1beta1.ServiceBinding,
) (*v1beta1.ServiceInstance, *v1beta1.ClusterServiceClass, *v1beta1.ClusterServicePlan, *v1beta1.ClusterServiceBroker, error) {
	instance, err := BindingToInstance(cl, binding)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	class, plan, err := InstanceToServiceClassAndPlan(cl, instance)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	broker, err := ServiceClassToBroker(cl, class)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return instance, class, plan, broker, nil
}
