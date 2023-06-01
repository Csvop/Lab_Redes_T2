package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

func main() {
	stop := make(chan bool)
	packetReceived := make(chan int)

	UDP_IP := "127.0.0.1"
	UDP_PORT := 5000
	// UDP_PORT_TO_SEND = 5001
	CHUNK_SIZE := 16

	// Create a UDP socket
	udpSocket := createUDPSocket(UDP_IP, UDP_PORT)

	// Listen for incoming acks
	go ackListener(packetReceived, udpSocket)

	// Transfer the file to bytes
	bytes := fileToBytes("send.txt")

	// Create the payload from the bytes
	payload := payloadSeparation(bytes, CHUNK_SIZE)

	// Send the payload
	go sender(udpSocket, payload, CHUNK_SIZE, stop, packetReceived)

	// Close the socket
	udpSocket.Close()
	<-stop
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

	// Use the bytes as needed
	fmt.Printf("Bytes: %v\n", bytes)
	fmt.Printf("Bytes as string: %d\n", len(bytes))

	return bytes
}

func createUDPSocket(UDP_IP string, UDP_PORT int) *net.UDPConn {
	// Resolve the UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", UDP_IP, UDP_PORT))
	if err != nil {
		fmt.Printf("Error resolving UDP address: %s\n", err)
		return nil
	}

	// Create a UDP socket
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Printf("Error creating UDP socket: %s\n", err)
		return nil
	}

	return udpConn
}

func sender(udpSocket *net.UDPConn, payload [][]byte, CHUNK_SIZE int, stop chan bool, packetReceived chan int) {
	currentPacket := 0
	for {
		if currentPacket == len(payload) {
			break
		}

		fmt.Printf("\nSending payload...\n")
		fmt.Printf("Payload: %v\n", payload[currentPacket])
		send(udpSocket, payload[currentPacket], packetReceived)

		// Wait for ack
		fmt.Printf("\nWaiting for ack...\n")
		<-packetReceived

		currentPacket = currentPacket + 1
	}
	stop <- true
}

func send(udpSocket *net.UDPConn, packet []byte, packetReceived chan int) {
	_, err := udpSocket.Write(packet)
	if err != nil {
		fmt.Printf("Error sending payload: %s\n", err)
		return
	}
	fmt.Printf("\nPayload sent successfully\n")
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

func ackListener(packetReceived chan int, udpSocket *net.UDPConn) {
	for {
		// // Listen for incoming acks
		// udpAddr, err := udpSocket.Read(payload)
		// if err != nil {
		// 	fmt.Printf("Error resolving UDP address: %s\n", err)
		// 	return
		// }

		// udpConn, err := net.ListenUDP("udp", udpAddr)

		fmt.Printf("\nPacket Ack: Sucess!")
		packetReceived <- 1
	}
}
