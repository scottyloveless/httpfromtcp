package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening on port: %s %s\n", port, err)
	}
	defer listener.Close()

	fmt.Println("Reading data from listener...")
	fmt.Println("=====================================")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error establishing connection: %v\n", err)
		}
		fmt.Println("connection accepted")
		linesChan := getLinesChannel(conn)
		for line := range linesChan {
			fmt.Println(line)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer fmt.Println("connection closed")
		defer close(lines)
		currentLineContents := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					lines <- fmt.Sprint(currentLineContents)
					currentLineContents = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				break
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}
