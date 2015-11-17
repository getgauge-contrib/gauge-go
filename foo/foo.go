package foo

import (
	"fmt"
	. "github.com/manuviswam/gauge-go/gauge"
)

func bar() {
	fmt.Println("Bar executed")
}

var _ = Describe("This is an awesome test", bar)
