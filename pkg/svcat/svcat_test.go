package svcat

import (
	"testing"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

func GetTestApp(t *testing.T) *App {
	t.Helper()

	svcat, err := NewApp("", "")
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
