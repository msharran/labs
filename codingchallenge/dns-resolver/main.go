package main

import (
	"dns-resolver/internal/dns"
	"log"
	"net"
)

var GoogleDNSAddr = "8.8.8.8:53"

// https://codingchallenges.fyi/challenges/challenge-dns-resolver
func main() {
	// DNS header is 12 bytes (96 bits)
	m := dns.Message{
		Header: dns.Header{
			ID:      22,
			Flag:    dns.NewHeaderFlag(0, 0, 0, 0, 1, 0, 0, 0), // set RD (recursion desired) to 1
			QdCount: 1,
		},
		Question: dns.Question{
			QName:  "dns.google.com",
			QType:  1,
			QClass: 1,
		},
	}

	log.Printf("DNS request message: \n%s\n", m)

	// dial Google's public DNS server (udp)
	// at 8.8.8.8 and port 53.
	// then send `h` and `q` as a single message.
	log.Printf("Dialing Google's public DNS server at %s\n", GoogleDNSAddr)
	conn, err := net.Dial("udp", GoogleDNSAddr)
	if err != nil {
		log.Fatalf("> Dial failed: %v\n", err)
	}
	defer conn.Close()

	// convert message to DNS's binary format
	log.Printf("> Sending message to Google's public DNS server\n")
	buf, err := dns.MarshalMessage(m)
	if err != nil {
		log.Fatal(err)
	}

	// send the message
	n, err := conn.Write(buf)
	if err != nil {
		log.Fatalf("> Failed to write: %v\n", err)
	}
	log.Printf("> Sent %d bytes\n", n)

	// read the response
	log.Printf("< Reading response from Google's public DNS server\n")
	buf = make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		log.Fatalf("< Failed to read: %v", err)
	}
	log.Printf("< Read %d bytes\n", n)
	resp := buf[:n]

	// parse the response message
	log.Printf("Response (hex): %x\n", resp)
	respMsg := dns.Message{}
	if err := dns.UnmarshalMessage(resp, &respMsg); err != nil {
		log.Fatal(err)
	}

	log.Printf("DNS response message: \n%s\n", respMsg)
}
