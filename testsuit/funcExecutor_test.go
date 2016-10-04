package testsuit

import (
	"os"
	"testing"

	"github.com/getgauge-contrib/gauge-go/constants"
	"github.com/stretchr/testify/assert"
)

func TestShouldRunCustomScreenshotFn(t *testing.T) {
	called := false
	screenshotStr := "hello"
	os.Setenv(constants.ScreenshotOnFailure, "true")
	fn := func() []byte {
		called = true
		return []byte(screenshotStr)
	}
	CustomScreenShot = &fn

	screenshot := getScreenshot()

	assert.True(t, called)
	assert.Equal(t, screenshotStr, string(screenshot))
}
