package main

import (
	"log"

	"github.com/anacrolix/tagflag"
)

const (
	batchSize = 32
)

var flags = struct {
	IsServer bool
	Remote   string // 19-ffaa:1:c3f,[10.0.0.2]
	Local    string
	Type     string
	tagflag.StartPos
}{
	IsServer: true,
	Type:     "batchconn",
}

func main() {
	tagflag.Parse(&flags)

	if flags.IsServer {
		runServer()
	} else {
		runClient()
	}
}

func runServer() {
	if flags.Type == "batchconn" {
		b := NewBatchConn()
		err := b.Listen(flags.Local)
		if err != nil {
			log.Fatal(err)
		}

		for {
			b.Read()
		}
	} else {
		b := NewPacketConn()
		err := b.Listen(flags.Local)
		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, 1400)
		for {
			b.Read(buf)
		}
	}

}

func runClient() {
	if flags.Type == "batchconn" {
		b := NewBatchConn()
		err := b.Dial(flags.Local, flags.Remote)
		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, 1400)
		for {
			b.Write(buf)
		}
	} else {
		b := NewPacketConn()
		err := b.Dial(flags.Local, flags.Remote)
		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, 1400)
		for {
			b.Write(buf)
		}
	}
}
