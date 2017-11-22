package main

import (
	"os"

	"github.com/Azure/service-catalog-cli/pkg/binding"
	"github.com/Azure/service-catalog-cli/pkg/broker"
	"github.com/Azure/service-catalog-cli/pkg/class"
	"github.com/Azure/service-catalog-cli/pkg/instance"
	"github.com/Azure/service-catalog-cli/pkg/kube"
	"github.com/Azure/service-catalog-cli/pkg/plan"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:          "svc-cat",
		Short:        "The Kubernetes Service-Catalog Command Line Interface (CLI)",
		SilenceUsage: true,
	}

	flags := cmd.PersistentFlags()
	kubeConfigLocation := flags.String(
		"config",
		kube.DefaultConfigLocation(),
		"the location of the Kubernetes configuration file",
	)
	kubeContext := flags.String(
		"context",
		"",
		"the context to use in the Kubernetes configuration file",
	)
	flags.Parse(os.Args)

	cfg, err := kube.ConfigForContext(*kubeConfigLocation, *kubeContext)
	if err != nil {
		logger.Fatalf("Error getting Kubernetes configuration (%s)", err)
	}
	cl, err := clientset.NewForConfig(cfg)
	if err != nil {
		logger.Fatalf("Error connecting to Kubernetes (%s)", err)
	}

	cmd.AddCommand(newGetCmd(cl))

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newGetCmd(cl *clientset.Clientset) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "List a resource, optionally filtered by name",
	}
	cmd.AddCommand(binding.NewGetCmd(cl))
	cmd.AddCommand(broker.NewGetCmd(cl))
	cmd.AddCommand(class.NewGetCmd(cl))
	cmd.AddCommand(instance.NewGetCmd(cl))
	cmd.AddCommand(plan.NewGetCmd(cl))

	return cmd
}
