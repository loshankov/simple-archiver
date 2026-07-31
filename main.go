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

const (
	minRun  = 4
	stopRun = 3
	maxLen  = 127
)

func (sa *SimpleArchiver) groupCompress(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	var result, hand []byte

	runLength := func(i int) int {
		n := 1
		for i+n < len(data) && data[i+n] == data[i] {
			n++
		}
		return n
	}

	flushHand := func() {
		for len(hand) > 0 {
			n := min(len(hand), maxLen)
			result = append(result, sa.createControlByte(n, false))
			result = append(result, hand[:n]...)
			hand = hand[n:]
		}
	}

	for i := 0; i < len(data); {
		run := runLength(i)

		if run >= minRun {
			flushHand()
			for r := run; r > 0; {
				n := min(r, maxLen)
				result = append(result, sa.createControlByte(n, true), data[i])
				r -= n
			}
			i += run
			continue
		}

		hand = append(hand, data[i])
		i++

		if i < len(data) && runLength(i) >= stopRun {
			flushHand()
		}

		if len(hand) >= maxLen {
			flushHand()
		}
	}

	flushHand()
	return result
}

func main() {
	sa := NewArchiver("input.txt")
	sa.groupCompress([]byte("AAAABCD"))
	sa.groupCompress([]byte("ABCCDE"))
}
