package output

import (
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

// NewListTable builds a table formatted to list a set of results.
func NewListTable2(w io.Writer) *tablewriter.Table {
	t := tablewriter.NewWriter(w)
	t.SetBorder(false)
	t.SetColumnSeparator(" ")
	return t
}

// NewListTable builds a table formatted to list a set of results.
func NewListTable() *tablewriter.Table {
	return NewListTable2(os.Stdout)
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
