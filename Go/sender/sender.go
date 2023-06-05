package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

func main() {
	CHUNK_SIZE := 16

	// Create a UDP socket
	udpSocket := createUDPSocket()
	defer udpSocket.Close()

	// Transfer the file to bytes
	bytes := fileToBytes("send.zip")

	// Create the payload from the bytes
	payload := payloadSeparation(bytes, CHUNK_SIZE)

	// Starts the sending mechanism
	slowStart(udpSocket, payload)
}

func createUDPSocket() *net.UDPConn {
	// Server address
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Failed to resolve server address: %s\n", err)
		return nil
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Printf("Failed to create UDP connection: %s\n", err)
		return nil
	}

	return conn
}

func fileToBytes(filePath string) []byte {
	// Read the contents of the .txt file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return nil
	}

	// Convert the content to bytes
	bytes := []byte(content)

	return bytes
}

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

	fmt.Printf("Payload: %v\n", matrix)

	return matrix
}

func sender(udpSocket *net.UDPConn, payload [][]byte) {
	size := fmt.Sprint(len(payload))
	udpSocket.Write([]byte(size))

	for i := 0; i < len(payload); i++ {
		// Data to send
		message := payload[i]

		// Send the UDP packet
		udpSocket.Write(message)

		fmt.Println("\nUDP packet sent successfully!")

		// Buffer to store ACK message
		ackBuffer := make([]byte, 1024)

		// Receive the ACK from the receiver
		n, _, err := udpSocket.ReadFromUDP(ackBuffer)
		if err != nil {
			fmt.Printf("Failed to receive ACK: %s\n", err)
			return
		}

		// Print the received ACK message
		fmt.Printf("Received ACK: %s\n", string(ackBuffer[:n]))
	}
}

func slowStart(udpSocket *net.UDPConn, payload [][]byte) {
	size := fmt.Sprint(len(payload))
	udpSocket.Write([]byte(size))

	packetsToSend := 1
	packetsSent := 0
	caCounter := 4
	for i := 0; i < len(payload); i++ {
		// Data to send
		message := payload[i]

		// Send the UDP packet
		udpSocket.Write(message)
		packetsSent++

		fmt.Println("\nUDP packet sent successfully!")

		if packetsSent == packetsToSend {
			for j := 0; j < packetsSent; j++ {
				ack(udpSocket)
			}

			if packetsToSend == caCounter {
				packetsSent = 0
				packetsToSend++
				caCounter++
			} else {
				packetsSent = 0
				packetsToSend *= 2
			}
		} else if (packetsSent < packetsToSend) && (i == len(payload)-1) {
			for j := 0; j < packetsSent; j++ {
				ack(udpSocket)
			}
		}
	}
}

func ack(udpSocket *net.UDPConn) {
	// Buffer to store ACK message
	ackBuffer := make([]byte, 1024)

	// Receive the ACK from the receiver
	n, _, err := udpSocket.ReadFromUDP(ackBuffer)
	if err != nil {
		fmt.Printf("Failed to receive ACK: %s\n", err)
		return
	}

	// Print the received ACK message
	fmt.Printf("Received ACK: %s\n", string(ackBuffer[:n]))
}
