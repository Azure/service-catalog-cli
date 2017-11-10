package main

import (
	"os"

	"github.com/Azure/service-catalog-cli/pkg/binding"
	"github.com/Azure/service-catalog-cli/pkg/broker"
	"github.com/Azure/service-catalog-cli/pkg/catalog"
	"github.com/Azure/service-catalog-cli/pkg/instance"
	"github.com/Azure/service-catalog-cli/pkg/kube"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:          "svc-cat",
		Short:        "The Kubernetes Service-Catalog Command Line Interface (CLI)",
		SilenceUsage: true,
	}

	kubeContext := "" // TODO: make this configurable
	cfg, err := kube.ConfigForContext(kubeContext)
	if err != nil {
		logger.Fatalf("Error getting Kubernetes configuration (%s)", err)
	}
	cl, err := clientset.NewForConfig(cfg)
	if err != nil {
		logger.Fatalf("Error connecting to Kubernetes (%s)", err)
	}
	cmd.AddCommand(broker.NewRootCmd(cl))
	cmd.AddCommand(catalog.NewRootCmd(cl))
	cmd.AddCommand(instance.NewRootCmd(cl))
	cmd.AddCommand(binding.NewRootCmd(cl))

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

}
