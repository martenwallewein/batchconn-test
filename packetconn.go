package main

import "net"

func NewPacketConn() *PacketConn {
	pc := &PacketConn{}

	return pc
}

type PacketConn struct {
	local   *net.UDPAddr
	remote  *net.UDPAddr
	rawConn *net.UDPConn
}

func (pc *PacketConn) Listen(local string) error {
	addr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		return err
	}

	pc.local = addr
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	pc.rawConn = conn
	return nil
}

func (pc *PacketConn) Dial(local, remote string) error {
	addr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		return err
	}

	raddr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		return err
	}

	pc.local = addr
	pc.remote = raddr
	conn, err := net.DialUDP("udp", addr, raddr)
	if err != nil {
		return err
	}
	pc.rawConn = conn
	return nil
}

func (pc *PacketConn) Read(buf []byte) ([]byte, int, error) {
	count, err := pc.rawConn.Read(buf)
	return nil, count, err
}

func (pc *PacketConn) Write(buf []byte) (int, error) {
	return pc.rawConn.Write(buf)
}
