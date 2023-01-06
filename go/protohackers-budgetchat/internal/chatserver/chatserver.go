package chatserver

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

// errors
var (
	ErrNameNotProvided = fmt.Errorf("name is not provided; members will be joined only if name is provided")
)

type ServerArgs struct {
	Log     *slog.Logger
	Timeout time.Duration
}

type chatServer struct {
	log     *slog.Logger
	timeout time.Duration
	room    map[net.Conn]*member
	mu      sync.Mutex
}

type member struct {
	name   string
	joined bool
}

type Listener interface {
	Listen(addr string) error
}

func NewServer(reqs ServerArgs) Listener {
	return &chatServer{log: reqs.Log, timeout: reqs.Timeout, room: make(map[net.Conn]*member, 0)}
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
		go b.enterRoom(conn)
	}
}

func (c *chatServer) enterRoom(conn net.Conn) {
	defer c.exitRoom(conn)

	conn.SetDeadline(time.Now().Add(c.timeout))

	for c.reconcileMember(conn) {
		// reset the current deadline and add new one
		// for each reconciliation
		conn.SetDeadline(time.Time{})
		conn.SetDeadline(time.Now().Add(c.timeout))
	}
}

func (c *chatServer) reconcileMember(conn net.Conn) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	mem := c.room[conn]

	// if member has not joined , ask name
	if mem == nil {
		_, err := conn.Write([]byte("What is your name?"))
		if err != nil {
			c.log.Error("failed to ask member name", err)
			return false
		}

		// add the member to  without joining them (joined=false)
		c.room[conn] = &member{}

		return true
	}

	nameBuf := [16]byte{}
	nameN, err := conn.Read(nameBuf[:])
	if err != nil {
		c.log.Error("failed to get name from member", err)
		return false
	}

	name := string(nameBuf[:nameN])
	log := c.log.With("member", name)

	// kick out member if name is not provided
	if name == "\n" {
		log.Error("", ErrNameNotProvided)

		_, err := conn.Write([]byte(ErrNameNotProvided.Error()))
		if err != nil {
			log.Error("failed to write back to connection", nil)
			return false
		}
		return false
	}

	// add new member to the room
	if !mem.joined {
		mem.name = name

		// get other joined members
		var mm []string
		for _, m := range c.room {
			if m.name != name {
				mm = append(mm, m.name)
			}
		}

		// show members list to new member
		msg := fmt.Sprintf("room contains: %s", strings.Join(mm, ","))
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Error("failed to write back to connection", nil)
			return false
		}
		mem.joined = true
		return true
	}

	return true
}

func (c *chatServer) exitRoom(conn net.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	mem := c.room[conn]

	err := conn.Close()
	if err != nil {
		c.log.Error("error while closing the ", err, "member", mem.name)
		return
	}

	c.log.Info("member left ", "member", mem.name)
	delete(c.room, conn)

	// inform all members
	msg := fmt.Sprintf("%s left the ", mem.name)
	for mConn, m := range c.room {
		_, err := mConn.Write([]byte(msg))
		if err != nil { // only print error in server if writing fails. it can be soft failure
			c.log.Error("error while writing 'member left' message", err, "failed_member", m.name)
		}
	}
}
