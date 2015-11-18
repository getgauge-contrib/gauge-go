package constants

const HelloWorldImplTemplate = `
package foo

import (
	"fmt"

	. "github.com/manuviswam/gauge-go/gauge"
	"github.com/manuviswam/gauge-go/models"
)

var _ = Describe("Say <greeting> to <product name>", func(greeting string, name string) {
	fmt.Println(greeting + ", " + name)
})

var _ = Describe("Step that takes a table <table>", func(table models.Table) {
	for _, column := range table.Columns() {
		fmt.Println(column)
	}
	for _, row := range table.Rows() {
		fmt.Println(row)
	}
})
`

const InitTestTemplate = "package gauge\n/* Do not delete or modify this file */"
