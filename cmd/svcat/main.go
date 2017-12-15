package main

import (
	"fmt"
	"os"

	"github.com/Azure/service-catalog-cli/pkg/binding"
	"github.com/Azure/service-catalog-cli/pkg/broker"
	"github.com/Azure/service-catalog-cli/pkg/class"
	"github.com/Azure/service-catalog-cli/pkg/command"
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
	cxt := &command.Context{}
	env := environment.EnvSettings{}

	// root command flags
	var opts struct {
		Version bool
	}

	cmd := &cobra.Command{
		Use:          "svcat",
		Short:        "The Kubernetes Service Catalog Command-Line Interface (CLI)",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Enable tests to swap the output
			cxt.Output = cmd.OutOrStdout()

			// Initialize a service catalog client
			env.Init()
			_, cl, err := getKubeClient(env)
			if err != nil {
				return fmt.Errorf("Error connecting to Kubernetes (%s)", err)
			}
			cxt.Client = cl

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Version {
				printVersion()
				return nil
			}

			fmt.Print(cmd.UsageString())
			return nil
		},
	}

	cmd.Flags().BoolVarP(&opts.Version, "version", "v", false, "Show the application version")
	env.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(newGetCmd(cxt))
	cmd.AddCommand(newDescribeCmd(cxt))
	cmd.AddCommand(newSyncCmd(cxt))

	return cmd
}

func printVersion() {
	if commit == "" { // commit is empty for Homebrew builds
		fmt.Printf("svcat %s\n", version)
	} else {
		fmt.Printf("svcat %s (%s)\n", version, commit)
	}
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

func newSyncCmd(cxt *command.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "Syncs service catalog for a service broker",
		Aliases: []string{"relist"},
	}
	cmd.AddCommand(broker.NewSyncCmd(cxt))

	return cmd
}

func newGetCmd(cxt *command.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "List a resource, optionally filtered by name",
	}
	cmd.AddCommand(binding.NewGetCmd(cxt))
	cmd.AddCommand(broker.NewGetCmd(cxt))
	cmd.AddCommand(class.NewGetCmd(cxt))
	cmd.AddCommand(instance.NewGetCmd(cxt))
	cmd.AddCommand(plan.NewGetCmd(cxt))

	return cmd
}

func newDescribeCmd(cxt *command.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Show details of a specific resource",
	}
	cmd.AddCommand(binding.NewDescribeCmd(cxt))
	cmd.AddCommand(broker.NewDescribeCmd(cxt))
	cmd.AddCommand(class.NewDescribeCmd(cxt))
	cmd.AddCommand(instance.NewDescribeCmd(cxt))
	cmd.AddCommand(plan.NewDescribeCmd(cxt))

	return cmd
}
