// pkg/utils/table.go
package utils

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

// TableFormatter provides consistent table formatting across the application
type TableFormatter struct {
	table *tablewriter.Table
}

// NewTableFormatter creates a new table formatter with consistent styling
func NewTableFormatter(headers []string) *TableFormatter {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	// Set consistent styling
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetTablePadding("  ")
	table.SetNoWhiteSpace(true)

	return &TableFormatter{
		table: table,
	}
}

// AppendRow adds a row to the table
func (t *TableFormatter) AppendRow(row []string) {
	t.table.Append(row)
}

// SetColumnMinWidth sets minimum width for a column
func (t *TableFormatter) SetColumnMinWidth(column int, width int) {
	t.table.SetColMinWidth(column, width)
}

// Render displays the table
func (t *TableFormatter) Render() {
	t.table.Render()
}
