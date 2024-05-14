package main

import (
	"net"
	"testing"

	"golang.org/x/net/nettest"
)

func TestServer(t *testing.T) {
	ln := newLocalListener(t)
	s := server{}
	s.Serve()

	// Start a connection between two endpoints.
	var c1, c2 net.Conn
	defer func() {
		c1.Close()
		c2.Close()
	}()
	var err1, err2 error
	done := make(chan bool)
	go func() {
		c2, err2 = ln.Accept()
		close(done)
	}()
	c1, err1 = net.Dial(ln.Addr().Network(), ln.Addr().String())
	<-done

	if err1 != nil {
		t.Error(err1)
	}
	if err2 != nil {
		t.Error(err2)
	}
}

func newLocalListener(t testing.TB) net.Listener {
	t.Helper()

	ln, err := nettest.NewLocalListener("tcp")
	if err != nil {
		t.Fatal(err)
	}
	return ln
}
