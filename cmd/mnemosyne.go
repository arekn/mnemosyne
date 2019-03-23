// +build linux

package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

const memInfoPath = "/proc/meminfo"

func main() {

	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Printf("%v + %v \n", f.Name(), f.IsDir())
	}

	//memInfoFile, openFileError := os.Open(memInfoPath)
	//
	//defer memInfoFile.Close()
	//if openFileError != nil {
	//	log.Fatal(openFileError)
	//}
	//
	//procfs.ParseProcFile(memInfoFile)
}
