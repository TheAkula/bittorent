package main

import (
	"log"
	"flag"
	"os"
	
	"github.com/theakula/bittorrent/tencoder"
	"github.com/theakula/bittorrent/torrent"
)

var tfile string
var port int
var out string

func main () {
	flag.Parse()

	if len(tfile) == 0 || port == 0 || len(out) == 0 {
		log.Fatalln("Usage: bittorent --file <torrent_file> --port <port> --out <out_file>")
	}
	
	f, err := os.Open(tfile)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	
	t, err := tencoder.Unmarshal(f)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(t.Length)
	trnt := torrent.New(t, uint16(port))
	buff, err := trnt.Run()
	if err != nil {
		log.Fatalln(err)
	}

	o, err := os.Open(out)
	if err != nil {
		log.Fatalln(err)
	}
	defer o.Close()

	_, err = o.Write(buff)
	if err != nil {
		log.Fatalln(err)
	}
}

func init () {
	flag.StringVar(&tfile, "file", "", "torrent file to download")
	flag.IntVar(&port, "port", 6881, "port on which bittorren will work")
	flag.StringVar(&out, "out", "", "out file")
}

