package procfs

import (
	"errors"
	"strconv"
	"strings"
)

type MemInfoFile map[string]string

type MemInfoError error

func (m MemInfoFile) MemTotal() (int, MemInfoError) {
	memTotalKB := m["MemTotal"]
	return kbField(memTotalKB)
}

func (m MemInfoFile) MemFree() (int, MemInfoError) {
	memFreeKB := m["MemFree"]
	return kbField(memFreeKB)
}

func kbField(field string) (int, error) {
	fieldSplit := strings.Fields(field)
	if len(fieldSplit) != 2 || fieldSplit[1] != "kB" {
		return -1, errors.New("Field is not kB type: " + field)
	}
	value, e := strconv.Atoi(fieldSplit[0])
	if e != nil {
		return -1, e
	}
	return value, nil
}
