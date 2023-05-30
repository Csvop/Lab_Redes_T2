package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

func main() {
	UDP_IP := "127.0.0.1"
	UDP_PORT := 34754
	CHUNK_SIZE := 1024

	// Transfer the file to bytes
	payload := fileToBytes("send.txt")

	// Create a UDP socket
	udpSocket := createUDPSocket(UDP_IP, UDP_PORT)

	// Send the payload
	send(udpSocket, payload, CHUNK_SIZE)

	// Close the socket
	udpSocket.Close()
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

	return bytes
}

func createUDPSocket(UDP_IP string, UDP_PORT int) *net.UDPConn {
	// Create the address
	UDP_ADDR := net.UDPAddr{
		Port: UDP_PORT,
		IP:   net.ParseIP(UDP_IP),
	}

	// Create the UDP socket
	udpSocket, err := net.DialUDP("udp", nil, &UDP_ADDR)
	if err != nil {
		fmt.Printf("Error creating UDP socket: %s\n", err)
		return nil
	}

	return udpSocket
}

func send(udpSocket *net.UDPConn, payload []byte, CHUNK_SIZE int) {
	_, err := udpSocket.Write(payload)
	if err != nil {
		fmt.Printf("Error sending payload: %s\n", err)
		return
	}
	fmt.Printf("Payload sent successfully\n")
}
