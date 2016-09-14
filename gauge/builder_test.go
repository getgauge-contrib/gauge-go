package gauge

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMainFileContents(tst *testing.T) {
	var b bytes.Buffer
	buffWriter := bufio.NewWriter(&b)
	importString := `github.com/apoorvam/foo/stepImpl
github.com/apoorvam/foo/stepImpl/impl`

	genGaugeTestFileContents(buffWriter, importString)
	buffWriter.Flush()

	expected := `package main
import (
	"github.com/manuviswam/gauge-go/gauge"
	_ "github.com/apoorvam/foo/stepImpl"
	_ "github.com/apoorvam/foo/stepImpl/impl"
)
func main() {
	gauge.Run()
}
`
	assert.Equal(tst, expected, b.String())
}
