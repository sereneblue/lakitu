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
	// add function to check if registry value exists
	if strings.Contains(strings.ToLower(args[0]), "test-registryvalue") {
		args[0] = `
			function Test-RegistryValue {
			    param (
				    [parameter(Mandatory=$true)]
					[ValidateNotNullOrEmpty()] $Path,

			    	[parameter(Mandatory=$true)]
					[ValidateNotNullOrEmpty()] $Value
			    )

			    try {
			        Get-ItemProperty -Path $Path | Select-Object -ExpandProperty $Value -ErrorAction Stop | Out-Null
			        return $true
		        } catch {
			        return $false
		        }
			}
		` + args[0]
	}

	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.psPath, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String())
	return
}
