package message

import (
	"encoding/binary"
)

type Question struct {
	Labels []string
	Record uint16
	Class  uint16
}

func DecodeQuestion(q []byte, startIndex uint16) (Question, uint16, error) {
	index := startIndex
	question := Question{}
	labels := []string{}

	var compressedIndex uint16

	for {
		// exit on null byte
		initialByte := q[index]
		if initialByte == 0 {
			if compressedIndex > 0 {
				index = compressedIndex
				compressedIndex = 0
				continue
			}
			break
		}

		// if 2 msb are enabled then content is compressed
		compressed := (initialByte & 192) == 192
		if compressed {
			compressedIndex = index + 2
			offset := binary.BigEndian.Uint16([]byte{initialByte & 63, q[index+1]})
			index = offset
			continue
		}

		len := initialByte

		index = index + 1

		content := q[index : index+uint16(len)]
		labels = append(labels, string(content))

		index = index + uint16(len)

	}

	question.Labels = labels

	index = index + 1

	// get record
	r1 := q[index]

	index = index + 1
	r2 := q[index]

	question.Record = binary.BigEndian.Uint16([]byte{r1, r2})

	// get class
	index = index + 1
	c1 := q[index]

	index = index + 1
	c2 := q[index]

	question.Class = binary.BigEndian.Uint16([]byte{c1, c2})

	return question, index + 1, nil
}

func (q Question) Encode() []byte {
	bs := []byte{}

	for _, label := range q.Labels {
		bs = append(bs, uint8(len(label)))
		bs = append(bs, []byte(label)...)
	}

	// null byte after labels
	bs = append(bs, uint8(0))

	bs = append(bs, uint8(q.Record>>8))
	bs = append(bs, uint8(q.Record))

	bs = append(bs, uint8(q.Class>>8))
	bs = append(bs, uint8(q.Class))

	return bs
}
