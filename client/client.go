package client

import (
	"net"
	"time"
	"fmt"
	
	"github.com/theakula/bittorrent/tracker"
	"github.com/theakula/bittorrent/bitfield"
	"github.com/theakula/bittorrent/tencoder"
	"github.com/theakula/bittorrent/handshake"
	"github.com/theakula/bittorrent/message"
)

type Client struct {
	Peer tracker.Peer
	Conn net.Conn
	Choked bool
	Interested bool
	Bitfield bitfield.Bitfield
}

func receiveBitfield(conn net.Conn) ([]byte, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})
	
	msg, err := message.Read(conn)
	if err != nil {
		return []byte{}, err
	}
	if msg == nil {
		return nil, fmt.Errorf("Keep alive message")
	}
	if msg.ID != message.MsgBitfield {
		return []byte{}, fmt.Errorf("Expected msg with id \"bitfield\" but got: ", msg)
	}

	return msg.Payload, nil
}

func New(p tracker.Peer, t *tencoder.TorrentFile, peerID [20]byte) (*Client, error) {
	con, err := net.DialTimeout("tcp", p.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	req_hndshk := handshake.New(t.InfoHash, peerID)
	con.Write(req_hndshk.Serialize())

	res_hndshk, err := handshake.Read(con)
	if err != nil {
		return nil, err
	}

	if res_hndshk == nil {
		return nil, fmt.Errorf("No response handshake")
	}

	if res_hndshk.InfoHash != req_hndshk.InfoHash {
		return nil, fmt.Errorf("Handshake failed")
	}

	bf, err := receiveBitfield(con)
	if err != nil {
		return nil, err
	}
	
	return &Client{
		Peer: p,
		Conn: con,
		Choked: true,
		Interested: false,
		Bitfield: bf,
	}, nil
}

func (c *Client) SendInterested() error {
	msg := message.Message{
		ID: message.MsgInterested,
	}
	_, err := c.Conn.Write(msg.Serialize())
	
	if err != nil {
		return err
	}
	c.Interested = true
	return nil
}

func (c *Client) SendUnchoke() error {
	msg := message.Message{
		ID: message.MsgUnchoke,
	}
	_, err := c.Conn.Write(msg.Serialize())
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendRequest(index int, begin int, length int) error {
	msg := message.FormatRequest(index, begin, length)
	_, err := c.Conn.Write(msg.Serialize())
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendHave(index int) error {
	msg := message.FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	if err != nil {
		return err
	}
	return nil
}
