package s3

type S3 struct {
	mocks3 map[string][]byte
}

func NewStore() *S3 {
	return &S3{mocks3: make(map[string][]byte)}
}

func (m *S3) Read(key string) ([]byte, error) {
	return m.mocks3[key], nil
}

func (m *S3) Write(key string, data []byte) error {
	m.mocks3[key] = data
	return nil
}
