package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
)

var (
	users = make(map[net.Conn]struct{})
	mu    sync.RWMutex
)

var port = flag.Int("p", 8080, "client port number")

func main() {
	flag.Parse()
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
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
	name := conn.RemoteAddr().String()
	mu.Lock()
	users[conn] = struct{}{}
	mu.Unlock()

	mu.RLock()
	broadcast(conn, fmt.Sprintf("%s has entered the chat", name))
	mu.RUnlock()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
		mu.RLock()
		broadcast(conn, fmt.Sprintf("%s: %s", name, scanner.Text()))
		mu.RUnlock()
	}

	mu.Lock()
	delete(users, conn)
	mu.Unlock()
	mu.RLock()
	broadcast(conn, fmt.Sprintf("%s has left the chat", name))
	mu.RUnlock()
	conn.Close()
}

func broadcast(src net.Conn, message string) {
	log.Println(message)
	for conn := range users {
		if src != conn {
			go sendMessage(conn, []byte(message))
		}
	}
}

func sendMessage(conn net.Conn, message []byte) {
	message = append(message, '\n')
	_, err := conn.Write(message)
	if err != nil {
		log.Println(err)
	}
}
