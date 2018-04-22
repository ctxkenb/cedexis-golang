package main

import (
	"fmt"
	"io"
	"strings"
)

// Table formats for output data structured in rows and columns
type Table struct {
	Columns []string
	Rows    [][]string
}

// Print outputs the table to a writer with a given width (to avoid wrapping fields may be truncated)
func (t *Table) Print(w io.Writer, width int) {
	colWidths := make([]int, len(t.Columns))
	updateColWidths(t.Columns, &colWidths)
	for _, r := range t.Rows {
		updateColWidths(r, &colWidths)
	}

	totalWidth := 0
	largestCol := 0
	for i, c := range colWidths {
		if c > colWidths[largestCol] {
			largestCol = i
		}
		if i > 0 {
			totalWidth++
		}
		totalWidth += c
	}

	toCut := totalWidth - width
	if toCut > 0 {
		colWidths[largestCol] -= toCut
	}

	line := ""
	for i, c := range t.Columns {
		if i > 0 {
			line += " "
		}
		line += fmt.Sprintf("%*.*s", -colWidths[i], colWidths[i], c)
	}
	fmt.Fprintln(w, line)

	line = ""
	for i := range t.Columns {
		if i > 0 {
			line += " "
		}
		line += strings.Repeat("=", colWidths[i])
	}
	fmt.Fprintln(w, line)

	for _, r := range t.Rows {
		line = ""
		for i, c := range r {
			if i > 0 {
				line += " "
			}
			line += fmt.Sprintf("%*.*s", -colWidths[i], colWidths[i], c)
		}
		fmt.Fprintln(w, line)
	}
}

func updateColWidths(a []string, colWidths *[]int) {
	for i, v := range a {
		if len(v) > (*colWidths)[i] {
			(*colWidths)[i] = len(v)
		}
	}
}
