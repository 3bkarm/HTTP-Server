package server

import (
	"fmt"
	"httpserver/internal/headers"
	"httpserver/internal/request"
	"httpserver/internal/response"
	"io"
	"net"
)

type Handler func(w *response.Writer, rq *request.Request)

type Server struct {
	closed  bool
	handler Handler
}

func runConnection(srv *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	responseWriter := response.NewWriter(conn)
	rq, err := request.RequestFromReader(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		responseWriter.WriteHeaders(headers.GetDefaultHeaders(0))
		return
	}
	srv.handler(responseWriter, rq)
}

func runServer(srv *Server, listener net.Listener) {
	for {
		// TODO: explore more about this case
		if srv.closed {
			return
		}
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		go runConnection(srv, conn)
	}
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	srv := &Server{
		closed:  false,
		handler: handler,
	}
	go runServer(srv, listener)
	return srv, nil
}

func (srv *Server) Close() error {
	srv.closed = true
	return nil
}
