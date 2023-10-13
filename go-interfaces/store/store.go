package store

type Reader interface {
	Read(key string) ([]byte, error)
}

type Writer interface {
	Write(key string, data []byte) error
}

type StringWriter interface {
	WriteString(key string, data string) error
}

type Store interface {
	Reader
	Writer
	StringWriter
}
