package traverse

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

const (
	fieldServicePlanRef = "spec.clusterServicePlanRef.name"
)

func PlanToClass(cl *clientset.Clientset, plan *v1beta1.ClusterServicePlan,
) (*v1beta1.ClusterServiceClass, error) {
	// Retrieve the class as well because plans don't have the external class name
	class, err := cl.ServicecatalogV1beta1().ClusterServiceClasses().Get(plan.Spec.ClusterServiceClassRef.Name, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get class (%s)", err)
	}

	return class, nil
}

// PlanToInstances retrieves all instances of a plan.
func PlanToInstances(cl *clientset.Clientset, plan *v1beta1.ClusterServicePlan,
) ([]v1beta1.ServiceInstance, error) {
	planOpts := v1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(fieldServicePlanRef, plan.Name).String(),
	}
	instances, err := cl.ServicecatalogV1beta1().ServiceInstances("").List(planOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to list instances (%s)", err)
	}

	return instances.Items, nil
}
