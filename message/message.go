package message

import (
	"io"
	"fmt"
	"encoding/binary"
)

type messageID uint8

const (
    MsgChoke         messageID = 0
    MsgUnchoke       messageID = 1
    MsgInterested    messageID = 2
    MsgNotInterested messageID = 3
    MsgHave          messageID = 4
    MsgBitfield      messageID = 5
    MsgRequest       messageID = 6
    MsgPiece         messageID = 7
    MsgCancel        messageID = 8
)
	
type Message struct {
	ID messageID
	Payload []byte
}

func Read(r io.Reader) (*Message, error) {
	size_buff := make([]byte, 4)

	n, err := io.ReadFull(r, size_buff)
	if err != nil {
		return nil, err
	}

	if n != 4 {
		return nil, fmt.Errorf("Invalid message size length: ", n)
	}

	size := binary.BigEndian.Uint32(size_buff)
	if size == 0 {
		return nil, fmt.Errorf("Keep alive message")
	}

	buff := make([]byte, size)
	n, err = io.ReadFull(r, buff)

	if err != nil {
		return nil, err
	}
	return &Message{
		ID: messageID(buff[0]),
		Payload: buff[1:],
	}, nil
}

func (m *Message) Serialize() []byte {
	buff := make([]byte, 4 + 1 + len(m.Payload))
	binary.BigEndian.PutUint32(buff[:4], uint32(len(m.Payload)+1))
	buff[4] = byte(m.ID)
	copy(buff[5:], m.Payload)
	return buff
}

func (m *Message) ParseHave() (int, error) {
	if m.ID != MsgHave {
		return 0, fmt.Errorf("Message id is not equal to \"have\"")
	}

	if len(m.Payload) != 4 {
		return 0, fmt.Errorf("Message payload size is not equal to 4")
	}

	return int(binary.BigEndian.Uint32(m.Payload)), nil
}

func (m *Message) ParsePiece(index int, buffer []byte) (int, error) {
	if m.ID != MsgPiece {
		return 0, fmt.Errorf("Message id is not equal to \"piece\"")
	}
	parsedIndex := int(binary.BigEndian.Uint32(m.Payload[:4]))
	if parsedIndex != index {
		return 0, fmt.Errorf("Piece index is not the same")
	}
	begin := int(binary.BigEndian.Uint32(m.Payload[4:8]))
	if begin > len(buffer) {
		return 0, fmt.Errorf("Bad message begin index")
	}
	data := m.Payload[8:]

	copy(buffer[begin:], data)

	return len(data), nil
}

func FormatRequest(index int, begin int, length int) *Message {
	p := make([]byte, 12)
	binary.BigEndian.PutUint32(p[0:4], uint32(index))
	binary.BigEndian.PutUint32(p[4:8], uint32(begin))
	binary.BigEndian.PutUint32(p[8:12], uint32(length))
	return &Message{
		ID: MsgRequest,
		Payload: p,
	}
}

func FormatHave(index int) *Message {
	p := make([]byte, 4)
	binary.BigEndian.PutUint32(p, uint32(index))
	return &Message{
		ID: MsgHave,
		Payload: p,
	}
}
