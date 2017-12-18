package traverse

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

const (
	fieldServiceClassRef = "spec.clusterServiceClassRef.name"
)

// ClassToBroker retrieves the parent broker of a class.
func ClassToBroker(
	cl *clientset.Clientset,
	class *v1beta1.ClusterServiceClass,
) (*v1beta1.ClusterServiceBroker, error) {
	brokerName := class.Spec.ClusterServiceBrokerName
	broker, err := cl.ServicecatalogV1beta1().ClusterServiceBrokers().Get(brokerName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return broker, nil
}

// ClassToPlans retrieves all plans for a class.
func ClassToPlans(cl *clientset.Clientset, class *v1beta1.ClusterServiceClass,
) ([]v1beta1.ClusterServicePlan, error) {
	planOpts := v1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(fieldServiceClassRef, class.Name).String(),
	}
	plans, err := cl.ServicecatalogV1beta1().ClusterServicePlans().List(planOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to list plans (%s)", err)
	}

	return plans.Items, nil
}
