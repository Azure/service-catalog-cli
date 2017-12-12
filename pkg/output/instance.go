package output

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

func getInstanceStatusCondition(status v1beta1.ServiceInstanceStatus) v1beta1.ServiceInstanceCondition {
	if len(status.Conditions) > 0 {
		return status.Conditions[len(status.Conditions)-1]
	}
	return v1beta1.ServiceInstanceCondition{}
}

func getInstanceStatusFull(status v1beta1.ServiceInstanceStatus) string {
	lastCond := getInstanceStatusCondition(status)
	return formatStatusText(string(lastCond.Type), lastCond.Message, lastCond.LastTransitionTime)
}

func getInstanceStatusShort(status v1beta1.ServiceInstanceStatus) string {
	lastCond := getInstanceStatusCondition(status)
	return string(lastCond.Type)
}

// WriteInstanceList prints a list of instances.
func WriteInstanceList(instances ...v1beta1.ServiceInstance) {
	t := NewListTable()
	t.SetHeader([]string{
		"Name",
		"Namespace",
		"Class",
		"Plan",
		"Status",
	})

	for _, instance := range instances {
		t.Append([]string{
			instance.Name,
			instance.Namespace,
			instance.Spec.ClusterServiceClassExternalName,
			instance.Spec.ClusterServicePlanExternalName,
			getInstanceStatusShort(instance.Status),
		})
	}

	t.Render()
}

// WriteParentInstance prints identifying information for a parent instance.
func WriteParentInstance(instance *v1beta1.ServiceInstance) {
	fmt.Println("\nInstance:")
	t := NewDetailsTable()
	t.AppendBulk([][]string{
		{"Name:", instance.Name},
		{"Namespace:", instance.Namespace},
		{"Status:", getInstanceStatusShort(instance.Status)},
	})
	t.Render()
}

// WriteAssociatedInstances prints a list of instances associated with a plan.
func WriteAssociatedInstances(instances []v1beta1.ServiceInstance) {
	fmt.Println("\nInstances:")
	if len(instances) == 0 {
		fmt.Println("No instances defined")
		return
	}

	t := NewListTable()
	t.SetHeader([]string{
		"Name",
		"Namespace",
		"Status",
	})
	for _, instance := range instances {
		t.Append([]string{
			instance.Name,
			instance.Namespace,
			getInstanceStatusShort(instance.Status),
		})
	}
	t.Render()
}

// WriteInstanceDetails prints an instance.
func WriteInstanceDetails(instance *v1beta1.ServiceInstance) {
	t := NewDetailsTable()

	t.AppendBulk([][]string{
		{"Name:", instance.Name},
		{"Namespace:", instance.Namespace},
		{"Status:", getInstanceStatusFull(instance.Status)},
		{"Class:", instance.Spec.ClusterServiceClassExternalName},
		{"Plan:", instance.Spec.ClusterServicePlanExternalName},
	})

	t.Render()
}
