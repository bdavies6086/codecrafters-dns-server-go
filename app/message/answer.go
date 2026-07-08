package message

type Answer struct {
	Labels   []string
	Record   uint16
	Class    uint16
	Ttl      uint32
	RDLength uint16
	RData    []byte
}

func (a Answer) Encode() []byte {
	bs := []byte{}

	for _, label := range a.Labels {
		bs = append(bs, uint8(len(label)))
		bs = append(bs, []byte(label)...)
	}

	// null byte after labels
	bs = append(bs, uint8(0))

	bs = append(bs, uint8(a.Record>>8))
	bs = append(bs, uint8(a.Record))

	bs = append(bs, uint8(a.Class>>8))
	bs = append(bs, uint8(a.Class))

	bs = append(bs, uint8(a.Ttl>>24))
	bs = append(bs, uint8(a.Ttl>>16))
	bs = append(bs, uint8(a.Ttl>>8))
	bs = append(bs, uint8(a.Ttl))

	bs = append(bs, uint8(a.RDLength>>8))
	bs = append(bs, uint8(a.RDLength))

	bs = append(bs, a.RData...)

	return bs
}
