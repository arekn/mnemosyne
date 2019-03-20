package memory

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
)

type MemInfo map[string]int

func ParseMemInfo(memInfoSource io.Reader) (MemInfo, error) {

	scanner := bufio.NewScanner(memInfoSource)

	result := make(MemInfo)

	for scanner.Scan() {

		if scannerError := scanner.Err(); scannerError != nil {
			log.Println("memInfoSource scanner error")
			return nil, scannerError
		}
		line := scanner.Text()

		fields := strings.Fields(line)

		if len(fields) < 2 {

		}
		key := strings.TrimSuffix(fields[0], ":")
		value, conversionError := strconv.Atoi(fields[1])

		if conversionError != nil {
			log.Printf("converting value: %v", fields[1])
			return nil, conversionError
		}

		result[key] = value
	}
	return result, nil
}
