package output

import (
	"fmt"
	"strings"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/olekukonko/tablewriter"
)

// ClusterServiceClassHeaders sets the appropriate headers on t for displaying
// ClusterServiceClasses in t
func ClusterServiceClassHeaders(table *tablewriter.Table) {
	table.SetHeader([]string{
		"Name",
		"Description",
		"UUID",
		"Bindable",
		"Tags",
	})
}

// AppendClusterServiceClass appends class to t by calling t.Append.
// Ensure that you've called ClusterServiceClassHeaders on t before you call this function
func AppendClusterServiceClass(table *tablewriter.Table, class *v1beta1.ClusterServiceClass) {
	table.Append([]string{
		class.Spec.ExternalName,
		class.Spec.Description,
		class.Name,
		fmt.Sprintf("%t", class.Spec.Bindable),
		strings.Join(class.Spec.Tags, ", "),
	})
}
