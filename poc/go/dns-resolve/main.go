package main

import (
	"encoding/hex"
	"strings"
)

type message struct {
	header    header
	questions []question
}

type header struct {
	ID     uint16
	QR     uint8
	Opcode uint8
	AA     uint8
	TC     uint8
	// RD is the recursion desired flag
	RD      uint8
	RA      uint8
	Z       uint8
	RCODE   uint8
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func (h header) hex() []byte {

}

type qname string

func (q qname) encode() []byte {
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

type question struct {
	QNAME  qname
	QTYPE  uint16
	QCLASS uint16
}

// https://codingchallenges.fyi/challenges/challenge-dns-resolver
func main() {
	h := header{
		ID:      22,
		RD:      1,
		QDCOUNT: 1,
	}
}

func hexEncode(b []byte) string {
	return hex.EncodeToString(b)
}
