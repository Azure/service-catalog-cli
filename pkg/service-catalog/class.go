package servicecatalog

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

const (
	FieldExternalClassName = "spec.externalName"
)

func (cl *SDK) RetrieveClasses() ([]v1beta1.ClusterServiceClass, error) {
	classes, err := cl.ServicecatalogV1beta1().ClusterServiceClasses().List(v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list classes (%s)", err)
	}

	return classes.Items, nil
}

func (cl *SDK) RetrieveClassByName(name string) (*v1beta1.ClusterServiceClass, error) {
	opts := v1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(FieldExternalClassName, name).String(),
	}
	searchResults, err := cl.ServicecatalogV1beta1().ClusterServiceClasses().List(opts)
	if err != nil {
		return nil, fmt.Errorf("unable to search classes by name (%s)", err)
	}
	if len(searchResults.Items) == 0 {
		return nil, fmt.Errorf("class '%s' not found", name)
	}
	if len(searchResults.Items) > 1 {
		return nil, fmt.Errorf("more than one matching class found for '%s'", name)
	}
	return &searchResults.Items[0], nil
}

func (cl *SDK) RetrieveClassByID(uuid string) (*v1beta1.ClusterServiceClass, error) {
	class, err := cl.ServicecatalogV1beta1().ClusterServiceClasses().Get(uuid, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get class (%s)", err)
	}
	return class, nil
}

func (cl *SDK) RetrieveClassByPlan(plan *v1beta1.ClusterServicePlan,
) (*v1beta1.ClusterServiceClass, error) {
	// Retrieve the class as well because plans don't have the external class name
	class, err := cl.ServicecatalogV1beta1().ClusterServiceClasses().Get(plan.Spec.ClusterServiceClassRef.Name, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get class (%s)", err)
	}

	return class, nil
}
