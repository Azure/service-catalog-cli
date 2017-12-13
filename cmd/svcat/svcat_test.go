package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Azure/service-catalog-cli/internal/test"
)

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

			svcat := buildRootCommand()
			svcat.SetArgs(strings.Split(tc.cmd, " "))

			// Capture all output: stderr and stdout
			output := &bytes.Buffer{}
			svcat.SetOutput(output)

			// Ignoring errors, we only care about diffing output
			svcat.Execute()

			test.AssertEqualsGoldenFile(t, tc.golden, output.Bytes())
		})
	}
}
