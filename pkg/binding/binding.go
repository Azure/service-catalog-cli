package binding

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Azure/service-catalog-cli/pkg/client"
	traverse "github.com/Azure/service-catalog-cli/pkg/traverse"
	"github.com/hashicorp/go-multierror"
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

func unbind(cl *clientset.Clientset, ns, instanceName string) error {
	instance := &v1beta1.ServiceInstance{
		ObjectMeta: v1.ObjectMeta{
			Namespace: ns,
			Name:      instanceName,
		},
	}
	bindings, err := traverse.InstanceToBindings(cl, instance)
	if err != nil {
		return err
	}

	var g sync.WaitGroup
	errs := make(chan error, len(bindings))
	for _, binding := range bindings {
		g.Add(1)
		go func(binding v1beta1.ServiceBinding) {
			defer g.Done()
			errs <- deleteBinding(cl, binding.Namespace, binding.Name)
		}(binding)
	}

	g.Wait()
	close(errs)

	// Collect any errors that occurred into a single formatted error
	bindErr := &multierror.Error{
		ErrorFormat: func(errors []error) string {
			return joinErrors("could not remove some bindings:", errors, "\n  ")
		},
	}
	for err := range errs {
		bindErr = multierror.Append(bindErr, err)
	}

	return bindErr.ErrorOrNil()
}

func deleteBinding(cl *clientset.Clientset, ns, bindingName string) error {
	err := cl.ServicecatalogV1beta1().ServiceBindings(ns).Delete(bindingName, &v1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("remove binding %s/%s failed (%s)", ns, bindingName, err)
	}
	return nil
}

func joinErrors(groupMsg string, errors []error, sep string, a ...interface{}) string {
	if len(errors) == 0 {
		return ""
	}

	msgs := make([]string, 0, len(errors)+1)
	msgs = append(msgs, fmt.Sprintf(groupMsg, a...))
	for _, err := range errors {
		msgs = append(msgs, err.Error())
	}

	return strings.Join(msgs, sep)
}
