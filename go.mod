module github.com/theakula/bittorrent

go 1.20

replace github.com/theakula/bittorrent/tencoder => ./tencoder
replace github.com/theakula/bittorrent/tracker => ./tracker
replace github.com/theakula/bittorrent/handshake => ./handshake
replace github.com/theakula/bittorrent/message => ./message
replace github.com/theakula/bittorrent/torrent => ./torrent
replace github.com/theakula/bittorrent/client => ./client
replace github.com/theakula/bittorrent/bitfield => ./bitfield

require (
	github.com/jackpal/bencode-go v1.0.0 // indirect
	github.com/theakula/bittorrent/tencoder v1.0.0
	github.com/theakula/bittorrent/tracker v1.0.0
	github.com/theakula/bittorrent/handshake v1.0.0
	github.com/theakula/bittorrent/message v1.0.0
	github.com/theakula/bittorrent/torrent v1.0.0
	github.com/theakula/bittorrent/client v1.0.0
	github.com/theakula/bittorrent/bitfield v1.0.0
)