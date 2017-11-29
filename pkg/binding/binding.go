package binding

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func retrieveAll(cl *clientset.Clientset, ns string) (*v1beta1.ServiceBindingList, error) {
	bindings, err := cl.ServicecatalogV1beta1().ServiceBindings(ns).List(v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Error listing bindings (%s)", err)
	}

	return bindings, nil
}

func retrieveByName(cl *clientset.Clientset, ns, name string) (*v1beta1.ServiceBinding, error) {
	binding, err := cl.ServicecatalogV1beta1().ServiceBindings(ns).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get binding '%s.%s' (%+v)", ns, name, err)
	}
	return binding, nil
}
