package output

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	statusNone       = ""
	statusActive     = "Active"
	statusDeprecated = "Deprecated"
)

func formatStatusText(status, message string, timestamp v1.Time) string {
	message = strings.TrimRight(message, ".")
	return fmt.Sprintf("%s - %s @ %s", status, message, timestamp)
}
