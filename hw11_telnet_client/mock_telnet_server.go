package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

// MockTelnetServer created test telnet server.
type MockTelnetServer struct {
	Port     int
	Response string
	Messages []string
	listener net.Listener
}

// Start mock server.
func (m *MockTelnetServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", m.Port))
	if err != nil {
		return err
	}
	m.listener = listener
	m.Messages = []string{}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return // server stop
			}
			go m.handleConnection(conn)
		}
	}()

	return nil
}

func (m *MockTelnetServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Send hello.
	if m.Response != "" {
		conn.Write([]byte(m.Response))
	}

	// Read clients's messsages.
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Printf("Server read error: %v\n", err)
			}
			return
		}

		msg := strings.TrimSpace(string(buf[:n]))
		m.Messages = append(m.Messages, msg)

		// Send echo ansver
		if m.Response != "" {
			conn.Write([]byte(m.Response))
		}
	}
}

// Stop server.
func (m *MockTelnetServer) Stop() {
	if m.listener != nil {
		m.listener.Close()
	}
}

// GetLastMessage return last received message.
func (m *MockTelnetServer) GetLastMessage() string {
	if len(m.Messages) == 0 {
		return ""
	}
	return m.Messages[len(m.Messages)-1]
}
