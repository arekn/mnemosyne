package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseMemInfoFile(fileLocation string) {

	file, e := os.Open(fileLocation)

	if e != nil {
		log.Fatal(e)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	mem := make(map[string]int)

	for scanner.Scan() {

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
		line := scanner.Text()

		fields := strings.Fields(line)

		key := strings.TrimSuffix(fields[0], ":")
		value, conversionError := strconv.Atoi(fields[1])

		if conversionError != nil {
			fmt.Println(conversionError)
		}

		mem[key] = value
	}
}

// https://matthias-endler.de/2018/go-io-testing/
