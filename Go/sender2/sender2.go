package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

func payloadSeparation(payload []byte, partSize int) [][]byte {
	array := make([]byte, 0)
	array = append(array, payload...)

	if len(array)%partSize != 0 {
		filler := partSize - (len(array) % partSize)

		for i := 0; i < filler; i++ {
			array = append(array, byte(0))
		}
	}

	quantPackets := len(array) / partSize

	arrayIndex := 0
	matrix := make([][]byte, quantPackets)
	for i := 0; i < quantPackets; i++ {
		matrix[i] = make([]byte, partSize)
		for j := 0; j < partSize; j++ {
			matrix[i][j] = array[arrayIndex]
			arrayIndex++
		}
	}

	return matrix
}

func sender(payload [][]byte, i int, conn net.Conn, payloadSize int) {
	time.Sleep(2 * time.Second)
	for {
		// Adds the sequence number to the payload
		payload[i] = append([]byte(fmt.Sprintf("%d", i+1)+";"), payload[i]...)

		// Send the UDP packet
		fmt.Printf("Sending %d/%d\n", i+1, payloadSize)
		conn.Write(payload[i])

		// Wait for the server to respond for 1 second
		buffer := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		response, _ := conn.Read(buffer)

		// If the server responded, print the response
		if response != 0 {
			fmt.Printf("Response: %s for %d\n", string(buffer), i+1)
			break
		}
	}
}

func main() {
	// Channel to stop the main thread from exiting
	stop := make(chan bool)

	// Size of each chunk
	chunk_size := 300

	// Server address
	serverAddr, _ := net.ResolveUDPAddr("udp", "localhost:1234")

	// Create a UDP connection
	conn, _ := net.DialUDP("udp", nil, serverAddr)
	defer conn.Close()

	// Read the contents of the file
	content, _ := ioutil.ReadFile("simple.txt")

	// Convert the content to bytes
	fileInBytes := []byte(content)

	// Separate the payload into chunks
	payload := payloadSeparation(fileInBytes, chunk_size)

	// Gets the size of the payload
	payloadSize := len(payload)

	// Starts sync with receiver and send the size of the payload to the server
	conn.Write([]byte(fmt.Sprintf("%d", payloadSize)))

	// Wait for the server to confirm the sync
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	if string(buffer) == "CONFIRMED" {
		fmt.Printf("Sync confirmed\n")
	}

	// Send all the packages on the payload
	for i := 0; i < payloadSize; i++ {
		// Creates a Go routine to send the package
		go sender(payload, i, conn, payloadSize)
	}

	// Wait for the stop signal
	<-stop
}
