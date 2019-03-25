// +build linux

package main

import (
	"encoding/csv"
	"fmt"
	"github.com/arekn/mnemosyne/pkg/procfs"
	"log"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"
)

const outputFolder = "output"
const procPath = "/proc"
const memInfoFile = "meminfo"
const fileNameTimestamp = "2006.01.02"
const fileEntryTimestamp = "15:04"

func main() {

	createOutputFolder()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	oneMinuteTicker := time.NewTicker(1 * time.Minute)
	defer oneMinuteTicker.Stop()

	for {
		select {
		case <-stop:
			{
				log.Println("stopping application")
				os.Exit(0)
			}
		case <-oneMinuteTicker.C:
			{
				measureMemoryUsage()
			}
		}
	}
}

func createOutputFolder() {
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		log.Println("creating output folder")
		os.Mkdir(outputFolder, os.ModePerm)
	}
}

func measureMemoryUsage() {
	memTotal, memUsed := checkMemoryState()

	timestamp := time.Now()
	filename := fmt.Sprintf("meminfo-%v.csv", timestamp.Format(fileNameTimestamp))
	filePath := path.Join(outputFolder, filename)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(f)
	writeError := w.Write([]string{timestamp.Format(fileEntryTimestamp), strconv.Itoa(memUsed), strconv.Itoa(memTotal)})
	if writeError != nil {
		log.Fatal(writeError)
	}
	w.Flush()
	closeError := f.Close()
	if closeError != nil {
		log.Fatal(closeError)
	}
}

func checkMemoryState() (total int, used int) {
	memInfo := procfs.MemInfoFile(parseFile(path.Join(procPath, memInfoFile)))

	memFree, memFreeError := memInfo.MemFree()
	if memFreeError != nil {
		log.Fatal(memFreeError)
	}
	memTotal, memTotalError := memInfo.MemTotal()
	if memTotalError != nil {
		log.Fatal(memTotalError)
	}

	usedMemory := memTotal - memFree
	return memTotal, usedMemory
}

func parseFile(path string) map[string]string {
	file, openFileError := os.Open(path)

	if openFileError != nil {
		log.Fatal(openFileError)
	}

	procFile, parseError := procfs.ParseProcFile(file)
	file.Close()
	if parseError != nil {
		log.Fatal(parseError)
	}
	return procFile
}
