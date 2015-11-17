package foo

import (
	"fmt"
	. "github.com/manuviswam/gauge-go/gauge"
)

func foobar() {
	fmt.Println("Bar executed in a.go")
}

var _ = Describe("This is another awesome test", foobar)
