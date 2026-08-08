package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// start listerner: nc -u -l 42069 (may need -k if machine close after fine messages revicieved)
// Sender: go run ./cmd/udpsender

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Failed to resolve udp address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("err: %v", err)
		}

		conn.Write([]byte(line))
	}
}
