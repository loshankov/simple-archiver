package main

import "fmt"

const (
	minRun  = 4
	stopRun = 3
	maxLen  = 127
)

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

func (sa *SimpleArchiver) compress(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	var result []byte
	var hand []byte
	countRepeating := sa.countRepeating(data)
	for i := 0; i < len(countRepeating); {
		if countRepeating[i] >= 4 {
			result = append(result, sa.createControlByte(int(countRepeating[i]), true), countRepeating[i+1])
		} else {
			for range countRepeating[i] {
				hand = append(hand, countRepeating[i+1])
			}
		}
		if len(hand) == 0 {
			i += 2
			continue
		}
		if i+2 >= len(countRepeating) || countRepeating[i+2] >= 3 {
			result = append(result, sa.createControlByte(len(hand), false))
			result = append(result, hand...)
			hand = []byte{}
		}
		i += 2
	}

	return result
}

func (sa *SimpleArchiver) decompress(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	var result []byte

	for i := 0; i < len(data); {
		fmt.Printf("%#x\n", data[i])
		messageInfo := ""
		length := int(data[i] & 0x7F)
		if data[i]&0x80 != 0 {
			messageInfo += "cжатая, "
			i := i + 1
			value := data[i]
			for range length {
				result = append(result, value)
			}
		} else {
			messageInfo += "несжатая, "
		}
		fmt.Printf("%sдлина: %v\n", messageInfo, length)
		if data[i]&0x80 == 0 {
			i += 1 + length
		} else {
			i += 1
		}
	}
	return result
}

func main() {
	sa := NewArchiver("input.txt")
	//a := sa.compress([]byte("ABBBCCCCDE"))
	//a := sa.compress([]byte("AAAAAAA"))
	a := sa.decompress([]byte{0x85, 0x41, 0x03, 0x42, 0x43, 0x44})
	fmt.Println(a)
}
