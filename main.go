package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

var (
	counter int

	listenAddr = "localhost:8081"

	server = []string{
		"localhost:5001",
		"localhost:5002",
		"localhost:5003",
	}
)

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connetion: %s", err)
		}

		backend := selectBackend()
		fmt.Printf("counter:%d backend:%s\n", counter, backend)
		go func() {
			proxy(backend, conn)
			if err != nil {
				log.Printf("wARNING: proxing failed: %v", err)
			}
		}()

	}

}

func proxy(backend string, c net.Conn) error {

	bc, err := net.Dial("tcp", backend)
	if err != nil {
		return fmt.Errorf("failed to connect to backend %s: %v", backend, err)
	}
	go io.Copy(bc, c)
	go io.Copy(c, bc)

	return nil
}

func selectBackend() string {
	s := server[counter%len(server)]
	counter++
	return s
}
