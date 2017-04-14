package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"io"
	"net"
	"os"
	"path"
	"preprocess/modules/mlog"
	"strings"
)

func CheckSuffix(FileName string, suffixs ...string) bool {
	suffix := path.Ext(FileName)
	for _, stand := range suffixs {
		if suffix == "."+stand {
			return true
		}
	}
	return false
	/*
		if suffix != "xdr" && suffix != "XDR" {
			return false
		}
		return true
	*/
}
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

func DealFilePerline(fileName string, handler func(string) error) (int, error) {
	var n = 0
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err := handler(line); err != nil {
			mlog.Error("line:", line, "handler error:", err.Error())
		}
		n++
		if err != nil {
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
	return n, nil
}
