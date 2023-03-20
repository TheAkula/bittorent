package bitfield

type Bitfield []byte

func (b Bitfield) HasPiece(index int) bool {
	offset := index % 8
	byteIndex := index / 8
	return b[byteIndex]>>(7-offset)&1 != 0
}

func (b Bitfield) SetPiece(index int) {
	offset := index % 8
	b[index / 8] |= 1 << (7 - offset)
}
