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

		receivedHeader, err := message.DecodeHeader(buf)
		if err != nil {
			fmt.Println("Failed to decode header:", err)
		}

		responseCode := 0
		if receivedHeader.OpCode != 0 {
			responseCode = 4
		}

		head := message.Header{
			ID:                    receivedHeader.ID,
			Query:                 true,
			OpCode:                receivedHeader.OpCode,
			AuthorativeAnswer:     false,
			Truncation:            false,
			RecursionDesired:      receivedHeader.RecursionDesired,
			RecursionAvailable:    false,
			Reserved:              0,
			ResponseCode:          uint8(responseCode),
			QuestionCount:         1,
			AnswerRecordCount:     1,
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

		answer := message.Answer{
			Labels: []string{
				"codecrafters", "io",
			},
			Record:   1,
			Class:    1,
			Ttl:      60,
			RDLength: 4,
			RData: []byte{
				8, 8, 8, 8,
			},
		}

		hb := head.Encode()
		qb := question.Encode()
		ab := answer.Encode()

		by := []byte{}
		by = append(by, hb...)
		by = append(by, qb...)
		by = append(by, ab...)

		_, err = udpConn.WriteToUDP(by, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
