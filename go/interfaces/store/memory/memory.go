package memory

type Memory struct {
	data map[string][]byte
}

func NewStore() *Memory {
	return &Memory{data: make(map[string][]byte)}
}

func (m *Memory) Read(key string) ([]byte, error) {
	return m.data[key], nil
}

func (m *Memory) Write(key string, data []byte) error {
	m.data[key] = data
	return nil
}
