package kube

import (
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateNamespace with the specified name.
func (sdk *SDK) CreateNamespace(name string) (*v1.Namespace, error) {
	request := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	result, err := sdk.Core().Namespaces().Create(request)
	if err != nil {
		return nil, fmt.Errorf("unable to create namespace %s (%s)", name, err)
	}

	return result, nil
}

// DeleteNamespace by name with the default grace period.
func (sdk *SDK) DeleteNamespace(name string) error {
	err := sdk.Core().Namespaces().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("unable to delete namespace %s (%s)", name, err)
	}

	return nil
}
