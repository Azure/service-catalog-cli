package output

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

// NewListTable builds a table formatted to list a set of results.
func NewListTable() *tablewriter.Table {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetBorder(false)
	t.SetColumnSeparator(" ")
	return t
}

// NewDetailsTable builds a table formatted to list details for a single result.
func NewDetailsTable() *tablewriter.Table {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetBorder(false)
	t.SetColumnSeparator(" ")

	// tablewriter wraps based on "ragged text", not max column width
	// which is great for tables but isn't efficient for detailed views
	t.SetAutoWrapText(false)

	return t
}
