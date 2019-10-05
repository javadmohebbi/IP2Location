package main

import (
	"os"
	"runtime"
)

func main() {

	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// filename := os.Args[1] // get command line first parameter

	// filedirectory := filepath.Dir(filename)

	// thepath, err := filepath.Abs(filedirectory)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(thepath)
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
