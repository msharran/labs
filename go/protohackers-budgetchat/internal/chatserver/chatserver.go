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
	ErrNameNotProvided = fmt.Errorf("NameNotProvided")
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
	mem, ok := c.room[conn]
	c.mu.Unlock()

	switch {
	case !ok: // if member not in room , ask name
		_, err := conn.Write([]byte("* What is your name?"))
		if err != nil {
			c.log.Error("failed to ask member name", err)
			return false
		}

		// add the member to room without joining them (joined=false)
		c.mu.Lock()
		c.room[conn] = &member{}
		c.mu.Unlock()

		return true

	case !mem.joined: // validate and add new member to the room
		// get name from member
		nameBuf := [16]byte{}
		nameN, err := conn.Read(nameBuf[:])
		if err != nil {
			c.log.Error("failed to get name from member", err)
			return false
		}

		name := string(nameBuf[:nameN])
		name = strings.TrimSuffix(name, "\n")
		log := c.log.With("member", name)

		log.Info("adding new member to the room")

		// kick out member if name is not provided
		if name == "" {
			c.mu.Lock()
			delete(c.room, conn)
			c.mu.Unlock()
			log.Error("name not provided, kicked out member", ErrNameNotProvided)

			_, err := conn.Write([]byte("name is not provided, please enter a name when you join!\n"))
			if err != nil {
				log.Error("failed to write back to connection", nil)
				return false
			}
			return false
		}

		// now add the member to the room
		mem.name = name
		mem.joined = true

		// tell others new member has joined
		var others []string
		c.mu.Lock()
		for mc, m := range c.room {
			if mc != conn {
				others = append(others, m.name)

				msg := fmt.Sprintf("* %s joined the room\n", name)
				_, err = mc.Write([]byte(msg))
				if err != nil {
					log.Error("failed to write back to connection", nil)
				}
			}
		}
		c.mu.Unlock()

		// tell the member who are present in the room
		msg := fmt.Sprintf("room contains: %s\n", strings.Join(others, ","))
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Error("failed to write back to connection", nil)
		}

		return true

	default:
		name := mem.name
		log := c.log.With("member", name)

		// read message from member
		buf := [1000]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Error("error reading message from connection", nil)
			return false
		}

		// send to other members in the room
		c.mu.Lock()
		for mConn, m := range c.room {
			if mConn != conn {
				log.Info("sending message", "to", m.name)
				defer log.Info("sent message", "to", m.name)

				_, err := mConn.Write([]byte(fmt.Sprintf("[%s] %s", name, string(buf[:n]))))
				if err != nil {
					c.log.Error("error sending message", nil, "from", name, "to", m.name)
				}
			}
		}
		c.mu.Unlock()

		return true
	}
}

func (c *chatServer) exitRoom(conn net.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	mem, ok := c.room[conn]
	name := "unknown"
	if ok {
		name = mem.name
	}

	err := conn.Close()
	if err != nil {
		c.log.Error("error while closing the ", err, "member", name)
		return
	}

	c.log.Info("member left ", "member", name)
	delete(c.room, conn)

	// inform all members when a member leaves
	if ok {
		msg := fmt.Sprintf("* %s left the room\n", name)
		for mConn, m := range c.room {
			_, err := mConn.Write([]byte(msg))
			if err != nil { // only print error in server if writing fails. it can be soft failure
				c.log.Error("error while writing 'member left' message", err, "failed_member", m.name)
			}
		}
	}
}

func (c *chatServer) sendMessage(conn net.Conn, message string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	name := c.room[conn].name

	for mConn, m := range c.room {
		if mConn != conn {
			_, err := mConn.Write([]byte(fmt.Sprintf("[%s] %s", strings.TrimSuffix(name, "\n"), message)))
			if err != nil {
				c.log.Error("error sending message", nil, "from", name, "to", m.name)
			}
		}
	}
}
