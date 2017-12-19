package params

import (
	"fmt"
	"strings"
)

// ParseVariableAssignment converts a string array of variable assignments
// into a map of keys and values
// Example:
// [a=b c=abc1232===] becomes map[a:b c:abc1232===]
func ParseVariableAssignments(params []string) (map[string]string, error) {
	variables := map[string]string{}

	for _, p := range params {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid parameter (%s), must be in name=value format", p)
		}

		variable := strings.TrimSpace(parts[0])
		if variable == "" {
			return nil, fmt.Errorf("invalid parameter (%s), variable name is requried", p)
		}
		value := strings.TrimSpace(parts[1])

		variables[variable] = value
	}

	return variables, nil
}
