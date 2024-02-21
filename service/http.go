package service

import (
	"io"
	"log"
	"net/http"
)

const (
	Connected        = "200 Connected to go RPC"
	DefaultRPCPath   = "/_gorpc_"
	DefaultDebugPath = "/debug/gorpc"
)

func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = io.WriteString(w, "405 must CONNECT\n")
		return
	}
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Printf("rpc hijacking %s: %s", req.RemoteAddr, err.Error())
	}
	_, _ = io.WriteString(conn, "HTTP/1.0 "+Connected+"\n\n")
	server.ServeConn(conn)
}

func (server *Server) HandleHTTP() {
	http.Handle(DefaultRPCPath, server)
	http.Handle(DefaultDebugPath, debugHTTP{server})
	log.Println("rpc server debug path", DefaultDebugPath)
}

func HandleHTTP() {
	DefaultServer.HandleHTTP()
}
