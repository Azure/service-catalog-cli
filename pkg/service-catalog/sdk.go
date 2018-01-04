package servicecatalog

import "github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"

// SDK wrapper around the generated Go client for the Kubernetes Service Catalog
type SDK struct {
	*clientset.Clientset
}
