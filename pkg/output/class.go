package output

import (
	"fmt"
	"strings"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

func getClassStatusText(status v1beta1.ClusterServiceClassStatus) string {
	if status.RemovedFromBrokerCatalog {
		return statusDeprecated
	}
	return statusActive
}

// WriteClassList prints a list of classes.
func WriteClassList(classes ...v1beta1.ClusterServiceClass) {
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

// WriteParentClass prints identifying information for a parent class.
func WriteParentClass(class *v1beta1.ClusterServiceClass) {
	fmt.Println("\nClass:")
	t := NewDetailsTable()
	t.AppendBulk([][]string{
		{"Name:", class.Spec.ExternalName},
		{"UUID:", string(class.Name)},
		{"Status:", getClassStatusText(class.Status)},
	})
	t.Render()
}

// WriteClassDetails prints details for a single class.
func WriteClassDetails(class *v1beta1.ClusterServiceClass) {
	t := NewDetailsTable()
	t.AppendBulk([][]string{
		{"Name:", class.Spec.ExternalName},
		{"Description:", class.Spec.Description},
		{"UUID:", string(class.Name)},
		{"Status:", getClassStatusText(class.Status)},
		{"Tags:", strings.Join(class.Spec.Tags, ", ")},
		{"Broker:", class.Spec.ClusterServiceBrokerName},
	})
	t.Render()
}
