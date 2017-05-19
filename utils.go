package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"path"
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"
	"strings"
)

func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

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
	//mlog.Debug("sum[:]=", sum[:])
	return sum[:]
}

func DealFilePerline(fileName string, handlers perAlertFuncs) (int, error) {
	var n = 0
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
		line = strings.TrimSpace(line)
		if err := handlers.checkForm(line); err != nil {
			mlog.Error(fileName, "form error!")
			return n, err
		}
		if err := handlers.pushkafka(line); err != nil {
			mlog.Error("line:", line, "handler error:", err.Error())
		}
		n++
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func DeleteFile(fileName string) error {
	if err := os.Remove(fileName); err != nil {
		mlog.Error("remove file", fileName, "err:", err.Error())
		return err
	}
	mlog.Debug("remove file", fileName, "success")
	return nil
}

/**
 * oldFile is dir+file, newPath is dir only
 */
func RenameFile(oldFile string, newPath string) error {
	file := path.Base(oldFile)
	return os.Rename(oldFile, newPath+"/"+file)
}

func GetConfBool(result *bool, section string, option string, def bool) error {
	str, err := mconfig.Conf.String(section, option)
	if err != nil {
		*result = def
		return errors.New(section + option + " not configure")
	}
	switch str {
	case "true":
		*result = true
	case "false":
		*result = false
	default:
		*result = def
	}
	return nil
}

func StringToBool(tof string) (bool, error) {
	if tof == "true" {
		return true, nil
	} else if tof == "false" {
		return false, nil
	} else {
		return false, errors.New("parament error")
	}
}

func CreateDir(dir string) error {
	if absolute := path.IsAbs(dir); !absolute {
		mlog.Error(dir, "is not absolute")
		return nil
	}
	_, err := os.Stat(dir)
	if err == nil {
		return nil
	} else {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0777); err != nil {
				mlog.Error("create path ", dir, "err:", err.Error())
				return err
			}
		}
	}
	return err
}
