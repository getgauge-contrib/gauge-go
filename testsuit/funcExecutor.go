package testsuit

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/constants"
	"github.com/getgauge-contrib/gauge-go/gauge_messages"
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
	res = &m.ProtoExecutionResult{}
	T = &testingT{}

	rargs := make([]reflect.Value, len(args))
	for i, a := range args {
		if s, ok := a.(string); ok {
			castedVal, err := convertType(fn.Type().In(i), s)
			if err != nil {
				res.ScreenShot = getScreenshot()
				res.Failed = true
				res.ExecutionTime = 0
				res.ErrorType = gauge_messages.ProtoExecutionResult_VERIFICATION
				res.StackTrace = " " // make sure that the error message is displayed on the report html
				res.ErrorMessage = err.Error()
				return res
			}
			rargs[i] = castedVal
			continue
		}
		rargs[i] = reflect.ValueOf(a)
	}
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			res.ScreenShot = getScreenshot()
			res.Failed = true
			res.ExecutionTime = time.Since(start).Nanoseconds()
			res.StackTrace = strings.SplitN(string(debug.Stack()), "\n", 9)[8]
			res.ErrorMessage = fmt.Sprintf("%s", r)
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
	}
	res.ExecutionTime = time.Since(start).Nanoseconds()
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

func convertType(t reflect.Type, strVal string) (reflect.Value, error) {
	if t.Kind() == reflect.String {
		return reflect.ValueOf(strVal), nil
	}
	tBitSize := int(t.Size()) * 8

	val := reflect.New(t)
	var err error

	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, e := strconv.ParseInt(strVal, 10, tBitSize)
		if e == nil {
			val.Elem().SetInt(n)
		}
		err = e
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, e := strconv.ParseUint(strVal, 10, tBitSize)
		if e == nil {
			val.Elem().SetUint(n)
		}
		err = e
	case reflect.Float32, reflect.Float64:
		n, e := strconv.ParseFloat(strVal, tBitSize)
		if e == nil {
			val.Elem().SetFloat(n)
		}
		err = e
	default:
		return val, fmt.Errorf("cannot convert a string to a %s value", t.String())
	}
	if err != nil {
		err = fmt.Errorf("cannot convert %#v to a %s value: %s", strVal, t.String(), err.Error())
	}
	return val.Elem(), err
}
