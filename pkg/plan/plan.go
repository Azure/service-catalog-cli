package plan

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

const (
	fieldExternalName = "spec.externalName"
)

func retrieveByName(cl *clientset.Clientset, name string) (*v1beta1.ClusterServicePlan, error) {
	opts := v1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(fieldExternalName, name).String(),
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

func retrieveByUUID(cl *clientset.Clientset, uuid string) (*v1beta1.ClusterServicePlan, error) {
	plan, err := cl.ServicecatalogV1beta1().ClusterServicePlans().Get(uuid, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get plan by uuid '%s' (%s)", uuid, err)
	}
	return plan, nil
}
