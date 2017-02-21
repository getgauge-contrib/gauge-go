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

	expected := `package gauge_test_bootstrap
import (
	"os"
	"github.com/getgauge-contrib/gauge-go/gauge"
	_ "github.com/apoorvam/foo/stepImpl"
	_ "github.com/apoorvam/foo/stepImpl/impl"
)
func init() {
	gauge.Run()
	_, w, _ := os.Pipe()
	os.Stderr = w
	os.Stdout = w
}
`
	assert.Equal(tst, expected, b.String())
}
