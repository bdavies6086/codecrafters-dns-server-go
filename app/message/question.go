package message

type Question struct {
	Labels []string
	Record uint16
	Class  uint16
}

// func Decode(q []byte) (Question, error) {
// 	bs := []byte{}

// 	index := 0
// 	question := Question{}
// 	labels := []string{}

// 	for {
// 		// exit on null byte
// 		len := q[index]
// 		if len == 0 {
// 			break
// 		}

// 		content := q[index+1 : index+len]
// 		labels = append(labels, string(content))
// 		index = index + len
// 	}

// 	question.Labels = labels

// 	// get record
// 	index = index + 1
// 	r1 := q[index]

// 	index = index + 1
// 	r2 := q[index]

// 	question.Record = uint16(r1, r2)

// 	// get class
// 	index = index + 1
// 	c1 := q[index]

// 	index = index + 1
// 	c2 := q[index]

// 	question.Class = uint16(c1, c2)

// 	return question, nil
// }

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
