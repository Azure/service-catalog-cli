package output

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/olekukonko/tablewriter"
)

func getBindingStatusText(status v1beta1.ServiceBindingStatus) string {
	if len(status.Conditions) > 0 {
		lastCond := status.Conditions[len(status.Conditions)-1]
		return formatStatusText(string(lastCond.Type), lastCond.Message, lastCond.LastTransitionTime)
	}
	return statusNone
}

// BindingHeaders sets the headers for listing bindings in t
func BindingHeaders(t *tablewriter.Table) {
	t.SetHeader([]string{
		"Name",
		"Service Instance Name",
		"Status",
	})
}

// AppendBinding appends binding to t by calling t.Append.
// Ensure that you've called BindingHeaders on t before you call AppendBinding
func AppendBinding(t *tablewriter.Table, binding *v1beta1.ServiceBinding) {
	lastCond := ""
	if len(binding.Status.Conditions) > 0 {
		lastCond = binding.Status.Conditions[len(binding.Status.Conditions)-1].Reason
	}
	t.Append([]string{
		binding.Name,
		binding.Spec.ServiceInstanceRef.Name,
		lastCond,
	})
}

// WriteBindingDetails prints a binding.
func WriteBindingDetails(binding *v1beta1.ServiceBinding) {
	t := NewDetailsTable()

	t.AppendBulk([][]string{
		{"Name:", binding.Name},
		{"Namespace:", binding.Namespace},
		{"Status:", getBindingStatusText(binding.Status)},
		{"Instance:", binding.Spec.ServiceInstanceRef.Name},
	})

	t.Render()
}
