package output

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/olekukonko/tablewriter"
)

func getBrokerStatusText(status v1beta1.ClusterServiceBrokerStatus) string {
	if len(status.Conditions) > 0 {
		lastCond := status.Conditions[len(status.Conditions)-1]
		return formatStatusText(string(lastCond.Type), lastCond.Message, lastCond.LastTransitionTime)
	}
	return statusNone
}

// ClusterServiceBrokerHeaders sets the appropriate headers on t for displaying
// ClusterServiceBrokers in t
func ClusterServiceBrokerHeaders(t *tablewriter.Table) {
	t.SetHeader([]string{
		"Name",
		"URL",
		"Status Type",
		"Status Value",
		"Status Reason",
		"Status Message",
	})
}

// AppendClusterServiceBroker appends instance to t by calling t.Append.
// Ensure that you've called ClusterServiceBrokerHeaders on t before you call
// this function
func AppendClusterServiceBroker(t *tablewriter.Table, broker *v1beta1.ClusterServiceBroker) {
	condType := ""
	condStatus := ""
	condReason := ""
	condMessage := ""
	if len(broker.Status.Conditions) >= 1 {
		cond := &broker.Status.Conditions[len(broker.Status.Conditions)-1]
		condType = string(cond.Type)
		condStatus = string(cond.Status)
		condReason = cond.Reason
		condMessage = cond.Message
	}
	t.Append([]string{
		broker.Name,
		broker.Spec.URL,
		condType,
		condStatus,
		condReason,
		condMessage,
	})
}

// WriteBrokerDetails prints a broker to the console.
func WriteBrokerDetails(broker *v1beta1.ClusterServiceBroker) {
	t := NewDetailsTable()

	t.AppendBulk([][]string{
		{"Name:", broker.Name},
		{"URL:", broker.Spec.URL},
		{"Status:", getBrokerStatusText(broker.Status)},
	})

	t.Render()
}
