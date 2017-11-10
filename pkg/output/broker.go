package output

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/olekukonko/tablewriter"
)

// ClusterServiceBrokerHeaders sets the appropriate headers on t for displaying
// ClusterServiceBrokers in t
func ClusterServiceBrokerHeaders(t *tablewriter.Table) {
	t.SetHeader([]string{
		"Type",
		"Status",
		"Reason",
		"Message",
	})
}

// AppendClusterServiceBroker appends instance to t by calling t.Append.
// Ensure that you've called ClusterServiceBrokerHeaders on t before you call
// this function
func AppendClusterServiceBroker(t *tablewriter.Table, broker *v1beta1.ClusterServiceBroker) {
	t.SetHeader([]string{
		"Type",
		"Status",
		"Reason",
		"Message",
	})
}
