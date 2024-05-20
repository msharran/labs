package dns

import (
	"fmt"
	"strings"
)

type Question struct {
	QName  string
	QType  uint16
	QClass uint16
}

func encodeDomain(domain string) ([]byte, error) {
	parts := strings.Split(string(domain), ".")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid domain name, missing dot in %s", domain)
	}

	b := make([]byte, 0, len(domain)+len(parts)+1)
	for _, p := range parts {
		b = append(b, byte(len(p)))
		b = append(b, []byte(p)...)
	}
	b = append(b, 0x00)
	return b, nil
}

// decodeDomain reads a domain name from buffer b
// Returns the domain name d, no. of bytes read n, and error if any
func decodeDomain(b []byte) (d string, n int, err error) {
	var dd []string
	for {
		count := int(b[n])
		n += 1 // read 1 byte
		if count == 0 {
			break
		}

		if count >= len(b) {
			return "", 0, fmt.Errorf("invalid Qname: count=%d >= len(qname=%d)", count, len(b))
		}

		dd = append(dd, string(b[n:n+count]))
		n += count // read `count` bytes
	}

	d = strings.Join(dd, ".")
	return
}
