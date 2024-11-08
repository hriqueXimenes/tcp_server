package server

import (
	"fmt"
	"net"
)

type mockListener struct {
	shouldReturnErrorOnAccept bool
	shouldReturnErrorOnWrite  bool

	onAcceptCount int
	onWriteCount  int
}

func (m *mockListener) Accept() (net.Conn, error) {
	m.onAcceptCount++

	if m.shouldReturnErrorOnAccept {
		return nil, fmt.Errorf("mock error")
	}

	return nil, nil
}

func (m *mockListener) Write(conn net.Conn, req []byte) error {
	m.onWriteCount++

	if m.shouldReturnErrorOnWrite {
		return fmt.Errorf("mock error")
	}

	return nil
}
