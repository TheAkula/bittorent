package torrent

import (
	"crypto/sha1"
	"fmt"
	"runtime"
	"time"
	
	"github.com/theakula/bittorrent/tencoder"
	"github.com/theakula/bittorrent/tracker"
	"github.com/theakula/bittorrent/client"
	"github.com/theakula/bittorrent/message"
)

type piece struct {
	index int
	buff []byte
	downloaded int
	client *client.Client
	requested int
	backlog int
}

type pieceWork struct {
	index int
	hash [20]byte
	length int
}

type pieceResult struct {
	index int
	buff []byte
}

type Torrent struct {
	Port uint16
	TorrentFile *tencoder.TorrentFile
	PeerID [20]byte
}

var PieceMaxSize int = 16384
var MaxRequestsBacklog int = 5

func New(t *tencoder.TorrentFile, port uint16) *Torrent {
	return &Torrent{
		Port: port,
		TorrentFile: t,
		PeerID: [20]byte(sha1.Sum(make([]byte, 0))),
	}
}

func (t *Torrent) Run() ([]byte, error) {
	peers, err := tracker.GetPeers(t.TorrentFile, t.Port, t.PeerID)
	if err != nil {
		return []byte{}, err
	}

	if len(peers) == 0 {
		return []byte{}, fmt.Errorf("No peers found")
	}

	results := make(chan *pieceResult)
	works := make(chan *pieceWork, len(t.TorrentFile.PieceHashes))
	
	for i, hash := range t.TorrentFile.PieceHashes {
		works <- &pieceWork{
			index: i,
			length: t.TorrentFile.PieceLength,
			hash: hash,
		}
	}
	
	for _, peer := range peers {
		go t.startDownloadPiece(peer, results, works)
	}
	res := make([]byte, t.TorrentFile.Length)

	done := 0
	for done < len(t.TorrentFile.PieceHashes) {
		r := <- results
		copy(res[r.index*t.TorrentFile.PieceLength:], r.buff)
		done++

		peers_num := runtime.NumGoroutine() - 1
		percent := float64(done) / float64(len(t.TorrentFile.PieceHashes)) * 100
		fmt.Printf("(%0.2f%%) Downloaded piece #%d, %d peers\n", percent, r.index, peers_num)
	}
	close(results)

	return res, nil
}

func (t *Torrent) startDownloadPiece(peer tracker.Peer,
	rs chan *pieceResult, ws chan *pieceWork) {
	c, err := client.New(peer, t.TorrentFile, t.PeerID)
	if err != nil {
		fmt.Println("Client create error: ", err)
		return 
	}
	defer c.Conn.Close()

	c.SendUnchoke()
	c.SendInterested()
	
	for pw := range ws {
		if !c.Bitfield.HasPiece(pw.index) {
			ws <- pw
			continue
		}

		pr, err := t.downloadPiece(pw, c)
		if err != nil {
			ws <- pw
			fmt.Println(err)
			return
		}
		
		rs <- pr
		err = c.SendHave(pw.index)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (t *Torrent) downloadPiece(pw *pieceWork, c *client.Client) (*pieceResult, error) {
	c.Conn.SetDeadline(time.Now().Add(30*time.Second))
	defer c.Conn.SetDeadline(time.Time{})

	p := piece{
		client: c,
		requested: 0,
		backlog: 0,
		index: pw.index,
		downloaded: 0,
		buff: make([]byte, t.TorrentFile.PieceLength),
	}
	
	for p.downloaded < pw.length {
		if !p.client.Choked && p.requested < t.TorrentFile.PieceLength {
			if p.backlog < MaxRequestsBacklog {
				length := PieceMaxSize
				if t.TorrentFile.PieceLength - p.requested < PieceMaxSize {
					length = t.TorrentFile.PieceLength - p.requested
				} 
				c.SendRequest(p.index, p.requested, length)
				p.backlog++
				p.requested += length
			}
		}
		err := p.readMessage()
		if err != nil {
			return nil, err
		}
	}

	pr := pieceResult{
		index: p.index,
		buff: p.buff,
	}

	return &pr, nil
}

func (p *piece) readMessage() error {
	msg, err := message.Read(p.client.Conn)
	if err != nil {
		return err
	}
	
	switch msg.ID {
	case message.MsgChoke:
		p.client.Choked = true
	case message.MsgUnchoke:
		p.client.Choked = false
	case message.MsgHave:
		index, err := msg.ParseHave()
		if err != nil {
			return err
		}
		p.client.Bitfield.SetPiece(index)
	case message.MsgPiece:
		n, err := msg.ParsePiece(p.index, p.buff)
		if err != nil {
			return err
		}
		p.downloaded += n
		p.backlog--
	}

	return nil
}
