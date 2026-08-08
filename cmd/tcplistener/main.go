package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Failed to open listener: %v", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection: %v", err)
		}

		fmt.Println("Connection has been accepted")
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("read: %s\n", line)
		}
		fmt.Println("Connection has been closed")
	}
}

func getLinesChannel(conn net.Conn) <-chan string {
	ch := make(chan string)

	go func() {
		defer conn.Close()
		defer close(ch)

		buffer := make([]byte, 8)
		line := ""

		for {
			bytesRead, err := conn.Read(buffer)
			if err != nil {
				if line != "" {
					ch <- line
				}
				if err == io.EOF {
					break
				}
				log.Printf("Error during read: %v", err)
				break
			}

			chunk := buffer[:bytesRead]
			parts := bytes.Split(chunk, []byte("\n"))

			for i := 0; i < len(parts)-1; i++ {
				ch <- line + string(parts[i])
				line = ""
			}
			line += string(parts[len(parts)-1])
		}
	}()

	return ch
}
