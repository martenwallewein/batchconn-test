package main

import (
	"net"

	"golang.org/x/net/ipv4"
)

type batchConn interface {
	WriteBatch(ms []ipv4.Message, flags int) (int, error)
	ReadBatch(ms []ipv4.Message, flags int) (int, error)
}

type BatchConn struct {
	local   *net.UDPAddr
	remote  *net.UDPAddr
	rawConn net.PacketConn
	conn    batchConn
	msgs    []ipv4.Message
	tx      []ipv4.Message
}

func NewBatchConn() *BatchConn {
	msgs := make([]ipv4.Message, batchSize)
	tx := make([]ipv4.Message, batchSize)
	for k := range msgs {
		msgs[k].Buffers = [][]byte{make([]byte, 1500)}
	}

	bc := &BatchConn{
		msgs: msgs,
		tx:   tx,
	}

	return bc
}

func (bc *BatchConn) Listen(local string) error {
	addr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		return err
	}

	bc.local = addr
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	bc.rawConn = conn
	bc.transformConn()
	return nil
}

func (bc *BatchConn) Dial(local, remote string) error {
	addr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		return err
	}

	raddr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		return err
	}

	bc.local = addr
	bc.remote = raddr
	conn, err := net.DialUDP("udp", addr, raddr)
	if err != nil {
		return err
	}
	bc.rawConn = conn
	bc.transformConn()
	return nil
}

func (bc *BatchConn) transformConn() {
	var xconn batchConn
	xconn = ipv4.NewPacketConn(bc.rawConn)
	bc.conn = xconn
}

func (bc *BatchConn) Read() ([]byte, int, error) {
	count, err := bc.conn.ReadBatch(bc.msgs, 0)
	return nil, count, err
}

func (bc *BatchConn) Write(buf []byte) (int, error) {
	var msg ipv4.Message
	for i := 0; i < batchSize; i++ {
		msg.Buffers = [][]byte{buf}
		msg.Addr = bc.remote
		bc.tx[i] = msg
	}
	count, err := bc.conn.WriteBatch(bc.tx, 0)
	// bc.tx = make([]ipv4.Message, batchSize)
	return count, err
}
