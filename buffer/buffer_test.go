package buffer

import (
	"testing"
)

func TestManager(t *testing.T) {
	t.Run("should create and delete buffer", func(t *testing.T) {
		m := NewManager()
		_, err := m.TryNewBuffer(1)

		if err != nil {
			t.Error(err)
			return
		}
		if m.Count() != 1 {
			t.Error("count is not 1")
			return
		}

		m.CloseBuffer(1)
		if m.Count() != 0 {
			t.Error("count is not 0")
			return
		}
	})

	t.Run("should append to buffer", func(t *testing.T) {
		m := NewManager()
		b, _ := m.TryNewBuffer(1)

		if d := b.AppendBuffer([]int32{1}); len(d) != 1 || d[0] != 1 {
			t.Error("expected appended buffer")
		}
	})

	t.Run("should manage multiple buffers", func(t *testing.T) {
		m := NewManager()

		b1, _ := m.TryNewBuffer(1)
		b2, _ := m.TryNewBuffer(2)

		if d := b1.AppendBuffer([]int32{1}); len(d) != 1 || d[0] != 1 {
			t.Error("expected appended buffer 1")
		}
		if d := b2.AppendBuffer([]int32{2}); len(d) != 1 || d[0] != 2 {
			t.Error("expected appended buffer 2")
		}
	})

	t.Run("should return an id", func(t *testing.T) {
		m := NewManager()

		m.TryNewBuffer(1)
		m.TryNewBuffer(2)

		id, err := m.GetFirstID()
		if err != nil {
			t.Error(err)
			return
		}
		if id == -1 {
			t.Error("manager is empty")
			return
		}
	})

	t.Run("should return the same buffer if it exists", func(t *testing.T) {
		m := NewManager()

		b1, _ := m.TryNewBuffer(1)
		b1.AppendBuffer([]int32{1})

		b2 := m.NewBuffer(1)
		if b1.Buffer[0] != b2.Buffer[0] {
			t.Error("buffer data mismatch")
			return
		}
	})
}
