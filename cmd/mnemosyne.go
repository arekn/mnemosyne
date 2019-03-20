// +build linux

package main

import (
	"github.com/arekn/mnemosyne/internal/pkg"
	"log"
	"os"
)

const memInfoPath = "/proc/meminfo"

func main() {
	memInfoFile, openFileError := os.Open(memInfoPath)

	defer memInfoFile.Close()
	if openFileError != nil {
		log.Fatal(openFileError)
	}

	memory.ParseMemInfo(memInfoFile)
}
