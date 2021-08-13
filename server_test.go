package socks

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

func getConnection(t *testing.T) net.Conn {
	srv, err := NewSocks4Server()
	if err != nil {
		t.Fatalf("Could not start SOCKS server: %v", err)
	}

	go func() {
		srv.ListenAndServe("localhost:1080")
	}()

	conn, err := net.Dial("tcp", "localhost:1080")
	if err != nil {
		t.Fatalf("Could not connect to SOCKS server: %v", err)
	}

	return conn
}

func TestNewSocks4Server_InvalidProtocolVersion(t *testing.T) {
	conn := getConnection(t)

	// handshake
	conn.Write([]byte{5})

	// read reply
	expected := []byte{
		0,
		91,
	}
	out := make([]byte, len(expected))

	conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := io.ReadAtLeast(conn, out, len(out)); err != nil {
		t.Fatalf("Failed to read %v bytes: %v", expected, err)
	}

	if !bytes.Equal(out, expected) {
		t.Fatalf("Expected %v, got: %v", expected, out)
	}
}

func TestSocks4_Connect(t *testing.T) {
	conn := getConnection(t)

	// handshake
	conn.Write([]byte{4})

	// read reply
	expected := []byte{
		0,
		90,
	}
	out := make([]byte, len(expected))

	conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := io.ReadAtLeast(conn, out, len(out)); err != nil {
		t.Fatalf("Failed to read %v bytes: %v", expected, err)
	}

	if !bytes.Equal(out, expected) {
		t.Fatalf("Expected %v, got: %v", expected, out)
	}
}
