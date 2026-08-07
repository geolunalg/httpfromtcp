package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		buffer := make([]byte, 8)
		line := ""

		for {
			bytesRead, err := f.Read(buffer)
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
