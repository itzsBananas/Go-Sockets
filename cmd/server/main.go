package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type message struct {
	src  net.Conn
	data string
}

var chat = make(chan message)
var enter, leave = make(chan net.Conn), make(chan net.Conn)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Server cannot initialize: %s", err)
	}
	go acceptConnections(ln)

	users := make(map[net.Conn]struct{})
	for {
		select {
		case conn := <-enter:
			users[conn] = struct{}{}
		case conn := <-leave:
			delete(users, conn)
		case msg := <-chat:
			log.Println(msg.data)
			for conn := range users {
				if msg.src != conn {
					sendMessage(conn, []byte(msg.data))
				}
			}
		}
	}
}

func acceptConnections(ln net.Listener) {
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
	msg := message{
		src: conn,
	}
	enter <- conn
	msg.data = fmt.Sprintf("%s has entered the chat", name)
	chat <- msg

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
		msg.data = fmt.Sprintf("%s: %s", name, scanner.Text())
		chat <- msg
	}

	leave <- conn
	msg.data = fmt.Sprintf("%s has left the chat", name)
	chat <- msg
	conn.Close()
}

func sendMessage(conn net.Conn, message []byte) {
	message = append(message, '\n')
	_, err := conn.Write(message)
	if err != nil {
		log.Println(err)
	}
}
