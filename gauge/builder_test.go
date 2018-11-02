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
	importString := `foo/stepImpl
foo/stepImpl/impl`

	genGaugeTestFileContents(buffWriter, importString, "")
	buffWriter.Flush()

	expected := `package gauge_test_bootstrap
import (
	"os"
	"github.com/getgauge-contrib/gauge-go/gauge"
	_ "foo/stepImpl"
	_ "foo/stepImpl/impl"
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

func TestGeneratedMainFileContentsShouldFilterOutVendorPackages(tst *testing.T) {
	var b bytes.Buffer
	buffWriter := bufio.NewWriter(&b)
	importString := `github.com/foo/bar/stepImpl
github.com/foo/bar/stepImpl/impl
github.com/foo/bar/vendor/vendorStepImpl/impl`
	vendorString := `github.com/foo/bar/vendor/vendorStepImpl/impl`

	genGaugeTestFileContents(buffWriter, importString, vendorString)
	buffWriter.Flush()

	expected := `package gauge_test_bootstrap
import (
	"os"
	"github.com/getgauge-contrib/gauge-go/gauge"
	_ "github.com/foo/bar/stepImpl"
	_ "github.com/foo/bar/stepImpl/impl"
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

