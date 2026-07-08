package message

import (
	"encoding/binary"
	"fmt"
)

type Header struct {
	ID                    uint16 // Packet Identifier
	Query                 bool   // Query or response indicator
	OpCode                uint8  // Operation code (4 bits)
	AuthorativeAnswer     bool   // Does the DNS own the domain
	Truncation            bool   // is the message larger than 512 bytes
	RecursionDesired      bool   // should we recursively resolve the query
	RecursionAvailable    bool   // do we support recursion
	Reserved              uint8  // 3 bits
	ResponseCode          uint8  // 4 bits
	QuestionCount         uint16 // number of questions
	AnswerRecordCount     uint16 // number of answers
	AuthorityRecordCount  uint16 // number of authorities
	AdditionalRecordCount uint16 // number of records in additional section
}

func DecodeHeader(bytes []byte) (Header, error) {

	if len(bytes) < 12 {
		return Header{}, fmt.Errorf("invalid length of bytes %d", len(bytes))
	}

	id := binary.BigEndian.Uint16(bytes[0:2])

	queryMask := uint8(128)
	query := (bytes[2] & queryMask) > 0

	opcodeMask := uint8(120)
	opcode := (bytes[2] & opcodeMask) >> 3

	authorativeAnswerMask := uint8(4)
	authorativeAnswer := (bytes[2] & authorativeAnswerMask) > 0

	truncationMask := uint8(2)
	truncation := (bytes[2] & truncationMask) > 0

	recursionDesiredMask := uint8(1)
	recursionDesired := (bytes[2] & recursionDesiredMask) > 0

	recursionAvailableMask := uint8(128)
	recursionAvailable := (bytes[3] & recursionAvailableMask) > 0

	reserved := uint8(0)

	responseCodeMask := uint8(15)
	responseCode := (bytes[3] & responseCodeMask)

	numQuestions := binary.BigEndian.Uint16(bytes[4:6])

	numAnswerRecord := binary.BigEndian.Uint16(bytes[6:8])

	numAuthorityRecord := binary.BigEndian.Uint16(bytes[8:10])

	numAdditionalRecord := binary.BigEndian.Uint16(bytes[10:12])

	return Header{
		ID:                    id,
		Query:                 query,
		OpCode:                opcode,
		AuthorativeAnswer:     authorativeAnswer,
		Truncation:            truncation,
		RecursionDesired:      recursionDesired,
		RecursionAvailable:    recursionAvailable,
		Reserved:              reserved,
		ResponseCode:          responseCode,
		QuestionCount:         numQuestions,
		AnswerRecordCount:     numAnswerRecord,
		AuthorityRecordCount:  numAuthorityRecord,
		AdditionalRecordCount: numAdditionalRecord,
	}, nil
}

func (h Header) Encode() []byte {
	headerBytes := make([]byte, 12)

	headerBytes[0] = uint8(h.ID >> 8)
	headerBytes[1] = uint8(h.ID)

	b3 := boolToUint8(h.Query, 7)
	b3 = b3 | (h.OpCode << 3)
	b3 = b3 | boolToUint8(h.AuthorativeAnswer, 2)
	b3 = b3 | boolToUint8(h.Truncation, 1)
	b3 = b3 | boolToUint8(h.RecursionDesired, 0)
	headerBytes[2] = b3

	b4 := boolToUint8(h.RecursionAvailable, 7)
	b4 = b4 | (h.Reserved << 4)
	b4 = b4 | h.ResponseCode
	headerBytes[3] = b4

	headerBytes[4] = uint8(h.QuestionCount >> 8)
	headerBytes[5] = uint8(h.QuestionCount)

	headerBytes[6] = uint8(h.AnswerRecordCount >> 8)
	headerBytes[7] = uint8(h.AnswerRecordCount)

	headerBytes[8] = uint8(h.AuthorityRecordCount >> 8)
	headerBytes[9] = uint8(h.AuthorityRecordCount)

	headerBytes[10] = uint8(h.AdditionalRecordCount >> 8)
	headerBytes[11] = uint8(h.AdditionalRecordCount)

	return headerBytes
}

func boolToUint8(b bool, shift uint8) uint8 {
	if b {
		return uint8(1) << shift
	}
	return uint8(0)
}
