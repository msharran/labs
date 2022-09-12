package main

import (
	"github.com/stretchr/testify/require"
	"net"
	"sync"
	"testing"
	"time"
)

const (
	TestNetHost = "localhost"
	TestNetPort = "8081"
	TestNetAddr = TestNetHost + ":" + TestNetPort
)

type fixture struct {
	T   *testing.T
	svr *TCPServer
}

func newFixture(t *testing.T) *fixture {
	svr := NewTCPServer(TCPServerConfig{Host: TestNetHost, Port: TestNetPort})
	return &fixture{T: t, svr: svr}
}

func (f *fixture) Stop() {
	err := f.svr.Close()
	if err != nil {
		f.T.Fatal(err)
	}
}

func (f *fixture) RunTestServerInBackground() {
	// run TCP server in another goroutine
	go func() { f.svr.ListenAndServe() }()

	// wait for server to start
	<-time.After(1 * time.Second)
}

func (f *fixture) GetInsertPrices() [][]byte {
	return [][]byte{
		{0x49, 0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x65}, // I 12345 101
		{0x49, 0x00, 0x00, 0x30, 0x3a, 0x00, 0x00, 0x00, 0x66}, // I 12346 102
		{0x49, 0x00, 0x00, 0x30, 0x3b, 0x00, 0x00, 0x00, 0x64}, // I 12347 100
		{0x49, 0x00, 0x00, 0xa0, 0x00, 0x00, 0x00, 0x00, 0x05}, // I 40960 5
	}
}

func (f *fixture) GetQueryPrice() []byte {
	return []byte{0x51, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x40, 0x00} // Q 12288 16384
}

func TestInsertAndQuery(t *testing.T) {
	f := newFixture(t)
	f.RunTestServerInBackground()

	conn, err := net.Dial("tcp", TestNetAddr)
	require.NoError(t, err)
	defer func() { require.NoError(t, conn.Close()) }()

	for _, ins := range f.GetInsertPrices() {
		_, err = conn.Write(ins)
		require.NoError(t, err)
	}

	_, err = conn.Write(f.GetQueryPrice())
	require.NoError(t, err)

	readBuf := make([]byte, 4)
	_, err = conn.Read(readBuf)
	require.NoError(t, err)

	// should return 101
	require.Equal(t, string([]byte{0x00, 0x00, 0x00, 0x65}), string(readBuf))
}

func TestUnknownBehaviour(t *testing.T) {
	f := newFixture(t)
	f.RunTestServerInBackground()

	conn, err := net.Dial("tcp", TestNetAddr)
	require.NoError(t, err)
	defer func() { require.NoError(t, conn.Close()) }()

	_, err = conn.Write([]byte("Iaaaabbbb"))
	require.NoError(t, err)

	_, err = conn.Write([]byte("Daaaacccc"))
	require.NoError(t, err)

	_, err = conn.Write([]byte("Qaaaacccc"))
	require.NoError(t, err)

	readBuf := make([]byte, 8)
	_, err = conn.Read(readBuf)
	require.NoError(t, err)
	require.Equal(t, "aaaabbbb", string(readBuf))
}

func Test5ConcurrentConnections(t *testing.T) {
	f := newFixture(t)
	f.RunTestServerInBackground()

	var conns []net.Conn
	defer func() {
		for _, conn := range conns {
			require.NoError(t, conn.Close())
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			conn, err := net.Dial("tcp", TestNetAddr)
			require.NoError(t, err)
			conns = append(conns, conn)

			defer func() { require.NoError(t, conn.Close()) }()

			for _, ins := range f.GetInsertPrices() {
				_, err = conn.Write(ins)
				require.NoError(t, err)
			}

			_, err = conn.Write(f.GetQueryPrice())
			require.NoError(t, err)

			readBuf := make([]byte, 4)
			_, err = conn.Read(readBuf)
			require.NoError(t, err)

			// should return 101
			require.Equal(t, string([]byte{0x00, 0x00, 0x00, 0x65}), string(readBuf))
		}(i)
	}
	wg.Wait()
}
