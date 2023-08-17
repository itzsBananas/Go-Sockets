package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var port = flag.Int("p", 8080, "client port number")

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Client cannot initialize: %s", err)
	}
	defer conn.Close()

	go handleConnection(conn)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text() + "\n"
		_, err := conn.Write([]byte(text))
		if err != nil {
			log.Println(err)
		}
	}
	os.Stdin.Close()
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
