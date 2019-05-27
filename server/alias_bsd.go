// +build darwin freebsd

package server

import (
	"bytes"
	"io"
	"log"
	"os/exec"
)

func installNetworkAlias() ([]byte, error) {
	command := exec.Command("sysctl", "-w", "net.inet.ip.forwarding=1")
	sysctlOutput, err := command.CombinedOutput()
	log.Printf("sysctl: %s", sysctlOutput)
	if err != nil {
		return nil, err
	}

	command = exec.Command("pfctl", "-a", "com.apple/250.AWSVault", "-Ef", "-")

	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "rdr pass proto tcp from any to 169.254.169.254 port 80 -> 127.0.0.1 port 9099\n")
	}()

	pfctlOutput, err := command.CombinedOutput()
	log.Printf("pfctl: %s", pfctlOutput)
	if err != nil {
		return nil, err
	}

	output := bytes.Join([][]byte{sysctlOutput, pfctlOutput}, []byte("\n"))

	return output, nil
}
