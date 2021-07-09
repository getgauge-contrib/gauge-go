package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/getgauge/common"
)

const (
	CGO_ENABLED = "CGO_ENABLED"
	dotGauge    = ".gauge"
	plugins     = "plugins"
	GOARCH      = "GOARCH"
	GOOS        = "GOOS"
	X86         = "386"
	X86_64      = "amd64"
	ARM64       = "arm64"
	DARWIN      = "darwin"
	LINUX       = "linux"
	WINDOWS     = "windows"
	bin         = "bin"
	gauge       = "gauge"
	gaugeGo     = "gauge-go"
	deploy      = "deploy"
	targetDir   = "target"
)

var deployDir = filepath.Join(deploy, gaugeGo)

var install = flag.Bool("install", false, "Install to the specified prefix")
var pluginInstallPrefix = flag.String("plugin-prefix", "", "Specifies the prefix where gauge plugins will be installed")
var distro = flag.Bool("distro", false, "Creates distributables for gauge go")
var test = flag.Bool("test", false, "Runs tests")
var allPlatforms = flag.Bool("all-platforms", false, "Compiles or creates distributables for all platforms windows, linux, darwin both x86 and x86_64")
var binDir = flag.String("bin-dir", "", "Specifies OS_PLATFORM specific binaries to install when cross compiling")

var (
	platformEnvs = []map[string]string{
		{GOARCH: ARM64, GOOS: DARWIN, CGO_ENABLED: "0"},
		{GOARCH: X86_64, GOOS: DARWIN, CGO_ENABLED: "0"},
		{GOARCH: X86, GOOS: LINUX, CGO_ENABLED: "0"},
		{GOARCH: X86_64, GOOS: LINUX, CGO_ENABLED: "0"},
		{GOARCH: X86, GOOS: WINDOWS, CGO_ENABLED: "0"},
		{GOARCH: X86_64, GOOS: WINDOWS, CGO_ENABLED: "0"},
	}
)

func main() {
	flag.Parse()

	if *install {
		updatePluginInstallPrefix()
		installGaugeGo(*pluginInstallPrefix)
	} else if *distro {
		createGaugeDistro(*allPlatforms)
	} else if *test {
		compileGoPackage(gaugeGo)
		runGoUnitTests()
	} else {
		compileGaugeGo()
	}
}

func runGoUnitTests() {
	runCommand("go", "test", "./...")
}

func compileGaugeGo() {
	if *allPlatforms {
		compileGaugeGoAcrossPlatforms()
	} else {
		compileGoPackage(gaugeGo)
	}
}

func compileGaugeGoAcrossPlatforms() {
	for _, platformEnv := range platformEnvs {
		setEnv(platformEnv)
		fmt.Printf("Compiling for platform => OS:%s ARCH:%s \n", platformEnv[GOOS], platformEnv[GOARCH])
		compileGoPackage(gaugeGo)
	}
}

func compileGoPackage(packageName string) {
	runCommand("go", "build")
	moveBinaryToBinDir()
}

func createGaugeDistro(forAllPlatforms bool) {
	if forAllPlatforms {
		for _, platformEnv := range platformEnvs {
			setEnv(platformEnv)
			fmt.Printf("Creating distro for platform => OS:%s ARCH:%s \n", platformEnv[GOOS], platformEnv[GOARCH])
			createDistro()
		}
	} else {
		createDistro()
	}
}

func createDistro() {
	packageName := fmt.Sprintf("%s-%s-%s.%s", gaugeGo, getGaugeGoVersion(), getGOOS(), getArch())
	distroDir := filepath.Join(deploy, packageName)
	copyGaugeGoFiles(distroDir)
	createZipFromUtil(deploy, packageName)
	os.RemoveAll(distroDir)
}

func copyGaugeGoFiles(destDir string) {
	files := make(map[string]string)
	if getGOOS() == WINDOWS {
		files[filepath.Join(getBinDir(), "gauge-go.exe")] = bin
	} else {
		files[filepath.Join(getBinDir(), gaugeGo)] = bin
	}

	files[filepath.Join("go.json")] = ""
	files[filepath.Join("stepImpl", "stepImplementation.go")] = filepath.Join("skel")
	copyFiles(files, destDir)
}

func moveBinaryToBinDir() {
	srcFile := gaugeGo
	destDir := getBinDir()
	if getGOOS() == WINDOWS {
		srcFile = srcFile + ".exe"
	}
	destFile := filepath.Join(destDir, srcFile)
	if err := os.MkdirAll(destDir, common.NewDirectoryPermissions); err != nil {
		fmt.Printf("Failed to create directory %s. %s\n", destDir, err.Error())
	}
	if err := common.MirrorFile(srcFile, destFile); err != nil {
		fmt.Printf("Failed to copy file %s to %s. %s\n", srcFile, destFile, err.Error())
	}
	os.Remove(srcFile)
}

func installGaugeGo(installPrefix string) {
	defer os.RemoveAll(deploy)
	copyGaugeGoFiles(deployDir)
	goRunnerInstallPath := filepath.Join(installPrefix, "go", getGaugeGoVersion())
	log.Printf("Copying %s -> %s\n", deployDir, goRunnerInstallPath)
	common.MirrorDir(deployDir, goRunnerInstallPath)
}

func updatePluginInstallPrefix() {
	if *pluginInstallPrefix == "" {
		if runtime.GOOS == WINDOWS {
			*pluginInstallPrefix = os.Getenv("APPDATA")
			if *pluginInstallPrefix == "" {
				panic(fmt.Errorf("Failed to find AppData directory"))
			}
			*pluginInstallPrefix = filepath.Join(*pluginInstallPrefix, gauge, plugins)
		} else {
			userHome := getUserHome()
			if userHome == "" {
				panic(fmt.Errorf("Failed to find User Home directory"))
			}
			*pluginInstallPrefix = filepath.Join(userHome, dotGauge, plugins)
		}
	}
}

func getPluginProperties(jsonPropertiesFile string) (map[string]interface{}, error) {
	pluginPropertiesJSON, err := ioutil.ReadFile(jsonPropertiesFile)
	if err != nil {
		fmt.Printf("Could not read %s: %s\n", filepath.Base(jsonPropertiesFile), err)
		return nil, err
	}
	var pluginJSON interface{}
	if err = json.Unmarshal([]byte(pluginPropertiesJSON), &pluginJSON); err != nil {
		fmt.Printf("Could not read %s: %s\n", filepath.Base(jsonPropertiesFile), err)
		return nil, err
	}
	return pluginJSON.(map[string]interface{}), nil
}

func getGaugeGoVersion() string {
	goRunnerProperties, err := getPluginProperties("go.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to get gauge go properties file. %s", err))
	}
	return goRunnerProperties["version"].(string)
}

func getBinDir() string {
	if *binDir == "" {
		return filepath.Join(bin, fmt.Sprintf("%s_%s", getGOOS(), getArch()))
	}
	return filepath.Join(bin, *binDir)
}

func runCommand(command string, arg ...string) {
	cmd := exec.Command(command, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("Execute %v\n", cmd.Args)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func setEnv(envVariables map[string]string) {
	for k, v := range envVariables {
		os.Setenv(k, v)
	}
}

// key will be the source file and value will be the target
func copyFiles(files map[string]string, installDir string) {
	for src, dst := range files {
		base := filepath.Base(src)
		installDst := filepath.Join(installDir, dst)
		log.Printf("Copying %s -> %s\n", src, installDst)
		stat, err := os.Stat(src)
		if err != nil {
			panic(err)
		}
		if stat.IsDir() {
			_, err = common.MirrorDir(src, installDst)
		} else {
			err = common.MirrorFile(src, filepath.Join(installDst, base))
		}
		if err != nil {
			panic(err)
		}
	}
}

func getUserHome() string {
	return os.Getenv("HOME")
}

func getArch() string {
	arch := getGOARCH()
	if arch == ARM64 {
		return ARM64
	}
	if arch == X86 {
		return "x86"
	}
	return "x86_64"
}

func getGOARCH() string {
	goArch := os.Getenv(GOARCH)
	if goArch == "" {
		return runtime.GOARCH

	}
	return goArch
}

func getGOOS() string {
	os := os.Getenv(GOOS)
	if os == "" {
		return runtime.GOOS

	}
	return os
}

func isExecMode(mode os.FileMode) bool {
	return (mode & 0111) != 0
}

func createZipFromUtil(dir, name string) {
	wd, _ := os.Getwd()
	os.Chdir(filepath.Join(dir, name))
	runCommand("zip", "-r", filepath.Join("..", name+".zip"), ".")
	os.Chdir(wd)
}
