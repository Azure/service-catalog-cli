package output

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

func getBindingStatusCondition(status v1beta1.ServiceBindingStatus) v1beta1.ServiceBindingCondition {
	if len(status.Conditions) > 0 {
		return status.Conditions[len(status.Conditions)-1]
	}
	return v1beta1.ServiceBindingCondition{}
}

func getBindingStatusShort(status v1beta1.ServiceBindingStatus) string {
	lastCond := getBindingStatusCondition(status)
	return string(lastCond.Type)
}

func getBindingStatusFull(status v1beta1.ServiceBindingStatus) string {
	lastCond := getBindingStatusCondition(status)
	return formatStatusText(string(lastCond.Type), lastCond.Message, lastCond.LastTransitionTime)
}

// WriteBindingDetails prints a list of bindings.
func WriteBindingList(bindings ...v1beta1.ServiceBinding) {
	t := NewListTable()
	t.SetHeader([]string{
		"Name",
		"Namespace",
		"Instance",
		"Status",
	})

	for _, binding := range bindings {
		t.Append([]string{
			binding.Name,
			binding.Namespace,
			binding.Spec.ServiceInstanceRef.Name,
			getBindingStatusShort(binding.Status),
		})
	}

	t.Render()
}

// WriteBindingDetails prints details for a single binding.
func WriteBindingDetails(binding *v1beta1.ServiceBinding) {
	t := NewDetailsTable()

	t.AppendBulk([][]string{
		{"Name:", binding.Name},
		{"Namespace:", binding.Namespace},
		{"Status:", getBindingStatusFull(binding.Status)},
		{"Instance:", binding.Spec.ServiceInstanceRef.Name},
	})

	t.Render()
}
