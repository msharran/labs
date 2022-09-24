package store

type Reader interface {
	Read(key string) ([]byte, error)
}

type Writer interface {
	Write(key string, data []byte) error
}

type Store interface {
	Reader
	Writer
}
