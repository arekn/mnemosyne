// +build linux

package main

import (
	"encoding/csv"
	"fmt"
	"github.com/arekn/mnemosyne/pkg/procfs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"
)

const procPath = "/proc"
const memInfoFile = "meminfo"
const statusFile = "status"

func main() {
	now := time.Now()
	var memInfo procfs.MemInfoFile
	memInfo = parseFile(path.Join(procPath, memInfoFile))

	memFree, memFreeError := memInfo.MemFree()
	if memFreeError != nil {
		log.Fatal(memFreeError)
	}
	memTotal, memTotalError := memInfo.MemTotal()
	if memTotalError != nil {
		log.Fatal(memTotalError)
	}

	usedMemory := memTotal - memFree

	filename := fmt.Sprintf("meminfo-%v.csv", time.Now().Format("2006.01.02"))
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(f)
	w.Write([]string{now.Format("15:04:05"), strconv.Itoa(usedMemory), strconv.Itoa(memTotal)})
	w.Flush()

	//pidList := getProcessList()
	//var processFiles []procfs.StatusFile
	//for _, pid := range pidList {
	//	pathWithPid := path.Join(procPath, pid, statusFile)
	//	processFiles = append(processFiles, parseFile(pathWithPid))
	//}
	//
	//fmt.Println(len(pidList))
	//fmt.Println(len(processFiles))
	fmt.Println(time.Since(now))
}

func getProcessList() []string {
	files, err := ioutil.ReadDir(procPath)
	if err != nil {
		log.Fatal(err)
	}
	pidMatcher, compileError := regexp.Compile(`[0-9]`)
	if compileError != nil {
		log.Fatal(compileError)
	}
	var pidList []string
	for _, f := range files {

		if f.IsDir() && pidMatcher.MatchString(f.Name()) {
			pidList = append(pidList, f.Name())
		}
	}
	return pidList
}

func parseFile(path string) map[string]string {
	file, openFileError := os.Open(path)

	defer file.Close()
	if openFileError != nil {
		log.Fatal(openFileError)
	}

	procFile, parseError := procfs.ParseProcFile(file)
	if parseError != nil {
		log.Fatal(parseError)
	}
	return procFile
}
