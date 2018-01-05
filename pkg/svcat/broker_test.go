package svcat

import (
	"testing"

	"github.com/Azure/service-catalog-cli/internal/test"
)

func TestIntegrationBrokerGetAll(t *testing.T) {
	test.FilterTestSuite(t)
	t.Parallel()

	svcat := NewTestApp(t)

	brokers, err := svcat.RetrieveBrokers()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if len(brokers) == 0 {
		t.Fatalf("expected at least one broker")
	}
}

func TestIntegrationBrokerGet(t *testing.T) {
	test.FilterTestSuite(t)
	t.Parallel()

	svcat := NewTestApp(t)
	brokerName := GetTestBroker(t, svcat).Name

	broker, err := svcat.RetrieveBroker(brokerName)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if broker.Name != brokerName {
		t.Fatalf("wrong broker retrieved, expected: %s, got: %s", brokerName, broker.Name)
	}
}

func TestIntegrationBrokerSync(t *testing.T) {
	test.FilterTestSuite(t)
	t.Parallel()

	svcat := NewTestApp(t)
	brokerName := GetTestBroker(t, svcat).Name

	err := svcat.Sync(brokerName, 0)
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
