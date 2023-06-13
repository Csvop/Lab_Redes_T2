package main

import (
	"fmt"
	"net"
	"strconv"
)

func getSequenceNumber(buffer []byte) int {
	// Get the sequence number from the packet
	sequenceNumber := 0
	for i := 0; i < len(buffer); i++ {
		if buffer[i] == ';' {
			break
		}
		sequenceNumber = sequenceNumber*10 + int(buffer[i]-'0')
	}

	return sequenceNumber
}

func main() {
	// Size of each chunk
	chunk_size := 400

	// Local address to listen on
	localAddr, _ := net.ResolveUDPAddr("udp", ":1234")

	// Create a UDP connection to listen for incoming packets
	conn, _ := net.ListenUDP("udp", localAddr)

	// Sync with sender and receive the payload size
	buffer := make([]byte, 1024)
	payloadSize := 0
	for {
		_, addr, _ := conn.ReadFromUDP(buffer)
		response, _ := strconv.Atoi(string(buffer))

		if response != 0 {
			// Sets the payload size
			payloadSize = response

			// Confirms sync with sender and sends ack
			conn.WriteToUDP([]byte("CONFIRMED"), addr)
			break
		}
	}

	// Receives incoming packets
	i := 0
	for {
		// Buffer for incoming packets
		buffer := make([]byte, chunk_size)

		// Read incoming packet into buffer
		_, addr, _ := conn.ReadFromUDP(buffer)

		// Print the received packet
		println("Received from %d", getSequenceNumber(buffer))

		// Send a response back to sender
		conn.WriteToUDP([]byte("ACK"), addr)

		// Increment the counter
		i++

		// If the payload is complete, stop
		if payloadSize == i {
			fmt.Printf("\n\nFIN\n\n")
			break
		}
	}
}
