package output

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/olekukonko/tablewriter"
)

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
