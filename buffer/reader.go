package buffer

type ReaderManager struct {
	Buffers map[int32][]int32
}

func NewReaderManager() *ReaderManager {
	return &ReaderManager{
		Buffers: make(map[int32][]int32),
	}
}

func (m *ReaderManager) AppendBuffer(id int32, data []int32) []int32 {
	if _, ok := m.Buffers[id]; !ok {
		m.Buffers[id] = make([]int32, 0)
	}

	m.Buffers[id] = append(m.Buffers[id], data...)

	return m.Buffers[id]
}
func (m *ReaderManager) CloseBuffer(id int32) {
	delete(m.Buffers, id)
}
