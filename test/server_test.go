package test

import (
	"context"
	client2 "gorpc/client"
	service2 "gorpc/service"
	"net"
	"strings"
	"testing"
	"time"
)

type Bar int

func (b Bar) Timeout(argv int, reply *int) error {
	*reply = argv
	time.Sleep(time.Second * 2)
	return nil
}

func startServer(addr chan string) {
	var b Bar
	_ = service2.Register(&b)
	l, _ := net.Listen("tcp", ":0")
	addr <- l.Addr().String()
	service2.Accept(l)
}

func TestClient_Call(t *testing.T) {
	t.Parallel()
	addrCh := make(chan string)
	go startServer(addrCh)
	addr := <-addrCh
	time.Sleep(time.Second)
	t.Run("client timeout", func(t *testing.T) {
		client, _ := client2.Dial("tcp", addr)
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		var reply int
		err := client.Call(ctx, "Bar.Timeout", 1, &reply)
		_assert(err != nil && strings.Contains(err.Error(), ctx.Err().Error()), "expect a timeout error")
	})

	t.Run("server handle timeout", func(t *testing.T) {
		client, _ := client2.Dial("tcp", addr, &service2.Option{
			HandleTimeout: time.Second * 1,
		})
		var reply int
		err := client.Call(context.Background(), "Bar.Timeout", 1, &reply)
		_assert(err != nil, "expect a timeout error")
	})
}
