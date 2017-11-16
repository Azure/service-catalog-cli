package plan

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

type planLookup struct {
	m map[string][]*v1beta1.ClusterServicePlan
}

func (c *planLookup) get(name string) []*v1beta1.ClusterServicePlan {
	return c.m[name]
}

func plansByName(planList *v1beta1.ClusterServicePlanList) *planLookup {
	m := make(map[string][]*v1beta1.ClusterServicePlan)
	for _, plan := range planList.Items {
		// we need to copy the class since we're taking the address of it below
		planCopy := plan
		lst := m[plan.Spec.ExternalName]
		if lst == nil {
			lst = []*v1beta1.ClusterServicePlan{}
		}
		lst = append(lst, &planCopy)
		m[plan.Spec.ExternalName] = lst
	}

	return &planLookup{m: m}

}
