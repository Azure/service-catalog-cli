package output

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

// NewTable returns a new tablewriter with standard options set
func NewTable() *tablewriter.Table {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetBorder(false)
	t.SetColumnSeparator("\t")
	t.SetRowLine(true)
	return t
}
