// +build linux

package main

import (
	"encoding/csv"
	"fmt"
	"github.com/arekn/mnemosyne/procfs"
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

	stop := initStopChannel()

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

func initStopChannel() chan os.Signal {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)
	return stop
}

func createOutputFolder() {
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		log.Println("creating output folder")
		os.Mkdir(outputFolder, os.ModePerm)
	}
}

func measureMemoryUsage() {
	memTotal, memUsed := measureAllMemory()

	timestamp := time.Now()
	filename := fmt.Sprintf("meminfo-%v.csv", timestamp.Format(fileNameTimestamp))
	filePath := path.Join(outputFolder, filename)
	dataFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer closeDataFile(dataFile)
	if err != nil {
		log.Println(err)
		return
	}
	csvWriter := csv.NewWriter(dataFile)
	writeError := csvWriter.Write([]string{timestamp.Format(fileEntryTimestamp), strconv.Itoa(memUsed), strconv.Itoa(memTotal)})
	if writeError != nil {
		log.Println(writeError)
		return
	}
	csvWriter.Flush()
}

func closeDataFile(dataFile *os.File) {
	closeError := dataFile.Close()
	if closeError != nil {
		log.Println(closeError)
	}
}

func measureAllMemory() (total int, used int) {
	memInfo := procfs.MemInfoFile(parseFile(path.Join(procPath, memInfoFile)))

	memFree, memFreeError := memInfo.MemFree()
	if memFreeError != nil {
		log.Println(memFreeError)
	}
	memTotal, memTotalError := memInfo.MemTotal()
	if memTotalError != nil {
		log.Println(memTotalError)
	}

	usedMemory := memTotal - memFree
	return memTotal, usedMemory
}

func parseFile(path string) procfs.MemInfoFile {
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
