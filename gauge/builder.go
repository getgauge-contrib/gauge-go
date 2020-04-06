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
func LoadGaugeImpls(projectRoot string) error {
	var b bytes.Buffer
	buff := bufio.NewWriter(&b)

	os.Chdir(projectRoot)

	if err := util.RunCommand(os.Stdout, os.Stdout, constants.CommandGo, "build", "./..."); err != nil {
		buff.Flush()
		return fmt.Errorf("Build failed: %s\n", err.Error())
	}

	// get list of all packages in the projectRoot
	if err := util.RunCommand(buff, buff, constants.CommandGo, "list", "./..."); err != nil {
		buff.Flush()
		fmt.Printf("Failed to get the list of all packages: %s\n%s", err.Error(), b.String())
	}

	// if using any vendor tool, remove the vendor packages from list of packages
	vendorPackages, _ := getVendorPackageList(projectRoot)

	tempDir := common.GetTempDir()
	defer os.RemoveAll(tempDir)

	gaugeGoMainFile := filepath.Join(tempDir, constants.GaugeTestFileName)
	f, err := os.Create(gaugeGoMainFile)
	if err != nil {
		return fmt.Errorf("Failed to create main file in %s: %s", tempDir, err.Error())
	}

	genGaugeTestFileContents(f, b.String(), vendorPackages)
	f.Close()
	// Scan gauge methods
	if err := util.RunCommand(os.Stdout, os.Stdout, constants.CommandGo, "test", "-timeout", "0", "-v", gaugeGoMainFile); err != nil {
		return fmt.Errorf("Failed to compile project: %s\nPlease ensure the project is in GOPATH.\n", err.Error())
	}
	return nil
}

func getVendorPackageList(projectRoot string) (string, error) {
	vendorDir := filepath.Join(projectRoot, "vendor")
	if _, err := os.Stat(vendorDir) ; err != nil {
		return "", nil
	}
	err := os.Chdir(vendorDir)
	defer os.Chdir(projectRoot)

	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	buff := bufio.NewWriter(&b)

	if err := util.RunCommand(buff, buff, constants.CommandGo, "list", "./..."); err != nil {
		buff.Flush()
		fmt.Printf("Failed to get the list of vendor packages: %s\n%s", err.Error(), b.String())
		return "", err
	}

	return b.String(), nil
}

func genGaugeTestFileContents(fileWriter io.Writer, importString, vendorImports string) {
	type info struct {
		Imports []string
	}

	vendorImportMap := make(map[string]bool)
	for _, vi := range strings.Fields(vendorImports) {
		vendorImportMap[vi] = true
	}

	var validImports []string
	for _, i := range strings.Fields(importString) {
		if vendorImportMap[i] {
			continue
		}

		if strings.HasPrefix(i, "_") {
			validImports = append(validImports, strings.TrimPrefix(i, "_"))
		} else {
			validImports = append(validImports, i)
		}
	}
	gaugeTestRunnerTpl.Execute(fileWriter, info{Imports: validImports})
}

var gaugeTestRunnerTpl = template.Must(template.New("main").Parse(
	`package gauge_test_bootstrap
import (
	"os"
	"github.com/getgauge-contrib/gauge-go/gauge"
{{range $n, $i := .Imports}}	_ "{{$i}}"
{{end}})
func init() {
	gauge.Run()
	_, w, _ := os.Pipe()
	os.Stderr = w
	os.Stdout = w
}
`))
