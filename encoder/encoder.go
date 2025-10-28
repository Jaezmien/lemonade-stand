package encoder

import (
	"errors"
	"strings"
	"unicode"
)

func stringSatifiesEncoder(str string) bool {
	for _, char := range str {
		if char > unicode.MaxASCII {
			return false
		}
		if char < 0 {
			return false
		}
	}

	return true
}

func StringToBuffer(str string) ([]int32, error) {
	if !stringSatifiesEncoder(str) {
		return nil, errors.New("invalid string")
	}

	buffer := make([]int32, len(str))

	for idx, char := range str {
		buffer[idx] = int32(char)
	}

	return buffer, nil
}

func BufferToString(buffer []int32) (string, error) {
	var sb strings.Builder

	for _, idx := range buffer {
		err := sb.WriteByte(byte(idx))

		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}
