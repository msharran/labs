package dns

import "fmt"

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
