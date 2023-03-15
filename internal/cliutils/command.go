package cliutils

import "os/exec"

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
