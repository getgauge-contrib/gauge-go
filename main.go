package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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
	if requiresGoModuleFile(projectRoot) {
		fmt.Printf("Failed to start runner; go.mod is required when working outside the GOPATH.\nCreate it using `go mod init <module-name>`\n")
		os.Exit(1)
	}
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

	if requiresGoModuleFile(projectRoot) {
		showMessage("create", "go.mod")

		cmd := exec.Command("go", "mod", "init")
		cmd.Dir = projectRoot

		err := cmd.Run()
		if err != nil {
			// go mod init without parameters only works if it can find an import hint or GoDep config or verdoring config or the git origin is on github
			// in other situations, it will fail, in which case the user must run it themselves with the correct module name
			fmt.Printf(" could not create go.mod; create it yourself using `go mod init <module name>`\n")
		}
	}
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

	if !checkIfInSrcPath(projectRoot) && !checkGoModulesAvailable() {
		fmt.Printf("Project folder must be a subfolder in GOPATH/src folder, or Go Modules must be available (go1.11+)\n")
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

func getGoPaths() []string {
	var paths []string
	if runtime.GOOS == "windows" {
		paths = strings.Split(build.Default.GOPATH, ";")
	} else {
		paths = strings.Split(build.Default.GOPATH, ":")
	}
	return paths
}

func getGoSrcPaths() []string {
	var paths = getGoPaths()
	for i, p := range paths {
		paths[i] = filepath.Join(p, "src")
	}
	return paths
}

func checkIfInSrcPath(dirPath string) bool {
	for _, p := range getGoSrcPaths() {
		if filepath.HasPrefix(dirPath, p) {
			return true
		}
	}
	return false
}

func checkGoModulesAvailable() bool {
	minReleaseTag := "go1.11"
	for _, tag := range build.Default.ReleaseTags {
		if minReleaseTag == tag {
			return true
		}
	}
	return false
}

func requiresGoModuleFile(dirPath string) bool {
	// Module file is never required in GOPATH
	if checkIfInSrcPath(dirPath) {
		return false
	}

	// Look for a module file; if one exists somewhere up the tree everything is OK
	dir := dirPath
	// Traverse up to the root volume. On Windows VolumeName returns C:, on *nix it returns an empty string
	root := filepath.VolumeName(dirPath) + string(os.PathSeparator)
	for dir != root {
		modPath := filepath.Join(dir, `go.mod`)
		if common.FileExists(modPath) {
			return false
		}
		dir = filepath.Dir(dir)
	}

	// Outside of GOPATH and no module file could be found, so one is needed
	return true
}