package output

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/olekukonko/tablewriter"
)

func getPlanStatusText(status v1beta1.ClusterServicePlanStatus) string {
	if status.RemovedFromBrokerCatalog {
		return statusDeprecated
	}
	return statusActive
}

// ByAge implements sort.Interface for []Person based on
// the Age field.
type byClass []v1beta1.ClusterServicePlan

func (a byClass) Len() int {
	return len(a)
}
func (a byClass) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a byClass) Less(i, j int) bool {
	return a[i].Spec.ClusterServiceClassRef.Name < a[j].Spec.ClusterServiceClassRef.Name
}

// ClusterServicePlanHeaders sets the appropriate headers on t for displaying
// ClusterServicePlans in t
func ClusterServicePlanHeaders(table *tablewriter.Table) {
	table.SetHeader([]string{
		"Name",
		"Description",
		"UUID",
		"Class Name",
		"Class UUID",
	})
}

// AppendClusterServicePlan appends plan to t by calling t.Append.
// Ensure that you've called ClusterServicePlanHeaders on t before you call this function
func AppendClusterServicePlan(table *tablewriter.Table, plan *v1beta1.ClusterServicePlan) {
	table.Append([]string{
		plan.Spec.ExternalName,
		plan.Spec.Description,
		plan.Name,
		plan.Spec.ClusterServiceBrokerName,
		plan.Spec.ClusterServiceClassRef.Name,
	})
}

// WritePlanList prints a list of plans.
func WritePlanList(plans *v1beta1.ClusterServicePlanList, classes *v1beta1.ClusterServiceClassList) {
	classNames := map[string]string{}
	for _, class := range classes.Items {
		classNames[class.Name] = class.Spec.ExternalName
	}

	sort.Sort(byClass(plans.Items))

	t := NewListTable()
	t.SetHeader([]string{
		"Name",
		"Class",
		"Description",
		"UUID"})
	for _, plan := range plans.Items {
		t.Append([]string{
			plan.Spec.ExternalName,
			classNames[plan.Spec.ClusterServiceClassRef.Name],
			plan.Spec.Description,
			plan.Name})
	}
	t.Render()
}

// WriteAssociatedPlans prints a list of plans associated with a class.
func WriteAssociatedPlans(plans *v1beta1.ClusterServicePlanList) {
	fmt.Println("\nPlans:")
	if len(plans.Items) == 0 {
		fmt.Println("No plans defined")
		return
	}

	t := NewListTable()
	t.SetHeader([]string{
		"Name",
		"Description",
	})
	for _, plan := range plans.Items {
		t.Append([]string{
			plan.Spec.ExternalName,
			plan.Spec.Description,
		})
	}
	t.Render()
}

// WritePlanDetails prints a service plan to the console.
func WritePlanDetails(plan *v1beta1.ClusterServicePlan, class *v1beta1.ClusterServiceClass) {
	t := NewDetailsTable()

	t.AppendBulk([][]string{
		{"Name:", plan.Spec.ExternalName},
		{"Description:", plan.Spec.Description},
		{"UUID:", string(plan.Name)},
		{"Class:", class.Spec.ExternalName},
		{"Status:", getPlanStatusText(plan.Status)},
		{"Free:", strconv.FormatBool(plan.Spec.Free)},
	})

	t.Render()
}
