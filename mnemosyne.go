// +build linux

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/arekn/mnemosyne/procfs"
	"io/ioutil"
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
const defaultConfigFile = "config.json"

func main() {

	config, configError := loadConfig(defaultConfigFile)
	if configError != nil {
		log.Println(configError)
		log.Println("using default config")
	}

	createOutputFolder(config.OutputFolder)

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
				memory := measureAllMemory()
				saveMeasurementToFile(memory, config)
			}
		}
	}
}

type Memory struct {
	Total int
	Used  int
}

func loadConfig(configFile string) (Config, error) {
	defaultConfig := Config{OutputFolder: outputFolder, FilePrefix: "mnemosyne"}

	jsonFileConfig, openFileError := os.Open(configFile)
	if openFileError != nil {
		return defaultConfig, openFileError
	}
	defer jsonFileConfig.Close()

	dd, readAllError := ioutil.ReadAll(jsonFileConfig)
	if readAllError != nil {
		return defaultConfig, readAllError
	}

	unmarshalError := json.Unmarshal(dd, &defaultConfig)
	if unmarshalError != nil {
		return defaultConfig, unmarshalError
	}

	return defaultConfig, nil
}

func initStopChannel() chan os.Signal {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)
	return stop
}

func createOutputFolder(outputFolder string) {
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		log.Println("creating output folder")
		os.Mkdir(outputFolder, os.ModePerm)
	}
}

func saveMeasurementToFile(memory Memory, config Config) {

	timestamp := time.Now()
	filename := fmt.Sprintf("%v-%v.csv", config.FilePrefix, timestamp.Format(fileNameTimestamp))
	filePath := path.Join(config.OutputFolder, filename)
	dataFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer closeDataFile(dataFile)
	if err != nil {
		log.Println(err)
		return
	}
	csvWriter := csv.NewWriter(dataFile)
	writeError := csvWriter.Write([]string{timestamp.Format(fileEntryTimestamp), strconv.Itoa(memory.Used), strconv.Itoa(memory.Total)})
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

func measureAllMemory() Memory {
	memInfo := parseFile(path.Join(procPath, memInfoFile))

	memFree, memFreeError := memInfo.MemFree()
	if memFreeError != nil {
		log.Println(memFreeError)
	}
	memTotal, memTotalError := memInfo.MemTotal()
	if memTotalError != nil {
		log.Println(memTotalError)
	}

	usedMemory := memTotal - memFree
	return Memory{memTotal, usedMemory}
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
