package main

import (
	"fmt"
	"os"

	"github.com/Azure/service-catalog-cli/cmd/svcat/binding"
	"github.com/Azure/service-catalog-cli/cmd/svcat/broker"
	"github.com/Azure/service-catalog-cli/cmd/svcat/class"
	"github.com/Azure/service-catalog-cli/cmd/svcat/command"
	"github.com/Azure/service-catalog-cli/cmd/svcat/instance"
	"github.com/Azure/service-catalog-cli/cmd/svcat/plan"
	"github.com/Azure/service-catalog-cli/pkg/environment"
	"github.com/Azure/service-catalog-cli/pkg/svcat"
	"github.com/spf13/cobra"
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

			// Initialize flags from environment variables
			env.Init()

			app, err := svcat.NewApp(env.KubeConfig, env.KubeContext)
			cxt.App = app

			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Version {
				printVersion(cxt)
				return nil
			}

			fmt.Fprint(cxt.Output, cmd.UsageString())
			return nil
		},
	}

	cmd.Flags().BoolVarP(&opts.Version, "version", "v", false, "Show the application version")
	env.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(newGetCmd(cxt))
	cmd.AddCommand(newDescribeCmd(cxt))
	cmd.AddCommand(instance.NewProvisionCmd(cxt))
	cmd.AddCommand(instance.NewDeprovisionCmd(cxt))
	cmd.AddCommand(binding.NewBindCmd(cxt))
	cmd.AddCommand(binding.NewUnbindCmd(cxt))
	cmd.AddCommand(newSyncCmd(cxt))

	return cmd
}

func printVersion(cxt *command.Context) {
	if commit == "" { // commit is empty for Homebrew builds
		fmt.Fprintf(cxt.Output, "svcat %s\n", version)
	} else {
		fmt.Fprintf(cxt.Output, "svcat %s (%s)\n", version, commit)
	}
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
