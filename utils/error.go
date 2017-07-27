package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// ExitWithError prints out error and exits program with status code 1
func ExitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

// PrintJSON prints to stdout marshaled JSON
//  Given argument should be any marshal - compatibile data structure
func PrintJSON(data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		ExitWithError(err)
	}

	fmt.Println(string(json))
}

// PrintTable prints to stdout an ascii table using github.com/olekukonko/tablewriter library
// Arguments are:
//   data: [][]string, data that should be printed as table
//   header: []string, headers of the table
func PrintTable(data [][]string, header []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(false)
	table.AppendBulk(data)

	table.Render()
}
