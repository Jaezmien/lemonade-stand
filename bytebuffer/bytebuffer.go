package bytebuffer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func BufferToBytes(data []int32) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, v := range data {
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
func BytesToBuffer(data []byte) ([]int32, error) {
	result := make([]int32, 0)
	reader := bytes.NewReader(data)

	// FIXME: Is this safe?
	size := 4
	buff := make([]byte, uint(size))
	for {
		read, err := reader.Read(buff)
		if read > 0 {
			value := binary.LittleEndian.Uint32(buff[:size])
			result = append(result, int32(value))
		}

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
