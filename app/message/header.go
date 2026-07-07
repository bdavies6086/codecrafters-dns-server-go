package message

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

func (h Header) Encode() []byte {
	headerBytes := make([]byte, 12)

	b1 := uint8(h.ID)
	b2 := uint8(h.ID << 8)

	b3 := boolToUint8(h.Query, 0)
	b3 = b3 | (h.OpCode >> 1)
	b3 = b3 | boolToUint8(h.AuthorativeAnswer, 5)
	b3 = b3 | boolToUint8(h.Truncation, 6)
	b3 = b3 | boolToUint8(h.RecursionDesired, 7)

	b4 := boolToUint8(h.RecursionAvailable, 0)
	b4 = b4 | (h.Reserved >> 1)
	b4 = b4 | (h.ResponseCode >> 4)

	b5 := uint8(h.QuestionCount)
	b6 := uint8(h.QuestionCount << 8)

	b7 := uint8(h.AnswerRecordCount)
	b8 := uint8(h.AnswerRecordCount << 8)

	b9 := uint8(h.AuthorityRecordCount)
	b10 := uint8(h.AuthorityRecordCount << 8)

	b11 := uint8(h.AdditionalRecordCount)
	b12 := uint8(h.AdditionalRecordCount << 8)

	return append(headerBytes, b1, b2, b3, b4, b5, b6, b7, b8, b9, b10, b11, b12)

}

func boolToUint8(b bool, shift uint8) uint8 {
	if b {
		return uint8(1) >> shift
	}
	return uint8(0)
}
