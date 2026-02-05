package main

import (
	"fmt"
	"httpserver/internal/headers"
	"httpserver/internal/request"
	"httpserver/internal/response"
	"httpserver/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	srv, err := server.Serve(port, func(w *response.Writer, rq *request.Request) {
		status := response.StatusOk
		h := headers.GetDefaultHeaders(0)
		body := response.OkRespond()
		if rq.RequestLine.RequestTarget == "/yourproblem" {
			status = response.StatusBadRequest
			body = response.BadRequestRespond()
		} else if rq.RequestLine.RequestTarget == "/myproblem" {
			status = response.StatusInternalServerError
			body = response.InternalServerErrorRespond()
		}
		w.WriteStatusLine(status)
		h.Replace("content-length", fmt.Sprintf("%d", len(body)))
		h.Replace("content-type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody(body)
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer srv.Close()
	log.Println("Server started on port", port)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
