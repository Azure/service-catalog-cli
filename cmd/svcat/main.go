package main

import (
	"fmt"
	"os"

	"github.com/Azure/service-catalog-cli/pkg/binding"
	"github.com/Azure/service-catalog-cli/pkg/broker"
	"github.com/Azure/service-catalog-cli/pkg/class"
	"github.com/Azure/service-catalog-cli/pkg/environment"
	"github.com/Azure/service-catalog-cli/pkg/instance"
	"github.com/Azure/service-catalog-cli/pkg/kube"
	"github.com/Azure/service-catalog-cli/pkg/plan"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
)

// These are build-time values, set during an official release
var (
	commit  string
	version string
)

func main() {
	cmd := buildRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func buildRootCommand() *cobra.Command {
	// root command context
	var opts struct {
		Version bool
	}

	cmd := &cobra.Command{
		Use:          "svcat",
		Short:        "The Kubernetes Service Catalog Command-Line Interface (CLI)",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Version {
				if commit == "" { // commit is empty for Homebrew builds
					fmt.Printf("svcat %s\n", version)
				} else {
					fmt.Printf("svcat %s (%s)\n", version, commit)
				}
				return nil
			}

			fmt.Print(cmd.UsageString())
			return nil
		},
	}

	cmd.Flags().BoolVarP(&opts.Version, "version", "v", false, "Show the application version")

	flags := cmd.PersistentFlags()

	// adds the appropriate persistent flags, parses them, and negotiates values based on existing environment variables
	vars := environment.New(os.Args, flags)

	_, cl, err := getKubeClient(vars)
	if err != nil {
		logger.Fatalf("Error connecting to Kubernetes (%s)", err)
	}

	cmd.AddCommand(newGetCmd(cl))
	cmd.AddCommand(newDescribeCmd(cl))
	cmd.AddCommand(newSyncCmd(cl))

	return cmd
}

// configForContext creates a Kubernetes REST client configuration for a given kubeconfig context.
func configForContext(vars environment.EnvSettings) (*rest.Config, error) {
	config, err := kube.GetConfig(vars.KubeContext, vars.KubeConfig).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", vars.KubeContext, err)
	}
	return config, nil
}

// getKubeClient creates a Kubernetes config and client for a given kubeconfig context.
func getKubeClient(vars environment.EnvSettings) (*rest.Config, *clientset.Clientset, error) {
	config, err := configForContext(vars)
	if err != nil {
		logger.Fatalf("Error getting Kubernetes configuration (%s)", err)
	}
	client, err := clientset.NewForConfig(config)
	return nil, client, err
}

func newSyncCmd(cl *clientset.Clientset) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "Syncs service catalog for a service broker",
		Aliases: []string{"relist"},
	}
	cmd.AddCommand(broker.NewSyncCmd(cl))

	return cmd
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

func newDescribeCmd(cl *clientset.Clientset) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Show details of a specific resource",
	}
	cmd.AddCommand(binding.NewDescribeCmd(cl))
	cmd.AddCommand(broker.NewDescribeCmd(cl))
	cmd.AddCommand(class.NewDescribeCmd(cl))
	cmd.AddCommand(instance.NewDescribeCmd(cl))
	cmd.AddCommand(plan.NewDescribeCmd(cl))

	return cmd
}
