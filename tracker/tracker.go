package tracker

import (
	"strconv"
	"net/url"
	"net"
	"net/http"
	"encoding/binary"
	
	"github.com/jackpal/bencode-go"

	"github.com/theakula/bittorrent/tencoder"
)

type PeersResponse struct {
	Interval int
	Peers string
}

type Peer struct {
	IP net.IP
	Port uint16
}

// TODO: add udp torrent support
// type UDPPeersRequest struct {
// 	ProtocolID int64 `bencode:"protocol_id"`
// 	Action int32 `bencode:"action"`
// 	TransactionID int32 `bencode:"transaction_id"`
// }

// type UDPPeersResponse struct {
// 	Action int32 `bencode:"action"`
// 	TransactionID int32 `bencode:"transaction_id"`
// 	ConnectionID int64 `bencode:"connection_id"`
// }

func GetPeers(t *tencoder.TorrentFile, p uint16, peerID [20]byte) ([]Peer, error)  {
	peers := make([]Peer, 0)
	base, err := url.Parse(t.Announce)
	if err != nil {
		return peers, err
	}

	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
        "peer_id":    []string{string(peerID[:])},
        "port":       []string{strconv.Itoa(int(p))},
        "uploaded":   []string{"0"},
        "downloaded": []string{"0"},
        "compact":    []string{"1"},
        "left":       []string{strconv.Itoa(t.Length)},
	}

	base.RawQuery = params.Encode()
	u := base.String()
	
	peersData := PeersResponse{}

	if base.Scheme == "https" || base.Scheme == "http" {
		resp, err := http.Get(u)
		if err != nil {
			return peers, err
		}
		// fmt.Println(t.Length)
		// fmt.Println(resp)
		err = bencode.Unmarshal(resp.Body, &peersData)
		// fmt.Println(peersData)
		if err != nil {
			return peers, err
		}
	} else if base.Scheme == "udp" {
		// TODO: get peers logic here (udp) 
		// // fmt.Println(base.Host)
		// s, err := net.ResolveUDPAddr("udp4", base.Host)
		// if err != nil {
		// 	return peers, err
		// }
		// // fmt.Println(s)
		// udpc, err := net.DialUDP("udp", nil, s)
		// if err != nil {
		// 	return peers, err
		// }
		// defer udpc.Close()

		// tr_id := rand.Int31()
		// req := UDPPeersRequest{
		// 	TransactionID: tr_id,
		// 	Action: 0,
		// 	ProtocolID: 0x41727101980,
		// }
		// // fmt.Println(req)
		// // breq := make([]byte, 16)
		// // binary.BigEndian.PutUint32(breq[0:8], uint32(req.ProtocolID))
		// // binary.BigEndian.PutUint16(breq[8:12], uint16(req.Action))
		// // binary.BigEndian.PutUint16(breq[12:16], uint16(req.TransactionID))
		// err = bencode.Marshal(udpc, req)
		// // fmt.Println(len(breq))
		// // fmt.Println(breq)
		// // udpc.Write(breq)
		// if err != nil {
		// 	return peers, err
		// }

		// ures := UDPPeersResponse{}
		// // buff := bufio.NewReader(udpc)
		// buff := make([]byte, 16)
		

		// n, err := udpc.Read(buff)
		// if err != nil {
		// 	return peers, err
		// }

		// if n != 16 {
		// 	// fmt.Println("Invalid response size")
		// 	return peers, nil
		// }

		// ures.Action = int32(binary.BigEndian.Uint32(buff[:4]))
		// ures.TransactionID = int32(binary.BigEndian.Uint32(buff[4:8]))
		// ures.ConnectionID = int64(binary.BigEndian.Uint64(buff[8:]))
		// // fmt.Println(buff)
		// // err = bencode.Unmarshal(buff, &ures)
		// // if err != nil {
		// // 	return peers, err
		// // }
		// // fmt.Println(ures)
	}
	
	for i := 0; i < len(peersData.Peers); i += 6 {
		p := Peer{}
		ip := peersData.Peers[i:i+4]
		p.IP = net.IPv4(ip[0], ip[1], ip[2], ip[3])
		p.Port = binary.BigEndian.Uint16([]byte(peersData.Peers[i+4:i+6]))
		peers = append(peers, p)
	}

	return peers, err
}

func (p *Peer) String() string {
	return p.IP.String() + ":" + strconv.Itoa(int(p.Port))
}
