package testsuit

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/getgauge-contrib/gauge-go/constants"
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	"github.com/getgauge-contrib/gauge-go/util"
	"github.com/getgauge/common"
)

const (
	screenshotFileName = "screenshot.png"
)

var CustomScreenShot *func() []byte

// TODO: Use gauge-go result object rather than ProtoExecutionResult
func executeFunc(fn reflect.Value, args ...interface{}) (res *m.ProtoExecutionResult) {
	rargs := make([]reflect.Value, len(args))
	for i, a := range args {
		rargs[i] = reflect.ValueOf(a)
	}
	res = &m.ProtoExecutionResult{}
	T = &testingT{}
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			res.ScreenShot = getScreenshot()
			res.Failed = true
			res.ExecutionTime = time.Since(start).Nanoseconds() / int64(time.Millisecond)
			res.StackTrace = caller(4)
			res.ErrorMessage = fmt.Sprintf("%s", r)
			res.RecoverableError = T.getContinueOnFailure()
		}
		T = &testingT{}
	}()
	fn.Call(rargs)
	res.Failed = false
	if len(T.errors) != 0 {
		res.ScreenShot = getScreenshot()
		res.Failed = true
		res.StackTrace = T.getStacktraces()
		res.ErrorMessage = T.getErrors()
		res.RecoverableError = T.getContinueOnFailure()
	}
	res.ExecutionTime = time.Since(start).Nanoseconds() / int64(time.Millisecond)
	return res
}

func getScreenshot() []byte {
	if os.Getenv(constants.ScreenshotOnFailure) == "true" {
		if *CustomScreenShot != nil {
			fn := reflect.ValueOf(*CustomScreenShot)
			screenShotBytes := fn.Call(make([]reflect.Value, 0))
			return screenShotBytes[0].Interface().([]byte)
		}
		tmpDir := common.GetTempDir()
		defer os.RemoveAll(tmpDir)
		var b bytes.Buffer
		buff := bufio.NewWriter(&b)
		screenshotFile := filepath.Join(tmpDir, screenshotFileName)
		util.RunCommand(buff, buff, constants.GaugeScreenshot, screenshotFile)
		bytes, err := ioutil.ReadFile(screenshotFile)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		return bytes
	}
	return nil
}

func caller(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("\t%s:%d", file, line)
}
