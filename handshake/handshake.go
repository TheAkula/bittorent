package handshake

import (
	"io"
	"fmt"
)

type Handshake struct {
	Pstr string
	InfoHash [20]byte
	PeerID [20]byte
}

func New(infoHash [20]byte, peerID [20]byte) *Handshake {
	return &Handshake{
		InfoHash: infoHash,
		PeerID: peerID,
		Pstr: "BitTorrent protocol",
	}
}

func (h *Handshake) Serialize() []byte {
	buff := make([]byte, len(h.Pstr)+49)
	buff[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buff[curr:], h.Pstr)
	curr += copy(buff[curr:], make([]byte, 8))
	curr += copy(buff[curr:], h.InfoHash[:])
	curr += copy(buff[curr:], h.PeerID[:])
	return buff
}

func Read(r io.Reader) (*Handshake, error) {
	h := Handshake{}

	hlen := 1+19+8+20+20
	
	buff := make([]byte, hlen)
	n, err := r.Read(buff)
	if err != nil {
		return nil, err
	}

	if n != hlen {
		return nil, nil
	}

	if int(buff[0]) > len(buff) {
		return nil, fmt.Errorf("Invalid handshake size: ", int(buff[0]))
	}
	h.Pstr = string(buff[1:int(buff[0]) + 1])
	h.InfoHash = [20]byte(buff[1+19+8:1+19+8+20])
	h.PeerID = [20]byte(buff[1+19+8+20:])
	return &h, nil
}
