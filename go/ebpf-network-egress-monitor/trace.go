package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
	"unsafe"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

// Generate Go bindings for egress_bpf.c.
// Run this after changing egress_bpf.c:
//
//	go generate ./...
//
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -g" egress egress_bpf.c -- -I. -I/usr/include/x86_64-linux-gnu

const (
	eventConnect     = 1
	eventDNSQuery    = 2
	eventDNSResponse = 3
)

type event struct {
	Type      uint32
	PID       uint32
	TGID      uint32
	Family    uint16
	DPort     uint16
	DataLen   uint16
	IPVersion uint8
	_         [1]byte
	Comm      [16]byte
	Dst       [16]byte
	Data      [512]byte
}

func runTrace() error {
	if err := rlimit.RemoveMemlock(); err != nil {
		return fmt.Errorf("remove memlock limit: %w", err)
	}

	objs := egressObjects{}
	if err := loadEgressObjects(&objs, nil); err != nil {
		return fmt.Errorf("load bpf objects: %w", err)
	}
	defer objs.Close()

	links := []link.Link{}
	defer func() {
		for _, l := range links {
			_ = l.Close()
		}
	}()

	enterConnect, err := link.Tracepoint("syscalls", "sys_enter_connect", objs.HandleEnter, nil)
	if err != nil {
		return fmt.Errorf("attach sys_enter_connect: %w", err)
	}
	links = append(links, enterConnect)

	sendto, err := link.Tracepoint("syscalls", "sys_enter_sendto", objs.HandleSendto, nil)
	if err != nil {
		return fmt.Errorf("attach sys_enter_sendto: %w", err)
	}
	links = append(links, sendto)

	recvfromEnter, err := link.Tracepoint("syscalls", "sys_enter_recvfrom", objs.HandleRecvfromEnter, nil)
	if err != nil {
		return fmt.Errorf("attach sys_enter_recvfrom: %w", err)
	}
	links = append(links, recvfromEnter)

	recvfromExit, err := link.Tracepoint("syscalls", "sys_exit_recvfrom", objs.HandleRecvfromExit, nil)
	if err != nil {
		return fmt.Errorf("attach sys_exit_recvfrom: %w", err)
	}
	links = append(links, recvfromExit)

	reader, err := ringbuf.NewReader(objs.Events)
	if err != nil {
		return fmt.Errorf("open ringbuf: %w", err)
	}
	defer reader.Close()

	dnsNames := map[string]string{}

	fmt.Println("listening for connect() and DNS events; run `curl https://example.com` in another terminal")
	fmt.Println("TYPE PID COMM DETAIL")

	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				return nil
			}
			return fmt.Errorf("read ringbuf: %w", err)
		}

		var e event
		size := int(unsafe.Sizeof(e))
		if len(record.RawSample) < size {
			continue
		}
		copy(unsafe.Slice((*byte)(unsafe.Pointer(&e)), size), record.RawSample[:size])

		comm := strings.TrimRight(string(e.Comm[:]), "\x00")
		dataLen := int(e.DataLen)
		if dataLen > len(e.Data) {
			dataLen = len(e.Data)
		}
		payload := e.Data[:dataLen]

		switch e.Type {
		case eventConnect:
			destIP := formatIP(e.IPVersion, e.Dst[:])
			dest := formatDestination(e.IPVersion, e.Dst[:], e.DPort)
			if dnsName, ok := dnsNames[destIP]; ok {
				dest = net.JoinHostPort(dnsName, fmt.Sprintf("%d", e.DPort))
			}
			fmt.Printf("CONNECT %d %s %s\n", e.PID, comm, dest)
		case eventDNSQuery:
			if q, ok := parseDNSQuery(payload); ok {
				fmt.Printf("DNS_QUERY %d %s %s\n", e.PID, comm, q)
			}
		case eventDNSResponse:
			q, ips, ok := parseDNSResponse(payload)
			if !ok {
				continue
			}
			fmt.Printf("DNS_QUERY %d %s %s\n", e.PID, comm, q)
			for _, ip := range ips {
				dnsNames[ip] = q
				fmt.Printf("DNS_ANSWER %d %s %s -> %s\n", e.PID, comm, q, ip)
			}
		}
	}
}

func parseDNSQuery(msg []byte) (string, bool) {
	if len(msg) < 12 || binary.BigEndian.Uint16(msg[4:6]) == 0 || binary.BigEndian.Uint16(msg[2:4])&0x8000 != 0 {
		return "", false
	}
	name, _, ok := dnsName(msg, 12)
	return name, ok && name != ""
}

func parseDNSResponse(msg []byte) (string, []string, bool) {
	if len(msg) < 12 || binary.BigEndian.Uint16(msg[2:4])&0x8000 == 0 {
		return "", nil, false
	}
	qd := int(binary.BigEndian.Uint16(msg[4:6]))
	an := int(binary.BigEndian.Uint16(msg[6:8]))
	if qd == 0 || an == 0 {
		return "", nil, false
	}
	qname, off, ok := dnsName(msg, 12)
	if !ok || off+4 > len(msg) {
		return "", nil, false
	}
	off += 4
	ips := []string{}
	for i := 0; i < an && off < len(msg); i++ {
		_, next, ok := dnsName(msg, off)
		if !ok || next+10 > len(msg) {
			break
		}
		typ := binary.BigEndian.Uint16(msg[next : next+2])
		rdlen := int(binary.BigEndian.Uint16(msg[next+8 : next+10]))
		rdata := next + 10
		if rdata+rdlen > len(msg) {
			break
		}
		switch typ {
		case 1:
			if rdlen == 4 {
				ips = append(ips, net.IP(msg[rdata:rdata+4]).String())
			}
		case 28:
			if rdlen == 16 {
				ips = append(ips, net.IP(msg[rdata:rdata+16]).String())
			}
		}
		off = rdata + rdlen
	}
	return qname, ips, len(ips) > 0
}

func dnsName(msg []byte, off int) (string, int, bool) {
	labels := []string{}
	start := off
	jumped := false
	for depth := 0; depth < 20; depth++ {
		if off >= len(msg) {
			return "", 0, false
		}
		l := int(msg[off])
		off++
		if l == 0 {
			if jumped {
				return strings.Join(labels, "."), start + 2, true
			}
			return strings.Join(labels, "."), off, true
		}
		if l&0xc0 == 0xc0 {
			if off >= len(msg) {
				return "", 0, false
			}
			ptr := ((l & 0x3f) << 8) | int(msg[off])
			if !jumped {
				start = off - 1
			}
			off = ptr
			jumped = true
			continue
		}
		if l&0xc0 != 0 || off+l > len(msg) {
			return "", 0, false
		}
		labels = append(labels, string(msg[off:off+l]))
		off += l
	}
	return "", 0, false
}

func formatDestination(version uint8, dst []byte, port uint16) string {
	ip := formatIP(version, dst)
	switch version {
	case 4:
		return ip + fmt.Sprintf(":%d", port)
	case 6:
		return net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	default:
		return fmt.Sprintf("?:%d", port)
	}
}

func formatIP(version uint8, dst []byte) string {
	switch version {
	case 4:
		return net.IPv4(dst[0], dst[1], dst[2], dst[3]).String()
	case 6:
		return net.IP(dst[:16]).String()
	default:
		return ""
	}
}
