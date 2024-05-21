package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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

type decoder struct {
	ptr  int
	hdr  *Header
	qn   *Question
	anss []*Answer
}

func (d decoder) decode(b []byte) (*Message, error) {
	// parse header
	d.parseHeader(b)

	// parse question
	err := d.parseQuestion(b)
	if err != nil {
		return nil, fmt.Errorf("parseQuestion error: %w", err)
	}

	err = d.parseAnswer(b)
	if err != nil {
		return nil, fmt.Errorf("parseAnswer error: %w", err)
	}

	return &Message{
		Header:   *d.hdr,
		Question: *d.qn,
		Answer:   *d.anss[0],
	}, nil
}

var defaultDecoder = decoder{}

func Decode(b []byte) (*Message, error) {
	return defaultDecoder.decode(b)
}

type Message struct {
	Header   Header
	Question Question
	Answer   Answer
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
	s += fmt.Sprintf("Answer:\n")
	s += fmt.Sprintf("	Name: %s\n", m.Answer.Name)
	s += fmt.Sprintf("	Type: %d\n", m.Answer.Type)
	s += fmt.Sprintf("	Class: %d\n", m.Answer.Class)
	s += fmt.Sprintf("	TTL: %d\n", m.Answer.TTL)
	s += fmt.Sprintf("	RDLength: %d\n", m.Answer.RDLength)
	s += fmt.Sprintf("	RData: %s\n", m.Answer.RData)
	return s
}

func (d *decoder) parseHeader(b []byte) {
	d.hdr = new(Header)
	d.hdr.ID = binary.BigEndian.Uint16(b[:2])
	d.hdr.Flag = HeaderFlag(binary.BigEndian.Uint16(b[2:4]))
	d.hdr.QdCount = binary.BigEndian.Uint16(b[4:6])
	d.hdr.AnCount = binary.BigEndian.Uint16(b[6:8])
	d.hdr.NsCount = binary.BigEndian.Uint16(b[8:10])
	d.hdr.ArCount = binary.BigEndian.Uint16(b[10:12])
	d.ptr = 12
}

func (d *decoder) parseQuestion(b []byte) error {
	n, ptr, err := decodeDomain(b, d.ptr)
	if err != nil {
		return fmt.Errorf("parseQuestion error: %w", err)
	}
	d.qn = new(Question)
	d.qn.QName = n
	d.qn.QType = binary.BigEndian.Uint16(b[ptr : ptr+2])
	d.qn.QClass = binary.BigEndian.Uint16(b[ptr+2 : ptr+4])
	d.ptr += 4 // 2 bytes for QType and 2 bytes for QCass
	return nil
}

func (d *decoder) parseAnswer(b []byte) error {
	n, ptr, err := decodeDomain(b, d.ptr)
	if err != nil {
		return fmt.Errorf("parseAnswer error: %w", err)
	}
	d.anss = append(d.anss, new(Answer))
	a := d.anss[0]
	a.Name = n
	a.Type = binary.BigEndian.Uint16(b[ptr : ptr+2])
	a.Class = binary.BigEndian.Uint16(b[ptr+2 : ptr+4])
	a.TTL = binary.BigEndian.Uint32(b[ptr+4 : ptr+8])
	a.RDLength = binary.BigEndian.Uint16(b[ptr+8 : ptr+10])
	a.RData = string(b[ptr+10 : ptr+10+int(a.RDLength)])
	d.ptr += 10 + int(a.RDLength)
	return nil
}

// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1
type Header struct {
	ID      uint16     // id is 16 bits
	Flag    HeaderFlag // flag is 16 bits
	QdCount uint16     // qdcount is 16 bits
	AnCount uint16     // ancount is 16 bits
	NsCount uint16     // nscount is 16 bits
	ArCount uint16     // arcount is 16 bits
}

type HeaderFlag uint16

func NewHeaderFlag(qr, opcode, aa, tc, rd, ra, z, rcode int) HeaderFlag {
	var f HeaderFlag
	f |= HeaderFlag(qr) << 15
	f |= HeaderFlag(opcode) << 11
	f |= HeaderFlag(aa) << 10
	f |= HeaderFlag(tc) << 9
	f |= HeaderFlag(rd) << 8
	f |= HeaderFlag(ra) << 7
	f |= HeaderFlag(z) << 6
	f |= HeaderFlag(rcode) << 0
	return f
}

func (f HeaderFlag) Parts() (qr, opcode, aa, tc, rd, ra, z, rcode int) {
	// qr is the first bit
	qr = int(f >> 15 & 0x1)     // 0x1 = 0000 0000 0000 0001
	opcode = int(f >> 11 & 0xf) // 0xf = 0000 0000 0000 1111
	aa = int(f >> 10 & 0x1)     // 0x1 = 0000 0000 0000 0001
	tc = int(f >> 9 & 0x1)      // 0x1 = 0000 0000 0000 0001
	rd = int(f >> 8 & 0x1)      // 0x1 = 0000 0000 0000 0001
	ra = int(f >> 7 & 0x1)      // 0x1 = 0000 0000 0000 0001
	z = int(f >> 6 & 0x1)       // 0x1 = 0000 0000 0000 0001
	rcode = int(f >> 0 & 0xf)   // 0xf = 0000 0000 0000 1111
	return
}

func (f HeaderFlag) String() string {
	qr, opcode, aa, tc, rd, ra, z, rcode := f.Parts()
	return fmt.Sprintf("QR=%d, OPCODE=%d, AA=%d, TC=%d, RD=%d, RA=%d, Z=%d, RCODE=%d", qr, opcode, aa, tc, rd, ra, z, rcode)
}

type Question struct {
	QName  string
	QType  uint16
	QClass uint16
}

type Answer struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    string
}
