package cli

import (
	"bytes"
	"os/exec"
	"strings"
)

type Shell struct {
	psPath string
}

func NewShell() *Shell {
	ps, _ := exec.LookPath("powershell.exe")
	return &Shell{
		psPath: ps,
	}
}

func (p *Shell) execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.psPath, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String())

	return stdOut, stdErr, err
}
