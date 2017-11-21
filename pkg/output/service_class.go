package output

import (
	"strings"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/olekukonko/tablewriter"
)

func getStatusText(status v1beta1.ClusterServiceClassStatus) string {
	if status.RemovedFromBrokerCatalog {
		return "deprecated"
	}
	return "active"
}

// ClusterServiceClassHeaders sets the appropriate headers on t for displaying
// ClusterServiceClasses in t
func ClusterServiceClassHeaders(table *tablewriter.Table) {
	table.SetHeader([]string{
		"Name",
		"Description",
		"UUID",
	})
}

// AppendClusterServiceClass appends class to t by calling t.Append.
// Ensure that you've called ClusterServiceClassHeaders on t before you call this function
func AppendClusterServiceClass(table *tablewriter.Table, class *v1beta1.ClusterServiceClass) {
	table.Append([]string{
		class.Spec.ExternalName,
		class.Spec.Description,
		class.Name,
	})
}

// WriteClusterServiceClassList prints a list of service class to the console.
func WriteClusterServiceClassList(classes []v1beta1.ClusterServiceClass) {
	t := NewListTable()
	t.SetHeader([]string{
		"Name",
		"Description",
		"UUID",
	})
	for _, class := range classes {
		t.Append([]string{
			class.Spec.ExternalName,
			class.Spec.Description,
			class.Name,
		})
	}
	t.Render()
}

// WriteClusterServiceClassDetails prints a service class to the console.
func WriteClusterServiceClass(class *v1beta1.ClusterServiceClass) {
	t := NewDetailsTable()
	t.AppendBulk([][]string{
		{"Name:", class.Spec.ExternalName},
		{"Description:", class.Spec.Description},
		{"UUID:", string(class.Name)},
		{"Status:", getStatusText(class.Status)},
		{"Tags:", strings.Join(class.Spec.Tags, ", ")},
	})
	t.Render()
}
