package client

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func RetrieveInstances(cl *clientset.Clientset, ns string) (*v1beta1.ServiceInstanceList, error) {
	instances, err := cl.ServicecatalogV1beta1().ServiceInstances(ns).List(v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list instances in %s (%s)", ns, err)
	}

	return instances, nil
}

func RetrieveInstance(cl *clientset.Clientset, ns, name string) (*v1beta1.ServiceInstance, error) {
	instance, err := cl.ServicecatalogV1beta1().ServiceInstances(ns).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get instance '%s.%s' (%s)", ns, name, err)
	}
	return instance, nil
}

func Provision(cl *clientset.Clientset, namespace, instanceName, className, planName string,
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

func Deprovision(cl *clientset.Clientset, namespace, instanceName string) error {
	err := cl.ServicecatalogV1beta1().ServiceInstances(namespace).Delete(instanceName, &v1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("deprovision request failed (%s)", err)
	}
	return nil
}
