package main

import (
	"os"
	"runtime"
)

func main() {

	dir := ""
	tmpDir := ""
	if runtime.GOOS == "windows" {
		dir = os.Getenv("ProgramFiles") + "\\IP2Location\\"
		tmpDir = os.Getenv("%temp%") + "\\"
	} else if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		dir = "/etc/ip2location/"
		tmpDir = "/tmp/"
	}

	validateArgsAndCallFuncs(dir, tmpDir)

	os.Exit(0)
}
