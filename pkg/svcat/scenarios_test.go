package svcat

import (
	"fmt"
	"testing"
	"time"

	"github.com/Azure/service-catalog-cli/internal/test"
)

func TestIntegrationLifecycle(t *testing.T) {
	test.FilterTestSuite(t)
	//t.Parallel()

	svcat := NewTestApp(t)

	const instanceName = "mysql"

	waitForOp := func() error {
		t.Log("Waiting for the instance operation to complete...")
		for {
			time.Sleep(time.Second * 30)

			instance, err := svcat.RetrieveInstance(Settings.Namespace, instanceName)
			if err != nil {
				return err
			}

			if !instance.Status.AsyncOpInProgress {
				return nil
			}
			t.Logf("... %s\n", instance.Status.Conditions[len(instance.Status.Conditions)-1].Reason)
		}
	}

	t.Log("Provisioning an instance...")
	instance, err := svcat.Provision(Settings.Namespace, instanceName, "azure-mysqldb", "basic50",
		map[string]string{
			"location":               "eastus",
			"resourceGroup":          "carolynvs-" + Settings.Namespace,
			"sslEnforcement":         "disabled",
			"firewallStartIPAddress": "0.0.0.0",
			"firewallEndIPAddress":   "255.255.255.255",
		}, nil)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if instanceName != instance.Name {
		t.Fatalf("expected an instance named %s but got %s", instanceName, instance.Name)
	}

	err = waitForOp()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// bind the instance multiple times
	numBindings := 3
	for i := 1; i <= numBindings; i++ {
		t.Log("Adding binding...")
		bindingName := fmt.Sprintf("%s-binding-%d", instanceName, i)
		secretName := fmt.Sprintf("%s-secret-%d", instanceName, i)
		binding, err := svcat.Bind(Settings.Namespace, bindingName, instanceName, secretName, nil, nil)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if bindingName != binding.Name {
			t.Fatalf("expected a binding named %s but got %s", bindingName, binding.Name)
		}
	}
	bindings, err := svcat.RetrieveBindingsByInstance(Settings.Namespace, instanceName)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if len(bindings) != numBindings {
		t.Fatalf("unexpected number of bindings created, want: %d, got: %d", numBindings, len(bindings))
	}

	// remove a particular binding
	t.Log("Removing a single binding...")
	err = svcat.DeleteBinding(Settings.Namespace, bindings[0].Name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	bindings, err = svcat.RetrieveBindingsByInstance(Settings.Namespace, instanceName)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if len(bindings) != numBindings-1 {
		t.Fatalf("unexpected number of bindings remaining, want: %d, got: %d", numBindings-1, len(bindings))
	}

	// Unbind the remaining bindings
	t.Log("Unbinding the instance...")
	err = svcat.Unbind(Settings.Namespace, instanceName)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	bindings, err = svcat.RetrieveBindingsByInstance(Settings.Namespace, instanceName)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if len(bindings) != 0 {
		t.Fatalf("unexpected number of bindings remaining, want: %s, got: %s", 0, len(bindings))
	}

	// Remove the instance
	t.Log("Deprovisioning the instance...")
	err = svcat.Deprovision(Settings.Namespace, instanceName)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	err = waitForOp()
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
