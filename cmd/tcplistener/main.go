package main

import (
	"fmt"
	"httpserver/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		rq, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", rq.RequestLine.Method)
		fmt.Printf("- Target: %s\n", rq.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", rq.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		rq.Headers.ForEach(func(k, v string) {
			fmt.Printf("- %s: %s\n", k, v)
		})
		fmt.Printf("Body:\n")
		fmt.Printf("%s\n", rq.Body)
	}
}
