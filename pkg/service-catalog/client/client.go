package client

import "github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"

type Client struct {
	*clientset.Clientset
}
