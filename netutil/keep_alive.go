package netutil

import (
	"net"
	"time"
)

type keepAliveListner struct {
	*net.TCPListener
}

func (ln keepAliveListner) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}

	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// WrapKeepAlive wraps the listener. If the listener is a TCP listener, it
// sets keep alive to 3 minute.
func WrapKeepAlive(ln net.Listener) net.Listener {
	tcpLis, ok := ln.(*net.TCPListener)
	if !ok {
		return ln
	}
	return keepAliveListner{tcpLis}
}
