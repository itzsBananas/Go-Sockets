package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatalf("Client cannot initialize: %s", err)
	}
	defer conn.Close()

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
