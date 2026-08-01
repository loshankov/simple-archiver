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
	i := 0

	for i < len(data) {
		count := runLength(data, i)

		// сжатый блок
		if count >= 4 {
			result = append(result, byte(0x80|count), data[i])
			i += count
			continue
		}

		start := i
		i += count

		for i < len(data) {
			next := runLength(data, i)
			if next >= 3 {
				break // впереди серия — она уйдёт в свой блок
			}
			if (i - start + next) > 127 {
				break // лимит длины группы
			}
			i += next
		}

		result = append(result, byte(i-start))
		result = append(result, data[start:i]...)
	}

	return result
}

// runLength — длина серии одинаковых байт начиная с i, максимум 127
func runLength(data []byte, i int) int {
	c := 1
	for i+c < len(data) && data[i+c] == data[i] && c < 127 {
		c++
	}
	return c
}

func main() {
	sa := NewArchiver("input.txt")
	a := sa.compress([]byte("AAAAAB"))
	fmt.Println(a)
}
