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
	if len(data) == 0 {
		return []byte{}
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

func (sa *SimpleArchiver) createControlByte(count int, isCompressed bool) byte {
	if count > 127 {
		count = 127
	}

	if isCompressed {
		return byte(128 + count)
	}
	return byte(count)
}

func (sa *SimpleArchiver) groupCompress(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	var result, hand []byte
	countRepeat := 1
	currentByte := data[0]

	flushHand := func() {
		if len(hand) != 0 {
			result = append(result, sa.createControlByte(len(hand), false))
			result = append(result, hand...)
			hand = hand[:0]
		}
	}

	closeRun := func(end int) {
		if countRepeat >= 4 {
			flushHand()
			result = append(result, sa.createControlByte(countRepeat, true), currentByte)
		} else {
			hand = append(hand, data[end-countRepeat:end]...)
		}
	}

	for i := 1; i < len(data); i++ {
		if currentByte == data[i] {
			countRepeat++
		} else {
			closeRun(i)
			currentByte = data[i]
			countRepeat = 1
		}
	}
	closeRun(len(data))
	flushHand()

	return result
}

func main() {
	sa := NewArchiver("input.txt")
	sa.groupCompress([]byte("AAAABCD"))
	sa.groupCompress([]byte("ABCCDE"))
}
