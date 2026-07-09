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
			fmt.Println("Failed to decode header")
			break
		}

		questions := []message.Question{}
		questionOffset := uint16(12)

		for i := uint16(0); i < receivedHeader.QuestionCount; i++ {

			// for now handle a single question
			// we need to send the whole buffer in because compression means we offset from the header
			q, qo, err := message.DecodeQuestion(buf, questionOffset)
			if err != nil {
				fmt.Printf("Failed to decode question")
				break
			}
			questionOffset = qo
			questions = append(questions, q)
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
			QuestionCount:         uint16(len(questions)),
			AnswerRecordCount:     uint16(len(questions)),
			AuthorityRecordCount:  0,
			AdditionalRecordCount: 0,
		}

		by := []byte{}
		hb := head.Encode()
		by = append(by, hb...)

		for _, q := range questions {
			qb := q.Encode()
			by = append(by, qb...)
		}

		for _, q := range questions {
			answer := message.Answer{
				Labels:   q.Labels,
				Record:   q.Record,
				Class:    q.Class,
				Ttl:      60,
				RDLength: 4,
				RData: []byte{
					8, 8, 8, 8,
				},
			}

			ab := answer.Encode()
			by = append(by, ab...)

		}

		_, err = udpConn.WriteToUDP(by, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
