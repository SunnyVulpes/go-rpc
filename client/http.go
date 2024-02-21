package client

import (
	"bufio"
	"fmt"
	"gorpc/service"
	"io"
	"net"
	"net/http"
	"strings"
)

func NewHTTPClient(conn net.Conn, opt *service.Option) (*Client, error) {
	_, _ = io.WriteString(conn, fmt.Sprintf("CONNECT %s HTTP/1.0\n\n", service.DefaultRPCPath))

	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	if err == nil && resp.Status == service.Connected {
		return NewClient(conn, opt)
	}
	if err == nil {
		err = fmt.Errorf("unexpected HTTP status: " + resp.Status)
	}
	return nil, err
}

func DialHTTP(network, addr string, opts ...*service.Option) (*Client, error) {
	return dialTimeout(NewHTTPClient, network, addr, opts...)
}

func XDial(rpcAddr string, opts ...*service.Option) (*Client, error) {
	parts := strings.Split(rpcAddr, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("rpc client: wrong format '%s', expected protocal@addr", rpcAddr)
	}
	protocal, addr := parts[0], parts[1]
	switch protocal {
	case "http":
		return DialHTTP("tcp", addr, opts...)
	default:
		return Dial("tcp", addr, opts...)
	}
}
