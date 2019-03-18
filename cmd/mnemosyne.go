// +build linux

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const meminfo = "/proc/meminfo"

func main() {
	file, e := os.Open(meminfo)

	if e != nil {
		log.Fatal(e)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		fmt.Println(scanner.Text())

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	}

}
