package buffer

import "sync"

type LemonadeBufferSet uint

const (
	BUFFER_PARTIAL LemonadeBufferSet = 0
	BUFFER_END     LemonadeBufferSet = 1
)

type WriterManager struct {
	Buffers map[int32][]*WriterBuffer
	mutex   sync.Mutex
}

func NewWriterManager() *WriterManager {
	return &WriterManager{
		Buffers: make(map[int32][]*WriterBuffer),
	}
}

func (m *WriterManager) GetFirstID() (int32, bool) {
	for id := range m.Buffers {
		return id, true
	}
	return -1, false
}

func (m *WriterManager) Queue(appid int32, buffer *WriterBuffer) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.Buffers[appid]; !ok {
		m.Buffers[appid] = make([]*WriterBuffer, 0)
	}

	m.Buffers[appid] = append(m.Buffers[appid], buffer)
}
func (m *WriterManager) Dequeue(appid int32) *WriterBuffer {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.Buffers[appid]) == 0 {
		return nil
	}

	buffer := m.Buffers[appid][0]
	m.Buffers[appid] = m.Buffers[appid][1:]

	if len(m.Buffers[appid]) == 0 {
		delete(m.Buffers, appid)
	}

	return buffer
}

type WriterBuffer struct {
	Buffer []int32
	Set    LemonadeBufferSet
}
