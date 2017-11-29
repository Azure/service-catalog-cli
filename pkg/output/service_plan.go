package output

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

func getPlanStatusShort(status v1beta1.ClusterServicePlanStatus) string {
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

// WritePlanList prints a list of plans.
func WritePlanList(plans []v1beta1.ClusterServicePlan, classes []v1beta1.ClusterServiceClass) {
	classNames := map[string]string{}
	for _, class := range classes {
		classNames[class.Name] = class.Spec.ExternalName
	}

	sort.Sort(byClass(plans))

	t := NewListTable()
	t.SetHeader([]string{
		"Name",
		"Class",
		"Description",
		"UUID"})
	for _, plan := range plans {
		t.Append([]string{
			plan.Spec.ExternalName,
			classNames[plan.Spec.ClusterServiceClassRef.Name],
			plan.Spec.Description,
			plan.Name})
	}
	t.Render()
}

// WriteAssociatedPlans prints a list of plans associated with a class.
func WriteAssociatedPlans(plans []v1beta1.ClusterServicePlan) {
	fmt.Println("\nPlans:")
	if len(plans) == 0 {
		fmt.Println("No plans defined")
		return
	}

	t := NewListTable()
	t.SetHeader([]string{
		"Name",
		"Description",
	})
	for _, plan := range plans {
		t.Append([]string{
			plan.Spec.ExternalName,
			plan.Spec.Description,
		})
	}
	t.Render()
}

// WriteParentPlan prints identifying information for a parent class.
func WriteParentPlan(plan *v1beta1.ClusterServicePlan) {
	fmt.Println("\nPlan:")
	t := NewDetailsTable()
	t.AppendBulk([][]string{
		{"Name:", plan.Spec.ExternalName},
		{"UUID:", string(plan.Name)},
		{"Status:", getPlanStatusShort(plan.Status)},
	})
	t.Render()
}

// WritePlanDetails prints details for a single plan.
func WritePlanDetails(plan *v1beta1.ClusterServicePlan, class *v1beta1.ClusterServiceClass) {
	t := NewDetailsTable()

	t.AppendBulk([][]string{
		{"Name:", plan.Spec.ExternalName},
		{"Description:", plan.Spec.Description},
		{"UUID:", string(plan.Name)},
		{"Class:", class.Spec.ExternalName},
		{"Status:", getPlanStatusShort(plan.Status)},
		{"Free:", strconv.FormatBool(plan.Spec.Free)},
	})

	t.Render()
}
