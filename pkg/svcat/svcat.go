package svcat

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/kube"
	"github.com/Azure/service-catalog-cli/pkg/service-catalog"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// App is the underlying application behind the svcat cli.
type App struct {
	*servicecatalog.SDK
	Kube *kube.SDK
}

// NewApp creates an svcat application.
func NewApp(kubeConfig, kubeContext string) (*App, error) {
	cfg, err := configForContext(kubeConfig, kubeContext)
	if err != nil {
		return nil, err
	}

	kubecl, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create a kubernetes client (%s)", err)
	}

	svcatcl, err := clientset.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create a service catalog client (%s)", err)
	}

	app := &App{
		SDK:  &servicecatalog.SDK{ServiceCatalogClient: svcatcl},
		Kube: &kube.SDK{KubeClient: kubecl},
	}

	return app, nil
}

// configForContext creates a Kubernetes REST client configuration for a given kubeconfig context.
func configForContext(kubeConfig, kubeContext string) (*rest.Config, error) {
	config, err := kube.GetConfig(kubeContext, kubeConfig).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", kubeContext, err)
	}
	return config, nil
}
