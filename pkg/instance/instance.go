package instance

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	waitForStatusPollInterval = time.Second
	waitForStatusRetries      = 5
)

func retrieveAll(cl *clientset.Clientset, ns string) (*v1beta1.ServiceInstanceList, error) {
	instances, err := cl.ServicecatalogV1beta1().ServiceInstances(ns).List(v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Error listing instances (%s)", err)
	}

	return instances, nil
}

func retrieveByName(cl *clientset.Clientset, ns, name string) (*v1beta1.ServiceInstance, error) {
	instance, err := cl.ServicecatalogV1beta1().ServiceInstances(ns).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get instance '%s.%s' (%s)", ns, name, err)
	}
	return instance, nil
}

func provision(cl *clientset.Clientset, namespace, instanceName, className, planName string, params map[string]string,
) (*v1beta1.ServiceInstance, error) {
	paramsJson, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal the request parameters %v (%s)", params, err)
	}

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
			Parameters: &runtime.RawExtension{
				Raw: paramsJson,
			},
		},
	}

	result, err := cl.ServicecatalogV1beta1().ServiceInstances(namespace).Create(request)
	if err != nil {
		return nil, fmt.Errorf("provision request failed (%s)", err)
	}
	return result, nil
}

// waitForStatus to be populated on an instance.
// The result of some operations only include a subset of the instance fields, notably missing status.
func waitForStatus(cl *clientset.Clientset, instance *v1beta1.ServiceInstance) (*v1beta1.ServiceInstance, error) {
	var err error

	for i := 0; i < waitForStatusRetries; i++ {
		if len(instance.Status.Conditions) > 0 {
			break
		}

		time.Sleep(waitForStatusPollInterval)

		instance, err = retrieveByName(cl, instance.Namespace, instance.Name)
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
}
