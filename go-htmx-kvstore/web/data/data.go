package data

type KeyValue struct {
	Key   string
	Value string
}

type KeyValues []*KeyValue

type User struct {
	Username string
	Password string
	Token    []byte
}

type Users []*User
