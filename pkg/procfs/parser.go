package procfs

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

const keyValueSeparator = ":"

func ParseProcFile(source io.Reader) (map[string]string, error) {

	scanner := bufio.NewScanner(source)

	result := make(map[string]string)

	for scanner.Scan() {

		if scannerError := scanner.Err(); scannerError != nil {
			return nil, scannerError
		}
		line := scanner.Text()

		index := strings.Index(line, ":")

		if index == -1 {
			message := fmt.Sprintf("Key - Value separator \"%s\" not found on line: %s", keyValueSeparator, line)
			return nil, errors.New(message)
		}

		key := line[:index]
		value := strings.TrimSpace(line[index+1:])

		result[key] = value
	}
	return result, nil
}
