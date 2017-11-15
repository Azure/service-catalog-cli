package kube

import (
	"fmt"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// DefaultConfigLocation returns the default location of the Kubernetes
// config file.
func DefaultConfigLocation() string {
	return fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
}

// ConfigForContext gets the Kubernetes REST configuration, for the given context,
// from the default Kubernetes configuration file on disk
func ConfigForContext(configLocation string, context string) (*rest.Config, error) {
	config, err := getConfig(context, configLocation).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", context, err)
	}
	return config, nil
}

// GetConfig returns a Kubernetes client config for a given context.
func getConfig(context string, kubeconfig string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}

	if context != "" {
		overrides.CurrentContext = context
	}

	if kubeconfig != "" {
		rules.ExplicitPath = kubeconfig
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}
