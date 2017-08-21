package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/getgauge-contrib/gauge-go/constants"
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge/common"
)

var pluginDir = ""
var projectRoot = ""
var start = flag.Bool("start", false, "Start go runner")
var initialize = flag.Bool("init", false, "Initialize Go project structure")

func main() {
	flag.Parse()

	setPluginAndProjectRoots()
	if *start {
		startGo()
	} else if *initialize {
		initGo()
	} else {
		printUsage()
	}
}

func startGo() {
	err := gauge.LoadGaugeImpls(projectRoot)
	if err != nil {
		fmt.Printf("Failed to build project: %s\nKilling go runner. \n", err.Error())
		os.Exit(1)
	}
}

func initGo() {
	stepImplDir := filepath.Join(projectRoot, constants.DefaultStepImplDir)
	createDirectory(stepImplDir)
	stepImplFile := filepath.Join(stepImplDir, constants.DefaultStepImplFileName)
	showMessage("create", stepImplFile)
	common.CopyFile(filepath.Join(constants.SkelDir, constants.DefaultStepImplFileName), stepImplFile)
}

func printUsage() {
	flag.PrintDefaults()
}

func showMessage(action, filename string) {
	fmt.Printf(" %s  %s\n", action, filename)
}

func setPluginAndProjectRoots() {
	var err error
	pluginDir, err = os.Getwd()
	if err != nil {
		fmt.Printf("Failed to find current working directory: %s \n", err)
		os.Exit(1)
	}
	projectRoot = os.Getenv(common.GaugeProjectRootEnv)
	if projectRoot == "" {
		fmt.Printf("Could not find %s env. Go Runner exiting...", common.GaugeProjectRootEnv)
		os.Exit(1)
	}

	goSrcPath := filepath.Join(os.Getenv("GOPATH"), "src")
	if !filepath.HasPrefix(projectRoot, goSrcPath) {
		fmt.Printf("Project folder must be a subfolder in GOPATH/src folder\n")
		os.Exit(1)
	}
}

func createDirectory(dirPath string) {
	showMessage("create", dirPath)
	if !common.DirExists(dirPath) {
		err := os.MkdirAll(dirPath, common.NewDirectoryPermissions)
		if err != nil {
			fmt.Printf("Failed to make directory. %s\n", err.Error())
		}
	} else {
		fmt.Println("skip ", dirPath)
	}
}
