package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Message struct {
	Header   Header
	Question Question
}

func MarshalMessage(m Message) ([]byte, error) {
	h := m.Header
	q := m.Question

	buf := new(bytes.Buffer)
	// write the header
	binary.Write(buf, binary.BigEndian, h.ID)
	binary.Write(buf, binary.BigEndian, h.Flag)
	binary.Write(buf, binary.BigEndian, h.QdCount)
	binary.Write(buf, binary.BigEndian, h.AnCount)
	binary.Write(buf, binary.BigEndian, h.NsCount)
	binary.Write(buf, binary.BigEndian, h.ArCount)

	// write the question
	dom, err := encodeDomain(q.QName)
	if err != nil {
		return nil, err
	}
	binary.Write(buf, binary.BigEndian, dom)
	binary.Write(buf, binary.BigEndian, q.QType)
	binary.Write(buf, binary.BigEndian, q.QClass)

	return buf.Bytes(), nil
}

func UnmarshalMessage(b []byte, m *Message) (err error) {
	if m == nil {
		return fmt.Errorf("nil message")
	}

	// parse header
	h := parseHeader(b[:12])

	// parse question
	q, _, err := parseQuestion(b[12:])
	if err != nil {
		return
	}

	m.Header = h
	m.Question = q
	return nil
}

func (m Message) String() (s string) {
	h := m.Header
	q := m.Question

	s += fmt.Sprintf("Header:\n")
	s += fmt.Sprintf("	ID: %d\n", h.ID)
	s += fmt.Sprintf("	Flag: %s\n", h.Flag)
	s += fmt.Sprintf("	QdCount: %d\n", h.QdCount)
	s += fmt.Sprintf("	AnCount: %d\n", h.AnCount)
	s += fmt.Sprintf("	NsCount: %d\n", h.NsCount)
	s += fmt.Sprintf("	ArCount: %d\n", h.ArCount)
	s += fmt.Sprintf("Question:\n")
	s += fmt.Sprintf("	QName: %s\n", q.QName)
	s += fmt.Sprintf("	QType: %d\n", q.QType)
	s += fmt.Sprintf("	QClass: %d\n", q.QClass)
	return s
}

func parseHeader(b []byte) (h Header) {
	h.ID = binary.BigEndian.Uint16(b[:2])
	h.Flag = HeaderFlag(binary.BigEndian.Uint16(b[2:4]))
	h.QdCount = binary.BigEndian.Uint16(b[4:6])
	h.AnCount = binary.BigEndian.Uint16(b[6:8])
	h.NsCount = binary.BigEndian.Uint16(b[8:10])
	h.ArCount = binary.BigEndian.Uint16(b[10:12])
	return
}

func parseQuestion(b []byte) (q Question, n int, err error) {
	d, n, err := decodeDomain(b)
	if err != nil {
		return
	}
	q.QName = d
	q.QType = binary.BigEndian.Uint16(b[n : n+2])
	q.QClass = binary.BigEndian.Uint16(b[n+2 : n+4])
	n += 4
	return
}
