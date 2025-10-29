package bytebuffer

import (
	"testing"
)

func TestBufferBytes(t *testing.T) {
	t.Run("should encode buffer to bytes", func(t *testing.T) {
		data := []int32{1, 2, 3}
		result, err := BufferToBytes(data)
		if err != nil {
			t.Error(err)
			return
		}

		expected := []byte{0x1, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0}
		if len(result) != len(expected) {
			t.Error("length mismatch")
			return
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("value mismatch: expected %x, got %x", v, result[i])
				return
			}
		}
	})
	t.Run("should encode bytes to buffer", func(t *testing.T) {
		data := []byte{0x1, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0}
		result, err := BytesToBuffer(data)
		if err != nil {
			t.Error(err)
			return
		}

		expected := []int32{1, 2, 3}
		if len(result) != len(expected) {
			t.Error("length mismatch")
			return
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("value mismatch: expected %d, got %d", v, result[i])
				return
			}
		}
	})
}
