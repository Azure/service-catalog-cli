package svcat

import (
	"testing"

	"github.com/Azure/service-catalog-cli/internal/test"
)

func TestIntegrationBrokerGetAll(t *testing.T) {
	test.FilterTestSuite(t)

	svcat, err := NewApp("", "")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	brokers, err := svcat.RetrieveBrokers()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if len(brokers) == 0 {
		t.Fatalf("expected at least one broker")
	}
}
