package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/getgauge/common"
	"github.com/manuviswam/gauge-go/constants"
	"io/ioutil"
)

var pluginDir = ""
var projectRoot = ""
var start = flag.Bool("start", false, "Start go runner")
var initialize = flag.Bool("init", false, "Initialize go specs")

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
	fmt.Println("Start called")
}

func initGo() {
	os.Chdir(projectRoot)
	createDirectory(constants.DefaultSpecImplDir)
	createFile(path.Join(constants.DefaultSpecImplDir, constants.DefaultStepImplFileName), constants.HelloWorldImplTemplate)
	createFile(path.Join(constants.DefaultSpecImplDir, constants.DefaultInitTestFileName), constants.InitTestTemplate)
}

func printUsage() {
	flag.PrintDefaults()
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
}

func createDirectory(dirPath string) {
	fmt.Println("create ", dirPath)
	if !common.DirExists(dirPath) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			fmt.Printf("Failed to make directory. %s\n", err.Error())
		}
	} else {
		fmt.Println("skip ", dirPath)
	}
}

func createFile(filePath, content string) {
	fmt.Println("create ", filePath)
	if !common.FileExists(filePath) {
		err := ioutil.WriteFile(filePath, []byte(content), common.NewFilePermissions)
		if err != nil {
			fmt.Println("Error creating file : ", err)
		}
	} else {
		fmt.Println("skip ", filePath)
	}

}
