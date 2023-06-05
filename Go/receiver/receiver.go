package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
)

func main() {
	udpSocket := createUDPSocket()
	defer udpSocket.Close()

	bytes := reciver(udpSocket)

	// Write the bytes to the file
	err := ioutil.WriteFile("receive.zip", bytes, 0644)
	if err != nil {
		fmt.Printf("Failed to write to file: %s\n", err)
		return
	}
}

func createUDPSocket() *net.UDPConn {
	// Local address to listen on
	localAddr, err := net.ResolveUDPAddr("udp", ":1234")
	if err != nil {
		fmt.Printf("Failed to resolve local address: %s\n", err)
		return nil
	}

	// Create a UDP connection to listen for incoming packets
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Printf("Failed to create UDP connection: %s\n", err)
		return nil
	}

	return conn
}

func reciver(conn *net.UDPConn) []byte {
	stringBuffer := ""

	// Receive the first packet
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Printf("Failed to receive UDP packet: %s\n", err)
		return nil
	}

	fmt.Printf("Received %d bytes from %s: %s\n", n, addr.String(), string(buffer[:n]))

	intVar, _ := strconv.Atoi(string(buffer[:n]))

	for i := 0; i < intVar; i++ {
		// Buffer to store received data
		buffer := make([]byte, 1024)

		// Receive the UDP packet
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Failed to receive UDP packet: %s\n", err)
			return nil
		}

		// Print the received packet details
		fmt.Printf("\nReceived %d bytes from %s: %s\n", n, addr.String(), string(buffer[:n]))

		// Sleep for 2 seconds to simulate a long-running process
		// time.Sleep(2 * time.Second)

		// Send the ACK back to the sender
		ackMessage := []byte(string(buffer[:n]))
		_, err = conn.WriteToUDP(ackMessage, addr)
		if err != nil {
			fmt.Printf("Failed to send ACK: %s\n", err)
			return nil
		}

		fmt.Println("ACK sent!")

		stringBuffer = stringBuffer + string(buffer[:n])
	}

	byteBuffer := cleanArray([]byte(stringBuffer))

	return byteBuffer
}

func cleanArray(stringBuffer []byte) []byte {
	var byteBuffer []byte

	// Clean the array starting from the end until the first non-zero byte
	for i := len(stringBuffer) - 1; i >= 0; i-- {
		if stringBuffer[i] != 0 {
			byteBuffer = stringBuffer[:i+1]
			break
		}
	}

	return byteBuffer
}
