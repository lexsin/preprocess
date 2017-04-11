package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"net"
	"preprocess/modules/mlog"
)

func Ipv4IntToString(n uint32) string {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, n)
	b := buf.Bytes()
	return net.IPv4(b[0], b[1], b[2], b[3]).String()
}

func Ipv4StringToInt(ip string) uint32 {
	bs := net.ParseIP(ip).To4()
	var n uint32
	buf := new(bytes.Buffer)
	buf.Write(bs)
	binary.Read(buf, binary.LittleEndian, &n)
	return n
}

func IntToBool(u uint32) bool {
	if u == 0 {
		return false
	} else {
		return true
	}
}

func Md5Sum(data []byte) []byte {
	sum := md5.Sum(data)
	mlog.Debug("sum[:]=", sum[:])
	return sum[:]
}
