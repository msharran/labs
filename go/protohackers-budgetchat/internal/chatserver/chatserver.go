package chatserver

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/exp/slog"
)

type ServerArgs struct {
	Log     *slog.Logger
	Timeout time.Duration
}

type chatServer struct {
	log     *slog.Logger
	timeout time.Duration
	channel map[net.Conn]member
	mu      sync.Mutex
}

type member struct {
	joined bool
	name   string
}

type Listener interface {
	Listen(addr string) error
}

func NewServer(reqs ServerArgs) Listener {
	return &chatServer{log: reqs.Log, timeout: reqs.Timeout, channel: make(map[net.Conn]member, 0)}
}

func (b *chatServer) Listen(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	b.log.Info("started budget-chat server", "addr", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go b.handleChat(conn)
	}
}

func (c *chatServer) handleChat(conn net.Conn) {
	id := fmt.Sprintf("m_%s", ulid.Make())
	log := slog.With("member", id)

	defer func() {
		conn.Close()
		log.Info("member exited")
	}()

	conn.SetDeadline(time.Now().Add(c.timeout))

	for {
		mem, ok := c.channel[conn]

		// if member not present in huddle, ask name
		if !ok {
			_, err := conn.Write([]byte("What is your name?"))
			if err != nil {
				log.Error("failed to get member name", err)
				return
			}

			// add the member as joined=false
			c.mu.Lock()
			defer c.mu.Unlock()
			c.channel[conn] = member{joined: false}

			continue
		}

		nameBuf := [16]byte{}
		nameN, err := conn.Read(nameBuf[:])
		if err != nil {
			log.Error("failed to get name from member", err)
			return
		}
		name := nameBuf[:nameN]
	}
}
