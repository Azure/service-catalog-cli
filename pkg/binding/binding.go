package binding

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/client"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func retrieveAll(cl *clientset.Clientset, ns string) (*v1beta1.ServiceBindingList, error) {
	bindings, err := cl.ServicecatalogV1beta1().ServiceBindings(ns).List(v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list bindings in %s (%s)", ns, err)
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

func bind(cl *clientset.Clientset, namespace, bindingName, instanceName, secretName string,
	params map[string]string, secrets map[string]string) (*v1beta1.ServiceBinding, error) {
	request := &v1beta1.ServiceBinding{
		ObjectMeta: v1.ObjectMeta{
			Name:      bindingName,
			Namespace: namespace,
		},
		Spec: v1beta1.ServiceBindingSpec{
			ServiceInstanceRef: v1beta1.LocalObjectReference{
				Name: instanceName,
			},
			SecretName:     secretName,
			Parameters:     client.BuildParameters(params),
			ParametersFrom: client.BuildParametersFrom(secrets),
		},
	}

	result, err := cl.ServicecatalogV1beta1().ServiceBindings(namespace).Create(request)
	if err != nil {
		return nil, fmt.Errorf("bind request failed (%s)", err)
	}

	return result, nil
}
