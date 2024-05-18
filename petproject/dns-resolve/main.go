package main

import (
	"log"
	"net"
)

var GoogleDNSAddr = "8.8.8.8:53"

// https://codingchallenges.fyi/challenges/challenge-dns-resolver
func main() {
	// DNS header is 12 bytes (96 bits)
	h := Header{
		ID:      22,
		RD:      1,
		QDCOUNT: 1,
	}
	q := Question{
		QNAME:  "dns.google.com",
		QTYPE:  1,
		QCLASS: 1,
	}
	m := Message{
		Header:   h,
		Question: q,
	}

	// sort h+q in big endian byte order
	log.Printf("DNS message (hex): %x\n", m.BigEndianBytes())

	// dial Google's public DNS server (udp)
	// at 8.8.8.8 and port 53.
	// then send `h` and `q` as a single message.
	log.Printf("Dialing Google's public DNS server at %s\n", GoogleDNSAddr)
	conn, err := net.Dial("udp", GoogleDNSAddr)
	if err != nil {
		log.Fatalf("> Dial failed: %v\n", err)
	}
	defer conn.Close()

	// send the message
	log.Printf("> Sending message to Google's public DNS server\n")
	n, err := conn.Write(m.BigEndianBytes())
	if err != nil {
		log.Fatalf("> Failed to write: %v\n", err)
	}
	log.Printf("> Sent %d bytes\n", n)

	// read the response
	log.Printf("< Reading response from Google's public DNS server\n")
	buf := make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		log.Fatalf("< Failed to read: %v", err)
	}
	log.Printf("< Read %d bytes\n", n)

	resp := buf[:n]
	// convert the response to a hex string
	log.Printf("Response: %x\n", resp)

}
