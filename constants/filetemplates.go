package constants

const HelloWorldImplTemplate = `
package specimpl

import (
	"fmt"

	m "github.com/manuviswam/gauge-go/models"
	. "github.com/manuviswam/gauge-go/runner"
)

var _ = Describe("Say <greeting> to <product name>", func(greeting, name string) {
	fmt.Println(greeting + ", " + name)
})

var _ = Describe("Step that takes a table <table>", func(tbl *m.Table) {
	for _, columnName := range tbl.Headers.Cells {
		fmt.Printf("%s,\t", columnName)
	}

	fmt.Printf("\n")
	for _, row := range tbl.Rows {
		for _, cell := range row.Cells {
			fmt.Printf("%s\t", cell)
		}
		fmt.Printf("\n")
	}
})

var _ = Describe("A context step which gets executed before every scenario", func() {
	fmt.Println("Context step executed")
})
`

const InitTestTemplate = `
/* Do not delete or modify this file */
package specimpl

import (
	. "github.com/manuviswam/gauge-go/runner"
	"testing"
)

func TestInitialize(t *testing.T) {
	Run()
}
`
