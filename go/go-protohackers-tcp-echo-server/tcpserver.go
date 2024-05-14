package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TCPServerConfig struct {
	Host string
	Port string
}

func NewTCPServer(cfg TCPServerConfig) *TCPServer {
	return &TCPServer{host: cfg.Host, port: cfg.Port, prices: map[net.Conn]map[uint32]uint32{}}
}

type TCPServer struct {
	host      string
	port      string
	prices    map[net.Conn]map[uint32]uint32
	pricesMux sync.Mutex
	net.Listener
}

func (svr *TCPServer) Addr() string {
	return fmt.Sprintf("%s:%s", svr.host, svr.port)
}

func (svr *TCPServer) ListenAndServe() {
	// listen for incoming connections
	l, err := net.Listen("tcp", svr.Addr())
	if err != nil {
		logrus.Fatal(err)
	}
	svr.Listener = l

	logrus.Info("accepting incoming connections")
	for {
		conn, err := l.Accept()
		if err != nil {
			logrus.Fatal("error accepting connection ", err.Error())
		}
		go svr.handleRequest(conn)
	}
}

func (svr *TCPServer) handleRequest(conn net.Conn) {
	log := logrus.WithField("request-id", uuid.NewString())

	log.Infof("handling connection")

	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("unable to close the connection")
		}
	}()

	r := bufio.NewReader(conn)
	for {
		log.Info("processing connection")
		msg := make([]byte, 9)
		_, err := r.Read(msg)
		if err != nil {
			// if EOF return
			if err.Error() == "EOF" {
				log.Info("connection closed by client")
				break
			}
			log.Error("error reading from connection", err)
			break
		}

		switch string(msg[0]) {
		case "I":
			svr.insertPrice(conn, msg, log)
		case "Q":
			meanPrice := svr.queryPrice(conn, msg, log)
			log.Infof("mean price: %v", meanPrice)

			data := make([]byte, 4)
			binary.BigEndian.PutUint32(data, meanPrice)

			_, err := conn.Write(data)
			if err != nil {
				log.Error("error writing to connection", err)
			}
		default:
			log.Errorf("undefined behavior: %s, try sending again", string(msg[0]))
		}
	}
}

func (svr *TCPServer) insertPrice(conn net.Conn, msg []byte, log *logrus.Entry) {
	svr.pricesMux.Lock()
	defer svr.pricesMux.Unlock()

	if svr.prices[conn] == nil {
		svr.prices[conn] = make(map[uint32]uint32)
	}

	ts := binary.BigEndian.Uint32(msg[1:5])
	price := binary.BigEndian.Uint32(msg[5:])
	svr.prices[conn][ts] = price

	log.Infof("inserted price. total: %d", len(svr.prices[conn]))
}

func (svr *TCPServer) queryPrice(conn net.Conn, msg []byte, log *logrus.Entry) uint32 {
	svr.pricesMux.Lock()
	prices := svr.prices[conn]
	svr.pricesMux.Unlock()

	log.Infof("queried prices: %v", len(prices))

	if prices == nil {
		log.Info("no prices added for the connection, returning 0")
		return 0
	}

	min := binary.BigEndian.Uint32(msg[1:5])
	max := binary.BigEndian.Uint32(msg[5:])

	var pp []uint32
	for ts, price := range prices {
		if ts >= min && ts <= max {
			pp = append(pp, price)
		}
	}

	return mean(pp)
}

func mean(prices []uint32) uint32 {
	var m uint32
	for _, price := range prices {
		m += price
	}
	m = m / uint32(len(prices))
	return m
}
