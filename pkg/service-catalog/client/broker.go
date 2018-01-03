package client

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func RetrieveBroker(cl *clientset.Clientset, name string) (*v1beta1.ClusterServiceBroker, error) {
	broker, err := cl.ServicecatalogV1beta1().ClusterServiceBrokers().Get(name, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get broker '%s' (%s)", name, err)
	}

	return broker, nil
}

// RetrieveBrokerByClass retrieves the parent broker of a class.
func RetrieveBrokerByClass(
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
