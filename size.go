package main

import (
	"fmt"
)

type Size int64

func (s Size) String() string {
	var gigabyte int64
	var megabyte int64
	var kilobyte int64
	var bytes int64

	bytes = int64(s)

	// bytes
	if bytes < 1024 {
		return fmt.Sprintf("%d byte", bytes)
	}

	// kilobytes
	kilobyte = bytes / 1024
	bytes %= 1024
	if kilobyte < 1024 {
		return fmt.Sprintf("%d.%d Kbyte", kilobyte, bytes)
	}

	// megabytes
	megabyte = kilobyte / 1024
	kilobyte %= 1024
	if megabyte < 1024 {
		return fmt.Sprintf("%d.%d Mbyte", megabyte, kilobyte)
	}

	// gigabytes
	gigabyte = megabyte / 1024
	megabyte %= 1024
	return fmt.Sprintf("%d.%d Gbyte", gigabyte, megabyte)
}
