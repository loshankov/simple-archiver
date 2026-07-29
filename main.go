package main

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

func (sa *SimpleArchiver) countRepeating(data []byte) []byte {
	data = sa.compressEmpty(data)
	if len(data) == 0 {
		return data
	}

	currentChar := data[0]
	count := 1
	var result []byte

	for i := 1; i < len(data); i++ {
		if currentChar == data[i] {
			count++
		} else {
			result = append(result, byte(count), currentChar)
			currentChar = data[i]
			count = 1
		}

	}

	result = append(result, byte(count), currentChar)

	return result
}

func main() {
	sa := NewArchiver("input.txt")
	//fmt.Printf("Архиватор создан, размер буфера: %d байт\nпуть к исходному файлу: %s", cap(sa.buffer), sa.inputPath)
	sa.countRepeating([]byte("AABBB"))
	sa.countRepeating([]byte("AABAB"))
}
