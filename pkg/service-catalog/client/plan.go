package client

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

const (
	FieldExternalPlanName = "spec.externalName"
	FieldServiceClassRef  = "spec.clusterServiceClassRef.name"
)

func (cl *Client) RetrievePlans() ([]v1beta1.ClusterServicePlan, error) {
	plans, err := cl.ServicecatalogV1beta1().ClusterServicePlans().List(v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list plans (%s)", err)
	}

	return plans.Items, nil
}

func (cl *Client) RetrievePlanByName(name string) (*v1beta1.ClusterServicePlan, error) {
	opts := v1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(FieldExternalPlanName, name).String(),
	}
	searchResults, err := cl.ServicecatalogV1beta1().ClusterServicePlans().List(opts)
	if err != nil {
		return nil, fmt.Errorf("unable to search plans by name '%s', (%s)", name, err)
	}
	if len(searchResults.Items) == 0 {
		return nil, fmt.Errorf("plan not found '%s'", name)
	}
	if len(searchResults.Items) > 1 {
		return nil, fmt.Errorf("more than one matching plan found for '%s'", name)
	}
	return &searchResults.Items[0], nil
}

func (cl *Client) RetrievePlanByID(uuid string) (*v1beta1.ClusterServicePlan, error) {
	plan, err := cl.ServicecatalogV1beta1().ClusterServicePlans().Get(uuid, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get plan by uuid '%s' (%s)", uuid, err)
	}
	return plan, nil
}

// RetrievePlansByClass retrieves all plans for a class.
func (cl *Client) RetrievePlansByClass(class *v1beta1.ClusterServiceClass,
) ([]v1beta1.ClusterServicePlan, error) {
	planOpts := v1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(FieldServiceClassRef, class.Name).String(),
	}
	plans, err := cl.ServicecatalogV1beta1().ClusterServicePlans().List(planOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to list plans (%s)", err)
	}

	return plans.Items, nil
}
