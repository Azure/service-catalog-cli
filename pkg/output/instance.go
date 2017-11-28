package output

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/olekukonko/tablewriter"
)

func getInstanceStatusCondition(status v1beta1.ServiceInstanceStatus) v1beta1.ServiceInstanceCondition {
	if len(status.Conditions) > 0 {
		return status.Conditions[len(status.Conditions)-1]
	}
	return v1beta1.ServiceInstanceCondition{}
}

func getInstanceStatusFull(status v1beta1.ServiceInstanceStatus) string {
	lastCond := getInstanceStatusCondition(status)
	return formatStatusText(string(lastCond.Type), lastCond.Message, lastCond.LastTransitionTime)
}

func getInstanceStatusShort(status v1beta1.ServiceInstanceStatus) string {
	lastCond := getInstanceStatusCondition(status)
	return string(lastCond.Type)
}

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

// WriteParentInstance prints identifying information for a parent instance.
func WriteParentInstance(instance *v1beta1.ServiceInstance) {
	fmt.Println("\nInstance:")
	t := NewDetailsTable()
	t.AppendBulk([][]string{
		{"Name:", instance.Name},
		{"Namespace:", instance.Namespace},
		{"Status:", getInstanceStatusShort(instance.Status)},
	})
	t.Render()
}

// WriteAssociatedInstances prints a list of instances associated with a service plan.
func WriteAssociatedInstances(instances *v1beta1.ServiceInstanceList) {
	fmt.Println("\nInstances:")
	if len(instances.Items) == 0 {
		fmt.Println("No instances defined")
		return
	}

	t := NewListTable()
	t.SetHeader([]string{
		"Name",
		"Namespace",
		"Status",
	})
	for _, instance := range instances.Items {
		t.Append([]string{
			instance.Name,
			instance.Namespace,
			getInstanceStatusFull(instance.Status),
		})
	}
	t.Render()
}

// WriteInstanceDetails prints an instance.
func WriteInstanceDetails(instance *v1beta1.ServiceInstance) {
	t := NewDetailsTable()

	t.AppendBulk([][]string{
		{"Name:", instance.Name},
		{"Namespace:", instance.Namespace},
		{"Status:", getInstanceStatusFull(instance.Status)},
		{"Class:", instance.Spec.ClusterServiceClassExternalName},
		{"Plan:", instance.Spec.ClusterServicePlanExternalName},
	})

	t.Render()
}
