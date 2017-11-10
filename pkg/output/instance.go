package output

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/olekukonko/tablewriter"
)

// InstanceHeaders sets the appropriate headers on t for displaying ServiceInstances
// in t
func InstanceHeaders(t *tablewriter.Table) {
	t.SetHeader([]string{
		"Name",
		"Service Class Name",
		"Service Class UUID",
		"Service Plan Name",
		"Service Plan UUID",
		"Status",
	})
}

// AppendInstance appends instance to t by calling t.Append.
// Ensure that you've called InstanceHeaders on t before you call AppendInstance
func AppendInstance(t *tablewriter.Table, instance *v1beta1.ServiceInstance) {
	latestCond := "None"
	if len(instance.Status.Conditions) >= 1 {
		latestCond = instance.Status.Conditions[len(instance.Status.Conditions)-1].Reason
	}
	t.Append([]string{
		instance.Name,
		instance.Spec.ClusterServiceClassExternalName,
		instance.Spec.ClusterServiceClassRef.Name,
		instance.Spec.ClusterServicePlanExternalName,
		instance.Spec.ClusterServicePlanRef.Name,
		latestCond,
	})
}
