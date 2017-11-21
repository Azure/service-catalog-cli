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

// NewListTable builds a table formatted to list a set of results.
func NewListTable() *tablewriter.Table {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetBorder(false)
	t.SetColumnSeparator("\t")
	t.SetRowLine(false)
	return t
}

// NewDetailsTable builds a table formatted to list details for a single result.
func NewDetailsTable() *tablewriter.Table {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetBorder(false)
	t.SetColumnSeparator("\t")
	t.SetRowLine(false)

	// tablewriter wraps based on "ragged text", not max column width
	// which is great for tables but isn't efficient for detailed views
	t.SetAutoWrapText(false)

	return t
}
