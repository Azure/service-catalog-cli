package environment

import (
	"os"

	"github.com/spf13/pflag"
)

// EnvSettings describes all of the environment settings.
type EnvSettings struct {
	// KubeContext is the name of the kube context.
	KubeContext string
	// KubeConfig is the name of the kubeconfig file.
	KubeConfig string
}

// New adds the appropriate persistent flags, parses them, and negotiates values based on existing environment variables
func New(args []string, flags *pflag.FlagSet) EnvSettings {
	var settings EnvSettings

	settings.addFlags(flags)

	flags.Parse(args)

	// set defaults from environment
	settings.init(flags)

	return settings
}

// addFlags binds flags to the given flagset.
func (s *EnvSettings) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.KubeContext, "kube-context", "", "name of the kube context to use")
	fs.StringVar(&s.KubeConfig, "kubeconfig", "", "path to kubeconfig file. Overrides $KUBECONFIG")
}

// init sets values from the environment.
func (s *EnvSettings) init(fs *pflag.FlagSet) {
	for name, envar := range envMap {
		setFlagFromEnv(name, envar, fs)
	}
}

// envMap maps flag names to envvars
var envMap = map[string]string{
	"kubeconfig": "KUBECONFIG",
}

func setFlagFromEnv(name, envar string, fs *pflag.FlagSet) {
	if fs.Changed(name) {
		return
	}
	if v, ok := os.LookupEnv(envar); ok {
		fs.Set(name, v)
	}
}
