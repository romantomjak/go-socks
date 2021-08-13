package socks

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

const (
	RequestGranted                   = 0x5A // 90
	RequestRejectedOrFailed          = 0x5B // 91
	RequestFailedIdentdNotRunning    = 0x5C // 92
	RequestFailedIdentdInvalidUserID = 0x5D // 93
)

type Client struct {
	addr string
}

// NewV4Client creates a new SOCKS v4 client
func NewV4Client(addr string) (*Client, error) {
	client := &Client{
		addr: addr,
	}
	return client, nil
}

// Connect instructs the SOCKS server to establish a connection
// to the server at addr. It is callers responsibility to close
// the connection.
func (c *Client) Connect(addr, username string) (net.Conn, error) {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return nil, err
	}

	ip, port, err := parseIPv4(addr)
	if err != nil {
		return nil, err
	}

	bufLen := 8 + len(username) + 1
	buf := make([]byte, bufLen)

	buf[0] = 4
	buf[1] = 1
	binary.BigEndian.PutUint16(buf[2:], port)
	binary.BigEndian.PutUint32(buf[4:], ip)
	copy(buf[8:], []byte(username))
	buf[8+len(username)] = 0

	n, err := conn.Write(buf)
	if err != nil {
		return nil, fmt.Errorf("write failed: %v", err)
	}

	if n != bufLen {
		return nil, fmt.Errorf("wrote %d of %d bytes", n, bufLen)
	}

	respBufLen := 8
	respBuf := make([]byte, respBufLen)
	n, err = conn.Read(respBuf)
	if err != nil {
		return nil, err
	}

	if n != respBufLen {
		return nil, fmt.Errorf("read %d of %d bytes", n, respBufLen)
	}

	switch respBuf[1] {
	case RequestGranted:
		return conn, nil
	case RequestRejectedOrFailed:
		return nil, fmt.Errorf("rejected or failed")
	case RequestFailedIdentdNotRunning:
		return nil, fmt.Errorf("identd not running (or not reachable from server)")
	case RequestFailedIdentdInvalidUserID:
		return nil, fmt.Errorf("identd could not confirm the user ID")
	}

	return nil, fmt.Errorf("unknown reply code %d", respBuf[1])
}

func parseIPv4(addr string) (uint32, uint16, error) {
	ipstr, portstr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, 0, err
	}

	ipaddr := net.ParseIP(ipstr)
	if ipaddr == nil {
		return 0, 0, fmt.Errorf("invalid ip address %s", ipstr)
	}

	ipbytes := ipaddr.To4()
	if ipbytes == nil {
		return 0, 0, fmt.Errorf("invalid ipv4 address %s", ipaddr)
	}

	port, err := strconv.Atoi(portstr)
	if err != nil {
		return 0, 0, err
	}

	return binary.BigEndian.Uint32(ipbytes), uint16(port), nil
}
