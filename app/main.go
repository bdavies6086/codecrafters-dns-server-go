package main

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/app/message"
)

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		head := message.Header{
			ID:                    1234,
			Query:                 true,
			OpCode:                0,
			AuthorativeAnswer:     false,
			Truncation:            false,
			RecursionDesired:      false,
			RecursionAvailable:    false,
			Reserved:              0,
			ResponseCode:          0,
			QuestionCount:         1,
			AnswerRecordCount:     0,
			AuthorityRecordCount:  0,
			AdditionalRecordCount: 0,
		}

		question := message.Question{
			Labels: []string{
				"codecrafters", "io",
			},
			Record: 1,
			Class:  1,
		}

		hb := head.Encode()
		qb := question.Encode()

		by := []byte{}
		by = append(by, hb...)
		by = append(by, qb...)

		_, err = udpConn.WriteToUDP(by, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
