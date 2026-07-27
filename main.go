package main

import "fmt"

type SimpleArchiver struct {
	inputPath  string
	outputPath string
	buffer     []byte
}

func NewArchiver(inputPath string) *SimpleArchiver {
	return &SimpleArchiver{
		inputPath: inputPath,
		buffer:    make([]byte, 1024*8),
	}
}

func (sa *SimpleArchiver) compressEmpty(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	return data
}

func main() {
	sa := NewArchiver("input.txt")
	fmt.Printf("Архиватор создан, размер буфера: %d байт\nпуть к исходному файлу: %s", cap(sa.buffer), sa.inputPath)
}
