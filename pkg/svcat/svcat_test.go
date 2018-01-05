package svcat

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/service-catalog-cli/pkg/environment"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

var (
	Settings = TestSettings{}
)

type TestSettings struct {
	environment.EnvSettings
	Namespace string
}

// Verify that we can connect to the test cluster before starting any tests
func (s TestSettings) Verify() *App {
	svcat, err := NewApp(Settings.KubeConfig, Settings.KubeContext)
	if err != nil {
		fmt.Printf("%s\nSettings: %+v\n", err, Settings.EnvSettings)
		os.Exit(1)
	}

	return svcat
}

func TestMain(m *testing.M) {
	// Load overrides to the cluster connection from environment variables
	Settings.Init()
	svcat := Settings.Verify()

	// Setup
	Settings.Namespace = CreateTestNamespace(svcat)

	retCode := m.Run()

	// Teardown
	//DeleteTestNamespace(svcat)

	os.Exit(retCode)
}

func NewTestApp(t *testing.T) *App {
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

func CreateTestNamespace(svcat *App) string {
	name := fmt.Sprintf("test-%d", time.Now().Unix())
	fmt.Printf("Creating test namespace: %s\n", name)
	namespace, err := svcat.Kube.CreateNamespace(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return namespace.Name
}

func DeleteTestNamespace(svcat *App) {
	fmt.Printf("Cleaning up test namespace: %s\n", Settings.Namespace)
	err := svcat.Kube.DeleteNamespace(Settings.Namespace)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
