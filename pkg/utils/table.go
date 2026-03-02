// pkg/utils/table.go
package utils

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// TableFormatter provides consistent table formatting across the application
type TableFormatter struct {
	table *tablewriter.Table
}

// NewTableFormatter creates a new table formatter with consistent styling
func NewTableFormatter(headers []string) *TableFormatter {
	rendition := tw.Rendition{
		Borders: tw.Border{
			Left:      tw.Off,
			Right:     tw.Off,
			Top:       tw.Off,
			Bottom:    tw.Off,
			Overwrite: true,
		},
		Symbols: tw.NewSymbols(tw.StyleNone),
		Settings: tw.Settings{
			Separators: tw.Separators{
				ShowHeader:     tw.Off,
				ShowFooter:     tw.Off,
				BetweenRows:    tw.Off,
				BetweenColumns: tw.Off,
			},
			Lines: tw.Lines{
				ShowTop:        tw.Off,
				ShowBottom:     tw.Off,
				ShowHeaderLine: tw.Off,
				ShowFooterLine: tw.Off,
			},
		},
	}

	table := tablewriter.NewTable(
		os.Stdout,
		tablewriter.WithHeader(headers),
		tablewriter.WithHeaderAlignment(tw.AlignLeft),
		tablewriter.WithRowAlignment(tw.AlignLeft),
		tablewriter.WithPadding(tw.Padding{Left: "  ", Right: "  ", Overwrite: true}),
		tablewriter.WithTrimSpace(tw.On),
		tablewriter.WithTrimLine(tw.On),
		tablewriter.WithRendition(rendition),
	)

	return &TableFormatter{
		table: table,
	}
}

// AppendRow adds a row to the table
func (t *TableFormatter) AppendRow(row []string) {
	// tablewriter v1 returns an error for type conversion failures; we keep the
	// caller API simple and just surface failures on stderr.
	if err := t.table.Append(row); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to append table row: %v\n", err)
	}
}

// SetColumnMinWidth sets minimum width for a column
func (t *TableFormatter) SetColumnMinWidth(column int, width int) {
	// tablewriter v1 doesn't expose a direct "min width" API.
	// Best-effort: set a per-column width constraint.
	t.table.Configure(func(cfg *tablewriter.Config) {
		if cfg.Widths.PerColumn == nil {
			cfg.Widths.PerColumn = tw.NewMapper[int, int]()
		}
		cfg.Widths.PerColumn[column] = width
	})
}

// Render displays the table
func (t *TableFormatter) Render() {
	if err := t.table.Render(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to render table: %v\n", err)
	}
}
