package buffer

import (
	"git.jaezmien.com/Jaezmien/lemonade-stand/chunk"
)

const (
	MAXIMUM_BUFFER_LENGTH = 29 - 3
)

func SplitBuffer(buffer []int32) [][]int32 {
	s := make([][]int32, 0)

	if len(buffer) <= MAXIMUM_BUFFER_LENGTH {
		s = append(s, buffer)
	} else {
		chunk.ChunkSlice(buffer, MAXIMUM_BUFFER_LENGTH, func(partialSlice []int32, _ bool) {
			s = append(s, partialSlice)
		})
	}

	return s
}
