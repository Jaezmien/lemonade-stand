package buffer

import (
	"errors"
	"sync"

	"git.jaezmien.com/Jaezmien/lemonade-stand/chunk"
	"git.jaezmien.com/Jaezmien/lemonade-stand/encoder"
)

var ErrBufferExists = errors.New("buffer already exists for id")
var ErrBufferNotExists = errors.New("buffer doesn't exist for id")
var ErrManagerEmpty = errors.New("manager is empty")

const (
	MAXIMUM_BUFFER_LENGTH = 29 - 3
)

type LemonadeBufferSet uint

const (
	BUFFER_PARTIAL LemonadeBufferSet = 0
	BUFFER_END     LemonadeBufferSet = 1
)

// --- //

type LemonadeBufferManager struct {
	Buffers map[int32]*LemonadeBuffer
	m       sync.Mutex
}

func NewManager() *LemonadeBufferManager {
	return &LemonadeBufferManager{
		Buffers: make(map[int32]*LemonadeBuffer),
	}
}

func (m *LemonadeBufferManager) Count() int {
	return len(m.Buffers)
}

func (m *LemonadeBufferManager) GetFirstID() (int32, error) {
	for id, _ := range m.Buffers {
		return id, nil
	}

	return -1, ErrManagerEmpty
}

// --- //

type LemonadeBuffer struct {
	Buffer []int32
	Set    LemonadeBufferSet
}

func NewBuffer() *LemonadeBuffer {
	return &LemonadeBuffer{
		Buffer: make([]int32, 0),
		Set:    BUFFER_END,
	}
}

func (b *LemonadeBuffer) AppendBuffer(data []int32) []int32 {
	b.Buffer = append(b.Buffer, data...)
	return b.Buffer
}
func (b *LemonadeBuffer) DecodeToString() (string, error) {
	return encoder.BufferToString(b.Buffer)
}

// --- //

func (m *LemonadeBufferManager) HasBuffer(id int32) bool {
	_, ok := m.Buffers[id]
	return ok
}
func (m *LemonadeBufferManager) TryNewBuffer(id int32) (*LemonadeBuffer, error) {
	if m.HasBuffer(id) {
		return nil, ErrBufferExists
	}

	m.m.Lock()
	defer m.m.Unlock()

	m.Buffers[id] = NewBuffer()
	return m.Buffers[id], nil
}
func (m *LemonadeBufferManager) NewBuffer(id int32) *LemonadeBuffer {
	if !m.HasBuffer(id) {
		m.TryNewBuffer(id)
	}

	b, ok := m.Buffers[id]
	if !ok {
		panic("trynewbuffer failed")
	}

	return b
}

func (m *LemonadeBufferManager) TryGetBuffer(id int32) (*LemonadeBuffer, error) {
	b, ok := m.Buffers[id]
	if !ok {
		return nil, ErrBufferNotExists
	}

	return b, nil
}
func (m *LemonadeBufferManager) CloseBuffer(id int32) {
	if !m.HasBuffer(id) {
		return
	}

	m.m.Lock()
	defer m.m.Unlock()

	delete(m.Buffers, id)
}

func SplitBuffer(buffer []int32) [][]int32 {
	s := make([][]int32, 0)

	if len(buffer) <= MAXIMUM_BUFFER_LENGTH {
		s = append(s, buffer)
	} else {
		chunk.ChunkSlice(buffer, MAXIMUM_BUFFER_LENGTH, func(partialSlice []int32, _ bool) {
			s = append(s, buffer)
		})
	}

	return s
}
