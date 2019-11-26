package httputil

import (
	"context"
	"net"
)

func unixSockSink(sinkAddr string) func(
	ctx context.Context, net, addr string,
) (net.Conn, error) {
	d := new(net.Dialer)
	return func(ctx context.Context, net, addr string) (net.Conn, error) {
		return d.DialContext(ctx, "unix", sinkAddr)
	}
}
