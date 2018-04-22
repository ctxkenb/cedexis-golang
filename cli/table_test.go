package main

import (
	"bytes"
	"testing"
)

func TestPrint(t *testing.T) {
	tbl := Table{
		Columns: []string{"Name", "Cost"},
		Rows: [][]string{
			{"This is a big name", "$1000"},
			{"name", "$$$ very large value"},
		},
	}

	expected := "Name               Cost       \n" +
		"================== ===========\n" +
		"This is a big name $1000      \n" +
		"name               $$$ very la\n"

	buff := bytes.Buffer{}
	tbl.Print(&buff, 30)
	s := buff.String()

	if s != expected {
		t.Errorf("Incorrect table formatting, got %s, want %s.", s, expected)
	}
}
