package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

func main() {
	udpSocket := createUDPSocket()
	defer udpSocket.Close()

	bytes := reciver(udpSocket)

	// Write the bytes to the file
	err := ioutil.WriteFile("send.txt", bytes, 0644)
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
	// Receive the first packet
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Printf("Failed to receive UDP packet: %s\n", err)
		return nil
	}

	fmt.Printf("Received %d bytes from %s: %s\n", n, addr.String(), string(buffer[:n]))

	intVar, _ := strconv.Atoi(string(buffer[:n]))

	byteBuffer := make([][]byte, intVar)

	ackAtual := 1

	for i := 0; i < intVar; i++ {
		// Buffer to store received data
		buffer := make([]byte, 1024)

		// Receive the UDP packet
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Failed to receive UDP packet: %s\n", err)
			return nil
		}

		dto := strings.Split(string(buffer[:n]), ";")
		packetIndex, _ := strconv.Atoi(dto[0])
		isRetransmission, _ := strconv.Atoi(dto[1])
		message := getMessageFromDto(dto)

		// Print the received packet details
		fmt.Printf("\nReceived %d bytes from %s: %s\n", n, addr.String(), message)
		fmt.Println("Packet index: " + strconv.Itoa(packetIndex))

		// Sleep for 2 seconds to simulate a long-running process
		// time.Sleep(2 * time.Second)

		if packetIndex == ackAtual {
			ackAtual++

			if isRetransmission == 1 {
				ackAtual += 3
			}
		}

		// Send the ACK back to the sender
		_, err = conn.WriteToUDP([]byte(strconv.Itoa(ackAtual)), addr)
		if err != nil {
			fmt.Printf("Failed to send ACK: %s\n", err)
			return nil
		}

		fmt.Println("ACK=" + strconv.Itoa(ackAtual))

		byteBuffer[packetIndex] = []byte(message)
	}

	sendBuffer := cleanMatrix(byteBuffer)

	return sendBuffer
}

func cleanMatrix(matrixBuffer [][]byte) []byte {
	var byteBuffer []byte

	buffer := matrixToArray(matrixBuffer)

	// Clean the array starting from the end until the first non-zero byte
	for i := len(buffer) - 1; i >= 0; i-- {
		if buffer[i] != 0 {
			byteBuffer = buffer[:i+1]
			break
		}
	}

	return byteBuffer
}

func matrixToArray(matrixBuffer [][]byte) []byte {
	var byteBuffer []byte

	for i := 0; i < len(matrixBuffer); i++ {
		for j := 0; j < len(matrixBuffer[i]); j++ {
			byteBuffer = append(byteBuffer, matrixBuffer[i][j])
		}
	}

	return byteBuffer
}

func getMessageFromDto(dto []string) string {
	message := ""
	for i := 2; i < len(dto); i++ {
		if i != 2 {
			message += (";" + dto[i])
		} else {
			message += dto[i]
		}

	}

	return message
}
