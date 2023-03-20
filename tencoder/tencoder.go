package tencoder

import (
	"io"
	"crypto/sha1"
	"bytes"
	
	"github.com/jackpal/bencode-go"
)

type bencodeFile struct {
	Path []string `bencode:"path"`
	Length int  `bencode:"length"`
}

type bencodeInfo struct {
    Pieces      string `bencode:"pieces"`
    PieceLength int    `bencode:"piece length"`
    Length      int    `bencode:"length"`
	// TODO: add multifile torrent support
	// Files       []bencodeFile `bencode:"files"`
    Name        string `bencode:"name"`
}

type bencodeTorrent struct {
    Announce string      `bencode:"announce"`
    Info     bencodeInfo `bencode:"info"`
}

type TorrentSubFile struct {
	Path []string
	Length int
}

type TorrentFile struct {
    Announce    string
    InfoHash    [20]byte
    PieceHashes [][20]byte
    PieceLength int
    Length      int
    Name        string
	// Files       []TorrentSubFile
} 

func Unmarshal(r io.Reader) (*TorrentFile, error) {
	bto := bencodeTorrent{}
    err := bencode.Unmarshal(r, &bto)
    if err != nil {
        return nil, err
    }

	t, err := bto.toTorrentFile()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (info *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *info)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (bto *bencodeTorrent) toTorrentFile() (*TorrentFile, error) {
    t := TorrentFile{}

	t.Announce = bto.Announce
	t.PieceLength = bto.Info.PieceLength
	t.Length = bto.Info.Length
	t.Name = bto.Info.Name
	info_hash, err := bto.Info.hash()
	if err != nil {
		return nil, err
	}
	t.InfoHash = info_hash
	for i := 0; i < len(bto.Info.Pieces); i += 20 {
		buf := bytes.Buffer{}
		buf.Write([]byte(bto.Info.Pieces[i:i+20]))
		t.PieceHashes = append(t.PieceHashes, [20]byte(buf.Bytes()))
	}
	// for _, v := range bto.Info.Files {
	// 	f := TorrentSubFile{}
	// 	f.Path = v.Path
	// 	f.Length = v.Length
	// 	t.Files = append(t.Files, f)
	// }
	
	return &t, nil
}
