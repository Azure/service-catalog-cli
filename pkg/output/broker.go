package output

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/olekukonko/tablewriter"
)

func getBrokerStatusCondition(status v1beta1.ClusterServiceBrokerStatus) v1beta1.ServiceBrokerCondition {
	if len(status.Conditions) > 0 {
		return status.Conditions[len(status.Conditions)-1]
	}
	return v1beta1.ServiceBrokerCondition{}
}

func getBrokerStatusShort(status v1beta1.ClusterServiceBrokerStatus) string {
	lastCond := getBrokerStatusCondition(status)
	return string(lastCond.Type)
}

func getBrokerStatusFull(status v1beta1.ClusterServiceBrokerStatus) string {
	lastCond := getBrokerStatusCondition(status)
	return formatStatusText(string(lastCond.Type), lastCond.Message, lastCond.LastTransitionTime)
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

// WriteParentBroker prints identifying information for a parent broker.
func WriteParentBroker(broker *v1beta1.ClusterServiceBroker) {
	fmt.Println("\nBroker:")
	t := NewDetailsTable()
	t.AppendBulk([][]string{
		{"Name:", broker.Name},
		{"Status:", getBrokerStatusShort(broker.Status)},
	})
	t.Render()
}

// WriteBrokerDetails prints details for a single broker.
func WriteBrokerDetails(broker *v1beta1.ClusterServiceBroker) {
	t := NewDetailsTable()

	t.AppendBulk([][]string{
		{"Name:", broker.Name},
		{"URL:", broker.Spec.URL},
		{"Status:", getBrokerStatusFull(broker.Status)},
	})

	t.Render()
}
