package svcat

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/service-catalog-cli/pkg/environment"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

var (
	Settings = TestSettings{}
)

type TestSettings struct {
	environment.EnvSettings
}

// Verify that we can connect to the test cluster before starting any tests
func (s TestSettings) Verify() {
	_, err := NewApp(Settings.KubeConfig, Settings.KubeContext)
	if err != nil {
		fmt.Printf("%s\nSettings: %+v\n", err, Settings.EnvSettings)
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	// Load overrides to the cluster connection from environment variables
	Settings.Init()
	Settings.Verify()

	os.Exit(m.Run())
}

func GetTestApp(t *testing.T) *App {
	t.Helper()

	svcat, err := NewApp(Settings.KubeConfig, Settings.KubeContext)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	return svcat
}

func GetTestBroker(t *testing.T, svcat *App) *v1beta1.ClusterServiceBroker {
	t.Helper()

	brokers, err := svcat.RetrieveBrokers()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if len(brokers) == 0 {
		t.Fatalf("no test broker found")
	}

	return &brokers[0]
}
