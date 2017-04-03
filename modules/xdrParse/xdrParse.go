package xdrParse

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"
)

type TlvFormExtend struct {
	baseInfo struct {
		TlvId      uint8
		ShortData  uint8
		TypeIdhigh uint8
		Reserve    uint8
		TlvLength  uint32
	}
	TlvData []byte
}

type TlvForm struct {
	baseInfo struct {
		TlvId      uint8
		ShortData  uint8
		TypeIdHigh uint8
		Reserve    uint8
		TlvLength  uint32
	}
	TlvData []byte
}

type TlvValue struct {
	TlvId     uint8
	ShortData uint8
	IsExtend  bool
	DataLen   uint32
	Data      []byte
}

const (
	XdrType = iota
	XdrHttpType
	XdrFileType
)

func parseOneXdr(xdr *TlvValue, obj *DpiXdr) error {
	return DecodeFuncMap[int(xdr.TlvId)](xdr, obj)
}

func rangeParseXdr(xdrs []TlvValue) *DpiXdr {
	var xdrOjbect DpiXdr
	for _, xdr := range xdrs {
		if err := parseOneXdr(&xdr, &xdrOjbect); err != nil {
			mlog.Error(fmt.Sprintf("parse xdr %s error:%s", xdr.TlvId, err.Error()))
		}
	}
	return &xdrOjbect
}

func ParseXdr(origiData []byte) ([]*DpiXdr, error) {
	var results []*DpiXdr
	tlvValues := RangeToObj(origiData)
	for _, tlv := range tlvValues {
		xdrs := RangeToObj(tlv.Data)
		obj := rangeParseXdr(xdrs)
		results = append(results, obj)
	}
	return results, nil
}

func RangeToObj(data []byte) []TlvValue {
	var list []TlvValue
	buf := new(bytes.Buffer)
	offset := 0
	for {
		if IsExtend(data) {
			var headExt TlvFormExtend
			buf.Reset()
			headsize := binary.Size(headExt.baseInfo)
			buf.Write(data[offset : offset+headsize])
			binary.Read(buf, binary.LittleEndian, &headExt.baseInfo)
			value := TlvValue{
				TlvId:     headExt.baseInfo.TlvId,
				ShortData: headExt.baseInfo.ShortData,
				IsExtend:  true,
				DataLen:   uint32(headExt.baseInfo.TlvLength - uint32(headsize)),
				Data:      data[headsize:headExt.baseInfo.TlvLength],
			}
			list = append(list, value)
			offset += int(headExt.baseInfo.TlvLength)
		} else {
			var head TlvFormExtend
			buf.Reset()
			headsize := binary.Size(head.baseInfo)
			buf.Write(data[offset : offset+headsize])
			binary.Read(buf, binary.LittleEndian, &head.baseInfo)
			value := TlvValue{
				TlvId:     head.baseInfo.TlvId,
				ShortData: head.baseInfo.ShortData,
				IsExtend:  false,
				DataLen:   uint32(head.baseInfo.TlvLength - uint32(headsize)),
				Data:      data[headsize:head.baseInfo.TlvLength],
			}
			list = append(list, value)
			offset += int(head.baseInfo.TlvLength)
		}
	}
}

/*
func GetIdAndData(data []byte) (int, []byte) {
	buf := new(bytes.Buffer)
	if IsExtend(data) {
		var headExt TlvFormExtend
		buf.Reset()
		buf.Write(data[:binary.Size(headExt.baseInfo)])
		binary.Read(buf, binary.LittleEndian, &headExt.baseInfo)
	} else {
		var head TlvForm
		buf.Reset()
		buf.Write(data[:binary.Size(head.baseInfo)])
		binary.Read(buf, binary.LittleEndian, &head.baseInfo)
	}
}
*/

func IsExtend(data []byte) bool {
	buf := new(bytes.Buffer)
	buf.Write(data[:4])
	var n int32
	binary.Read(buf, binary.LittleEndian, &n)
	if n&0x00008000 == 0 {
		return false
	}
	return true
}

func GetTlvLength(data []byte) (result int) {
	buf := new(bytes.Buffer)
	if !IsExtend(data) {
		var head TlvForm
		buf.Reset()
		buf.Write(data[:binary.Size(head.baseInfo)])
		binary.Read(buf, binary.LittleEndian, &head.baseInfo)
		result = int(head.baseInfo.TlvLength)
	} else {
		var headExt TlvFormExtend
		buf.Reset()
		buf.Write(data[:binary.Size(headExt.baseInfo)])
		binary.Read(buf, binary.LittleEndian, &headExt.baseInfo)
		result = int(headExt.baseInfo.TlvLength)
	}
	return
}

func ParseXdrHead(data []byte) ([][]byte, error) {
	var list [][]byte
	totallen := len(data)
	offset := 0
	for {
		tlvLen := GetTlvLength(data[offset:])
		list = append(list, data[offset:offset+tlvLen])
		offset += tlvLen
		if offset >= totallen {
			break
		}
	}
	return list, nil
}

func (this *DpiXdr) CheckType() int {
	//TODO
	return XdrType
}

func (this *DpiXdr) HashPartation() int32 {
	//init topic partition
	var err error
	_, err = mconfig.Conf.Int("kafka", "AgentNum")
	if err != nil {
		mlog.Error("app.conf AgentNum error")
		panic("app.conf AgentNum error")
	}
	//TODO
	return 0
}
