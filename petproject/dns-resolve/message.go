package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type Message struct {
	Header   Header
	Question Question
}

func (m Message) BigEndianBytes() []byte {
	h := m.Header
	q := m.Question

	buf := new(bytes.Buffer)
	// write the header
	binary.Write(buf, binary.BigEndian, h.ID)
	binary.Write(buf, binary.BigEndian, h.Flags())
	binary.Write(buf, binary.BigEndian, h.QDCOUNT)
	binary.Write(buf, binary.BigEndian, h.ANCOUNT)
	binary.Write(buf, binary.BigEndian, h.NSCOUNT)
	binary.Write(buf, binary.BigEndian, h.ARCOUNT)

	// write the question
	binary.Write(buf, binary.BigEndian, q.QNAME.Encode())
	binary.Write(buf, binary.BigEndian, q.QTYPE)
	binary.Write(buf, binary.BigEndian, q.QCLASS)

	return buf.Bytes()
}

func (m Message) String() string {
	h := m.Header
	q := m.Question

	b := bytes.NewBuffer(make([]byte, 12))
	fmt.Fprintf(b, "%04x", h.ID)
	fmt.Fprintf(b, "%04x", h.Flags())
	fmt.Fprintf(b, "%04x", h.QDCOUNT)
	fmt.Fprintf(b, "%04x", h.ANCOUNT)
	fmt.Fprintf(b, "%04x", h.NSCOUNT)
	fmt.Fprintf(b, "%04x", h.ARCOUNT)
	fmt.Fprintf(b, "%s", q.QNAME.Encode())
	fmt.Fprintf(b, "%04x", q.QTYPE)
	fmt.Fprintf(b, "%04x", q.QCLASS)

	return b.String()
}

// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1
type Header struct {
	// 16 bit
	ID uint16

	// flags is 16 bits
	QR     int
	Opcode int
	AA     int
	TC     int
	RD     int // RD is the recursion desired flag
	RA     int
	Z      int
	RCODE  int

	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func (h Header) Flags() uint16 {
	var flags uint16
	flags |= uint16(h.QR) << 15
	flags |= uint16(h.Opcode) << 11
	flags |= uint16(h.AA) << 10
	flags |= uint16(h.TC) << 9
	flags |= uint16(h.RD) << 8
	flags |= uint16(h.RA) << 7
	flags |= uint16(h.Z) << 6
	flags |= uint16(h.RCODE) << 0

	return flags
}

type Qname string

func (q Qname) Encode() []byte {
	parts := strings.Split(string(q), ".")
	if len(parts) == 0 {
		return nil
	}

	b := make([]byte, 0, len(q)+len(parts)+1)
	for _, p := range parts {
		b = append(b, byte(len(p)))
		b = append(b, []byte(p)...)
	}
	b = append(b, 0x00)
	return b
}

type Question struct {
	QNAME  Qname
	QTYPE  uint16
	QCLASS uint16
}
