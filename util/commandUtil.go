package util

import (
	"io"
	"os/exec"
)

func RunCommand(stdOut, stdErr io.Writer, command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	return cmd.Run()
}
