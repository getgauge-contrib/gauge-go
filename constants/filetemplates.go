package constants

const HelloWorldImplTemplate = `
package specimpl

import (
	"fmt"

	. "github.com/manuviswam/gauge-go/runner"
)

var _ = Describe("Say <greeting> to <product name>", func() {
	fmt.Println("Greetings")
})

var _ = Describe("Step that takes a table <table>", func() {
	fmt.Println("More greetings")
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
