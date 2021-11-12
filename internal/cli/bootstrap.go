package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	WS_2016_RELEASE_ID = "1607"
	WS_2019_RELEASE_ID = "1809"
)

func BootstrapParsec() error {
	if !isAWS() {
		return errors.New("Machine is not supported AWS instance.")
	}

	ps := NewShell()
	if !isSupportedWindowsServer(ps) {
		return errors.New("Machine is not Windows Server 2016 or Windows Server 2019.")
	}

	ps.execute(`
		[Net.ServicePointManager]::SecurityProtocol = "tls12, tls11, tls" 
		$ScriptWebArchive = "https://github.com/parsec-cloud/Parsec-Cloud-Preparation-Tool/archive/master.zip"  
		$LocalArchivePath = "$ENV:UserProfile\Downloads\Parsec-Cloud-Preparation-Tool"  
		(New-Object System.Net.WebClient).DownloadFile($ScriptWebArchive, "$LocalArchivePath.zip")  
		Expand-Archive "$LocalArchivePath.zip" -DestinationPath $LocalArchivePath -Force  
		CD $LocalArchivePath\Parsec-Cloud-Preparation-Tool-master\ | powershell.exe .\Loader.ps1 
	`)

	return nil
}

func isAWS() bool {
	fmt.Println("Checking if machine is supported AWS instance...")
	client := http.Client{Timeout: 3 * time.Second}
	res, err := client.Get("http://169.254.169.254/latest/meta-data/instance-type")

	if err != nil {
		return false
	}

	defer res.Body.Close()

	// check if supported instance
	body, err := ioutil.ReadAll(res.Body)
	instanceType := strings.Split(string(body), ".")

	if instanceType[0] == "g3s" || instanceType[0] == "g4dn" {
		return true
	}

	return false
}

func isSupportedWindowsServer(ps *Shell) bool {
	fmt.Println("Checking if machine is supported Windows Server...")
	stdOut, _, err := ps.execute(`
		(Get-ItemProperty "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion" -Name ReleaseID).ReleaseID
	`)

	if err != nil || (stdOut != WS_2016_RELEASE_ID && stdOut != WS_2019_RELEASE_ID) {
		return false
	}

	return true
}
