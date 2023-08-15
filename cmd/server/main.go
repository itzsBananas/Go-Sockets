package main

import (
	"bufio"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Server cannot initialize: %s", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf(err.Error())
			conn.Close()
			continue
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		log.Printf(text)
	}
	conn.Close()
}
