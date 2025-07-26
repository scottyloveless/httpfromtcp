package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddr := "localhost:42069"

	udpAddress, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving UDP address: %v\n", err)
		os.Exit(1)
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dialing UDP: %v\n", err)
		os.Exit(1)
	}
	defer udpConn.Close()

	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit\n", udpAddress)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
			os.Exit(1)
		}
		_, err = udpConn.Write([]byte(message))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error sending message: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Message sent: %s", message)
	}
}
