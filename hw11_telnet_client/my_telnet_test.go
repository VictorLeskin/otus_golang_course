package main

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockReadCloser имитирует io.ReadCloser

type MockReadCloser1 struct {
	closed bool
	Data   []string
	pos    int
}

func (m *MockReadCloser1) Read(p []byte) (n int, err error) {
	if m.closed {
		return 0, fmt.Errorf("stream had closed")
	}
	if m.pos < len(m.Data) {
		p = []byte(m.Data[m.pos])
		m.pos++
		return len(p), nil
	}

	return 0, io.EOF
}

func (m *MockReadCloser1) Close() error {
	m.closed = true
	return nil
}

// MockWriter имитирует io.Writer и сохраняет записанные данные
type MockWriter1 struct {
	buffer string
}

func (m *MockWriter1) Write(p []byte) (n int, err error) {
	m.buffer += string(p)
	return len(p), nil
}

func (m *MockWriter1) Free() (ret string) {
	ret = m.buffer
	m.buffer = ""
	return ret
}

type MockConn struct {
	net.Conn
}

func (c MockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (c MockConn) Write(b []byte) (n int, err error) {
	return 0, nil
}

type MyDialer struct {
	mockConn    MockConn
	err         error
	realTimeOut time.Duration

	network, address string
	timeout          time.Duration
}

var myDialer MyDialer

func MyDialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	myDialer.network = network
	myDialer.address = address
	myDialer.timeout = timeout

	if myDialer.err != nil {
		return nil, myDialer.err
	}

	if myDialer.realTimeOut > timeout {
		time.Sleep(myDialer.realTimeOut)
		return nil, fmt.Errorf("timeout error")
	}

	return myDialer.mockConn, myDialer.err
}

func Test_MockReadCloser1_Ctor(t *testing.T) {
	var t0 MockReadCloser1
	assert.Equal(t, false, t0.closed)
	assert.Equal(t, 0, len(t0.Data))
	assert.Equal(t, 0, t0.pos)
}

func Test_NewTelnetClient(t *testing.T) {
	in := &MockReadCloser1{}
	out := &MockWriter1{}
	t0 := NewTelnetClient("1.2.3.4:5", 11*time.Second, in, out)

	c, ok := t0.(*MyTelnetClient)

	assert.True(t, ok)
	assert.Equal(t, "1.2.3.4:5", c.address)
	assert.Equal(t, 11*time.Second, c.timeout)

	c1, ok1 := c.in.(*MockReadCloser1)
	assert.True(t, ok1)
	assert.Equal(t, in, c1)

	c2, ok2 := c.out.(*MockWriter1)
	assert.True(t, ok2)
	assert.Equal(t, out, c2)

	assert.NotNil(t, c.ctx)
	assert.NotNil(t, c.cancel)
}

func Test_MyTelnetClient_Connect(t *testing.T) {

	t.Run("connection is ok", func(t *testing.T) {
		t0 := MyTelnetClient{
			address: "1.2.3.4:5",
			timeout: 1 * time.Second}

		myDialer = MyDialer{}
		t0.dialer = MyDialTimeout

		assert.Nil(t, t0.conn)

		err := t0.Connect()

		assert.Nil(t, err)
		assert.Equal(t, "tcp", myDialer.network)
		assert.Equal(t, "1.2.3.4:5", myDialer.address)
		assert.Equal(t, 1*time.Second, myDialer.timeout)
		assert.NotNil(t, t0.conn)
	})

	t.Run("bad connection by unknown error", func(t *testing.T) {
		t0 := MyTelnetClient{
			address: "1.2.3.4:5",
			timeout: 1 * time.Second}

		myDialer = MyDialer{}
		t0.dialer = MyDialTimeout
		// t0.dialer parameters.
		myDialer.err = fmt.Errorf("unknown error")

		assert.Nil(t, t0.conn)

		err := t0.Connect()

		assert.Equal(t, "unknown error", err.Error())
		assert.Equal(t, "tcp", myDialer.network)
		assert.Equal(t, "1.2.3.4:5", myDialer.address)
		assert.Equal(t, 1*time.Second, myDialer.timeout)
		assert.Nil(t, t0.conn)
	})

	t.Run("dissconnected by timeout: asked 1 sec, waited more", func(t *testing.T) {
		t0 := MyTelnetClient{
			address: "1.2.3.4:5",
			timeout: 1 * time.Second} // 1 sec timeout

		myDialer = MyDialer{}
		t0.dialer = MyDialTimeout
		// t0.dialer parameters.
		myDialer.realTimeOut = 2 * time.Second // real timeout

		assert.Nil(t, t0.conn)

		err := t0.Connect()

		assert.Equal(t, "timeout error", err.Error())
		assert.Equal(t, "tcp", myDialer.network)
		assert.Equal(t, "1.2.3.4:5", myDialer.address)
		assert.Equal(t, 1*time.Second, myDialer.timeout)
		assert.Nil(t, t0.conn)
	})

}
