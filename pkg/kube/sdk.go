package kube

import (
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// SDK wrapper around the generated Go client for the Kubernetes API.
type SDK struct {
	KubeClient *kubernetes.Clientset
}

// Core is the underlying generated Kubernetes Core versioned interface.
// It should be used instead of accessing the client directly.
func (sdk *SDK) Core() corev1.CoreV1Interface {
	return sdk.KubeClient.CoreV1()
}
