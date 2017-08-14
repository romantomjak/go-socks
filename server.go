package go_socks

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

type Server struct {
	logger *log.Logger
}

type Connection struct {
	server *Server
	socket net.Conn
	r      io.Reader
}

func NewSocks4Server() (*Server, error) {
	server := &Server{
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	return server, nil
}

func (s *Server) ListenAndServe(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.serve(l)
}

func (s *Server) serve(l net.Listener) error {
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		c := s.newConn(conn)
		go c.serve()
	}
}

func (s *Server) newConn(conn net.Conn) *Connection {
	c := &Connection{
		server: s,
		socket: conn,
		r:      bufio.NewReader(conn),
	}
	return c
}

func (c *Connection) serve() {
	defer c.socket.Close()

	// check version
	version := []byte{0}
	_, err := c.r.Read(version)
	if err != nil {
		c.server.logger.Printf("Failed to read SOCKS version: %v", err)

		msg := make([]byte, 2)
		msg[0] = 0
		msg[1] = 91

		_, err := c.socket.Write(msg)
		if err != nil {
			c.server.logger.Printf("Failed to write to client socket: %v", err)
		}

		return
	}

	// check compatibility
	if version[0] != uint8(4) {
		c.server.logger.Printf("Unsupported SOCKS version: %v", version[0])

		msg := make([]byte, 2)
		msg[0] = 0
		msg[1] = 91

		_, err := c.socket.Write(msg)
		if err != nil {
			c.server.logger.Printf("Failed to write to client socket: %v", err)
		}

		return
	}

	// accept
	msg := make([]byte, 2)
	msg[0] = 0
	msg[1] = 90

	_, err = c.socket.Write(msg)
	if err != nil {
		c.server.logger.Printf("Failed to write to client socket: %v", err)
	}

	return
}
