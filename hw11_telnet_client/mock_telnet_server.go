package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

// MockTelnetServer создает тестовый telnet сервер
type MockTelnetServer struct {
	Port     int
	Response string
	Messages []string
	listener net.Listener
}

// Start запускает mock сервер
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
				return // Сервер остановлен
			}
			go m.handleConnection(conn)
		}
	}()

	return nil
}

func (m *MockTelnetServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Отправляем приветствие
	if m.Response != "" {
		conn.Write([]byte(m.Response))
	}

	// Читаем сообщения от клиента
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Server read error: %v\n", err)
			}
			return
		}

		msg := strings.TrimSpace(string(buf[:n]))
		m.Messages = append(m.Messages, msg)

		// Отправляем эхо-ответ
		if m.Response != "" {
			conn.Write([]byte(m.Response))
		}
	}
}

// Stop останавливает сервер
func (m *MockTelnetServer) Stop() {
	if m.listener != nil {
		m.listener.Close()
	}
}

// GetLastMessage возвращает последнее полученное сообщение
func (m *MockTelnetServer) GetLastMessage() string {
	if len(m.Messages) == 0 {
		return ""
	}
	return m.Messages[len(m.Messages)-1]
}
