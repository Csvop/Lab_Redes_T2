package main

import (
	"fmt"
	"net"
)

func main() {
	UDP_IP := "127.0.0.1"
	UDP_PORT := 5001

	// Create a UDP socket
	udpSocket := createUDPSocket(UDP_IP, UDP_PORT)

	// Listen for incoming packets
	go packetListener(udpSocket)

	// Close the socket
	udpSocket.Close()
}

func createUDPSocket(UDP_IP string, UDP_PORT int) *net.UDPConn {
	// Resolve the UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", UDP_IP, UDP_PORT))
	if err != nil {
		fmt.Printf("Error resolving UDP address: %s\n", err)
		return nil
	}

	// Create a UDP socket
	udpSocket, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Printf("Error listening to UDP: %s\n", err)
		return nil
	}

	return udpSocket
}

func packetListener(udpSocket *net.UDPConn) {
	// Listen for incoming packets
	for {
		// Create a buffer for the incoming packet
		buffer := make([]byte, 16)

		// Read the incoming packet
		n, _, err := udpSocket.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP: %s\n", err)
			return
		}

		// Print the packet
		fmt.Printf("Packet received: %s\n", buffer[:n])
	}
}
