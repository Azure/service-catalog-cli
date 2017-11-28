package instance

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
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
