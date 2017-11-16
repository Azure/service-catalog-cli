package class

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

type classLookup struct {
	m map[string][]*v1beta1.ClusterServiceClass
}

func (c *classLookup) get(name string) []*v1beta1.ClusterServiceClass {
	return c.m[name]
}

func classesByName(classList *v1beta1.ClusterServiceClassList) *classLookup {
	m := make(map[string][]*v1beta1.ClusterServiceClass)
	for _, class := range classList.Items {
		// we need to copy the class since we're taking the address of it below
		classCopy := class
		lst := m[class.Spec.ExternalName]
		if lst == nil {
			lst = []*v1beta1.ClusterServiceClass{}
		}
		lst = append(lst, &classCopy)
		m[class.Spec.ExternalName] = lst
	}

	return &classLookup{m: m}

}
