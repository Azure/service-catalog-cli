package output

import (
	"fmt"
	"io"

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
	return formatStatusShort(string(lastCond.Type), lastCond.Status, lastCond.Reason)
}

func getBindingStatusFull(status v1beta1.ServiceBindingStatus) string {
	lastCond := getBindingStatusCondition(status)
	return formatStatusFull(string(lastCond.Type), lastCond.Status, lastCond.Reason, lastCond.Message, lastCond.LastTransitionTime)
}

// WriteBindingDetails prints a list of bindings.
func WriteBindingList(w io.Writer, bindings ...v1beta1.ServiceBinding) {
	t := NewListTable(w)
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
func WriteBindingDetails(w io.Writer, binding *v1beta1.ServiceBinding) {
	t := NewDetailsTable(w)

	t.AppendBulk([][]string{
		{"Name:", binding.Name},
		{"Namespace:", binding.Namespace},
		{"Status:", getBindingStatusFull(binding.Status)},
		{"Instance:", binding.Spec.ServiceInstanceRef.Name},
	})

	t.Render()
}

// WriteAssociatedBindings prints a list of bindings associated with an instance.
func WriteAssociatedBindings(w io.Writer, bindings []v1beta1.ServiceBinding) {
	fmt.Fprintln(w, "\nBindings:")
	if len(bindings) == 0 {
		fmt.Fprintln(w, "No bindings defined")
		return
	}

	t := NewListTable(w)
	t.SetHeader([]string{
		"Name",
		"Status",
	})
	for _, binding := range bindings {
		t.Append([]string{
			binding.Name,
			getBindingStatusShort(binding.Status),
		})
	}
	t.Render()
}
