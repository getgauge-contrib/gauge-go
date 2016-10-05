package gauge

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/getgauge-contrib/gauge-go/constants"
	"github.com/getgauge-contrib/gauge-go/util"
	"github.com/getgauge/common"
)

// LoadGaugeImpls builds the go project and runs the generated go file,
// so that the gauge specific implementations get scanned
func LoadGaugeImpls() error {
	var b bytes.Buffer
	buff := bufio.NewWriter(&b)

	if err := util.RunCommand(os.Stdout, os.Stdout, constants.CommandGo, "build", "./..."); err != nil {
		buff.Flush()
		return fmt.Errorf("Build failed: %s\n", err.Error())
	}

	// get list of all packages in the projectRoot
	if err := util.RunCommand(buff, buff, constants.CommandGo, "list", "./..."); err != nil {
		buff.Flush()
		fmt.Printf("Failed to get the list of all packages: %s\n%s", err.Error(), b.String())
	}

	tempDir := common.GetTempDir()
	defer os.RemoveAll(tempDir)

	gaugeGoMainFile := filepath.Join(tempDir, constants.GaugeTestMainFileName)
	f, err := os.Create(gaugeGoMainFile)
	if err != nil {
		return fmt.Errorf("Failed to create main file in %s: %s", tempDir, err.Error())
	}

	genGaugeTestFileContents(f, b.String())
	f.Close()
	// Scan gauge methods
	if err := util.RunCommand(os.Stdout, os.Stdout, constants.CommandGo, "run", gaugeGoMainFile); err != nil {
		return fmt.Errorf("Failed to compile project: %s\nPlease ensure the project is in GOPATH.\n", err.Error())
	}
	return nil
}

func genGaugeTestFileContents(fileWriter io.Writer, importString string) {
	type info struct {
		Imports []string
	}
	var validImports []string
	for _, i := range strings.Fields(importString) {
		if strings.HasPrefix(i, "_") {
			validImports = append(validImports, strings.TrimPrefix(i, "_"))
		} else {
			validImports = append(validImports, i)
		}
	}
	tplMain.Execute(fileWriter, info{Imports: validImports})
}

var tplMain = template.Must(template.New("main").Parse(
	`package main
import (
	"github.com/getgauge-contrib/gauge-go/gauge"
{{range $n, $i := .Imports}}	_ "{{$i}}"
{{end}})
func main() {
	gauge.Run()
}
`))
