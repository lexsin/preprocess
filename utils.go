package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"net"
)

func IntToIpv4(n uint32) string {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, n)
	b := buf.Bytes()
	return net.IPv4(b[3], b[2], b[1], b[0]).String()
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
