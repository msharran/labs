package main

import (
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

type event struct {
	PID       uint32
	TGID      uint32
	Family    uint16
	DPort     uint16
	IPVersion uint8
	_         [1]byte
	Comm      [16]byte
	Dst       [16]byte
	_         [2]byte
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

	enterLink, err := link.Tracepoint("syscalls", "sys_enter_connect", objs.HandleEnter, nil)
	if err != nil {
		return fmt.Errorf("attach sys_enter_connect: %w", err)
	}
	defer enterLink.Close()

	reader, err := ringbuf.NewReader(objs.Events)
	if err != nil {
		return fmt.Errorf("open ringbuf: %w", err)
	}
	defer reader.Close()

	fmt.Println("listening for connect() calls; run `curl https://example.com` in another terminal")
	fmt.Println("PID COMM DEST")

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
		dest := formatDestination(e.IPVersion, e.Dst[:], e.DPort)
		fmt.Printf("%d %s %s\n", e.PID, comm, dest)
	}
}

func formatDestination(version uint8, dst []byte, port uint16) string {
	switch version {
	case 4:
		return net.IPv4(dst[0], dst[1], dst[2], dst[3]).String() + fmt.Sprintf(":%d", port)
	case 6:
		return net.JoinHostPort(net.IP(dst[:16]).String(), fmt.Sprintf("%d", port))
	default:
		return fmt.Sprintf("?:%d", port)
	}
}
