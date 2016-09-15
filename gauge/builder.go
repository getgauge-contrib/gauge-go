package gauge

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/getgauge/common"
	"github.com/manuviswam/gauge-go/constants"
)

// LoadGaugeImpls builds the go project and runs the generated go file,
// so that the gauge specific implementations get scanned
func LoadGaugeImpls() error {
	var b bytes.Buffer
	buff := bufio.NewWriter(&b)

	if err := runCommand(os.Stdout, os.Stdout, constants.CommandGo, "build", "./..."); err != nil {
		buff.Flush()
		return fmt.Errorf("Build failed: %s\n", err.Error())
	}

	// get list of all packages in the projectRoot
	if err := runCommand(buff, buff, constants.CommandGo, "list", "./..."); err != nil {
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
	if err := runCommand(os.Stdout, os.Stdout, constants.CommandGo, "run", gaugeGoMainFile); err != nil {
		return fmt.Errorf("Failed to run file %s: %s\n", gaugeGoMainFile, err.Error())
	}
	return nil
}

func genGaugeTestFileContents(fileWriter io.Writer, importString string) {
	imports := strings.Fields(importString)
	type info struct {
		Imports []string
	}
	tplMain.Execute(fileWriter, info{Imports: imports})
}

var tplMain = template.Must(template.New("main").Parse(
	`package main
import (
	"github.com/manuviswam/gauge-go/gauge"
{{range $n, $i := .Imports}}	_ "{{$i}}"
{{end}})
func main() {
	gauge.Run()
}
`))

func runCommand(stdOut, stdErr io.Writer, command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	return cmd.Run()
}
