package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/Azure/service-catalog-cli/internal/test"
	"github.com/pkg/errors"
)

var catalogRequestRegex = regexp.MustCompile("/apis/servicecatalog.k8s.io/v1beta1/(.*)")

func TestCommandOutput(t *testing.T) {
	testcases := []struct {
		name   string // Test Name
		cmd    string // Command to run
		golden string // Relative path to a golden file, compared to the command output
	}{
		{"list all brokers", "get brokers", "get-brokers.txt"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := executeCommand(t, tc.cmd)
			test.AssertEqualsGoldenFile(t, tc.golden, output)
		})
	}
}

// executeCommand runs a svcat command against a fake k8s api,
// returning the cli output.
func executeCommand(t *testing.T, cmd string) string {
	// Fake the k8s api server
	apisvr := newAPIServer()
	defer apisvr.Close()

	// Generate a test kubeconfig pointing at the server
	kubeconfig, err := writeTestKubeconfig(apisvr.URL)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.Remove(kubeconfig)

	// Setup the svcat command
	svcat := buildRootCommand()
	args := strings.Split(cmd, " ")
	args = append(args, "--kubeconfig", kubeconfig)
	svcat.SetArgs(args)

	// Capture all output: stderr and stdout
	output := &bytes.Buffer{}
	svcat.SetOutput(output)

	// Ignoring errors, we only care about diffing output
	svcat.Execute()

	return output.String()
}

func newAPIServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(apihandler))
}

// apihandler handles requests to the service catalog endpoint.
// When a request is received, it looks up the response from the testdata directory.
// Example:
// GET /apis/servicecatalog.k8s.io/v1beta1/clusterservicebrokers responds with testdata/clusterservicebrokers.json
func apihandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	match := catalogRequestRegex.FindStringSubmatch(r.RequestURI)

	if len(match) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "unexpected request: "+r.RequestURI)
		return
	}

	if r.Method != http.MethodGet {
		// Anything more interesting than a GET, i.e. it relies upon server behavior
		// probably should be an integration test instead
		w.WriteHeader(200)
		return
	}

	_, response, err := test.GetTestdata(match[1] + ".json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}
	w.Write(response)
}

func writeTestKubeconfig(fakeUrl string) (string, error) {
	_, configT, err := test.GetTestdata("kubeconfig.tmpl.yaml")
	if err != nil {
		return "", err
	}

	data := map[string]string{
		"Server": fakeUrl,
	}
	t := template.Must(template.New("kubeconfig").Parse(string(configT)))

	f, err := ioutil.TempFile("", "kubeconfig")
	if err != nil {
		return "", errors.Wrap(err, "unable to create a temporary kubeconfig file")
	}
	defer f.Close()

	err = t.Execute(f, data)
	return f.Name(), errors.Wrap(err, "error executing the kubeconfig template")
}
