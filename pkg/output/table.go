package output

import (
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

// NewTable returns a new tablewriter with standard options set
func NewTable(indent int) *tablewriter.Table {
	var wr io.Writer = os.Stdout
	if indent > 0 {
		prefix := make([]byte, indent)
		for i := 0; i < indent; i++ {
			prefix[i] = '\t'
		}
		wr = &prependWriter{
			prefix: prefix,
			wr:     os.Stdout,
		}
	}
	t := tablewriter.NewWriter(wr)
	t.SetBorder(false)
	t.SetColumnSeparator("\t")
	t.SetRowLine(true)
	return t
}
