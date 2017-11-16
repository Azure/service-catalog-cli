package environment

import (
	"os"

	"github.com/spf13/pflag"
)

// EnvSettings describes all of the environment settings.
type EnvSettings struct {
	// KubeContext is the name of the kubeconfig context.
	KubeContext string
	// KubeConfig is the name of the kubeconfig file.
	KubeConfig string
}

// AddFlags binds flags to the given flagset.
func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.KubeContext, "kube-context", "", "name of the kubeconfig context to use")
	fs.StringVar(&s.KubeConfig, "kubeconfig", "", "path to kubeconfig file. Overrides $KUBECONFIG")
}

// Init sets values from the environment.
func (s *EnvSettings) Init(fs *pflag.FlagSet) {
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
