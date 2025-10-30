package buffer

import (
	"testing"
)

func TestWriterManager(t *testing.T) {
	m := NewWriterManager()

	t.Run("should create and delete buffer", func(t *testing.T) {
		m.Queue(1, &WriterBuffer{
			Buffer: []int32{1, 2, 3},
			Set: BUFFER_END,
		})

		b := m.Dequeue(1)
		if len(b.Buffer) != 3 || b.Set != BUFFER_END {
			t.Error("invalid buffer")
			return
		}
	})

	t.Run("should return an id", func(t *testing.T) {
		m.Queue(1, &WriterBuffer{
			Buffer: []int32{1, 2, 3},
			Set: BUFFER_END,
		})
		m.Queue(2, &WriterBuffer{
			Buffer: []int32{1, 2, 3},
			Set: BUFFER_END,
		})

		id, ok := m.GetFirstID()
		if !ok {
			t.Error("got no id")
			return
		}
		if id == -1 {
			t.Error("no id")
			return
		}

		m.Dequeue(1)
		m.Dequeue(2)
	})

	t.Run("should return no id", func(t *testing.T) {
		id, ok := m.GetFirstID()
		if ok || id != -1 {
			t.Error("got id")
			return
		}
	})

	t.Run("should return no buffer", func(t *testing.T) {
		b := m.Dequeue(2)
		if b != nil {
			t.Error("got buffer")
			return
		}
	})
}
func TestReaderManager(t *testing.T) {
	m := NewReaderManager()
	
	t.Run("should append and clear buffer", func(t *testing.T) {
		m.AppendBuffer(1, []int32{1, 2, 3})
		m.CloseBuffer(1)
	})
}

func TestSplitter(t *testing.T) {
	t.Run("should split once", func(t *testing.T) {
		b := []int32{1, 2, 3, 4}

		split := SplitBuffer(b)
		if len(split) != 1 {
			t.Error("buffer isn't split")
			return
		}
		if len(split[0]) != 4 {
			t.Error("buffer slice has missing element")
			return
		}
	})

	t.Run("should split twice", func(t *testing.T) {
		var b []int32
		for i := range 32 {
			b = append(b, int32(i))
		}

		split := SplitBuffer(b)
		if len(split) != 2 {
			t.Error("unexpected split length")
			return
		}
		if len(split[0]) != MAXIMUM_BUFFER_LENGTH {
			t.Error("buffer slice has missing element")
			return
		}
		if len(split[1]) > MAXIMUM_BUFFER_LENGTH {
			t.Error("buffer slice overflowing count")
			return
		}
	})
}
