package server

import (
	"fmt"
	"net"
)

type Listener interface {
	Accept() (net.Conn, error)
}

type listenerImpl struct {
	listener net.Listener
}

func newListener(port int, addr, protocol string) (Listener, error) {
	newListener, err := net.Listen(protocol, fmt.Sprintf(`%s:%v`, addr, port))
	if err != nil {
		return nil, err
	}

	return &listenerImpl{
		listener: newListener,
	}, nil
}

func (l *listenerImpl) Accept() (net.Conn, error) {
	return l.listener.Accept()
}
