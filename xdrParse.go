package main

import (
	"bytes"
	"encoding/binary"
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"
)

func parseElements(data []byte) (*DpiXdr, error) {

}

func GetTlvId(data []byte) int32 {
	baseSize := binary.Size(head.baseInfo)
}

func parseEachXdr(xdr *TlvValue, obj *DpiXdr) error {
	switch xdr.TlvId {
	case 0:
		mlog.Debug("find xdr id=0")
	case XDR_SESSION_STATUS:
		if err := ParsSessionStatus(xdr, obj); err != nil {
			return err
		}
	case XDR_APP_ID:
		if err := parsAppId(); err != nil {
			return err
		}

	}
}

func rangeParseXdr(xdrs []TlvValue) *DpiXdr {
	var xdrOjbect DpiXdr
	for _, xdr := range xdrs {
		parseEachXdr(&xdr, &DpiXdr)
	}
}

func ParseXdr(origiData []byte) ([]*DpiXdr, error) {
	/*
		list, _ := ParseXdrHead(origiData)
		results := make([]*DpiXdr, 0)
		for _, data := range list {
			xdr := parseElements(data)
		}
		results = append(results, xdr)
	*/
	tlvValues := RangeToObj(origiData)
	for _, tlv := range tlvValues {
		xdrs := RangeToObj(tlv)
		obj := rangeParseXdr(xdrs)
	}
	return results, nil
}

type TlvFormExtend struct {
	baseInfo struct {
		TlvId      uint8
		ShortData  uint8
		TypeIdhigh uint8
		Reserve    uint8
		DataLen    uint32
	}
	TlvData []byte
}
type TlvForm struct {
	baseInfo struct {
		TlvId      uint8
		ShortData  uint8
		TypeIdHigh uint8
		Reserve    uint8
		DataLen    uint32
	}
	TlvData []byte
}

type TlvValue struct {
	TlvId     uint8
	ShortData uint8
	IsExtend  bool
	TlvLen    uint32
	Data      []byte
}

func RangeToObj(data []byte) []TlvValue {
	var list []TlvValue
	buf := new(bytes.Buffer)
	totallen := len(data)
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
				DataLen:   headExt.baseInfo.TlvLength - headsize,
				Data:      data[headsize:headExt.baseInfo.TlvLength],
			}
			list = append(list, value)
			offset += headExt.baseInfo.TlvLength
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
				TlvLen:    head.baseInfo.TlvLength - headsize,
				Data:      data[headsize:head.baseInfo.TlvLength],
			}
			list = append(list, value)
			offset += headExt.baseInfo.TlvLength
		}
	}
}

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

func IsExtend(data []byte) bool {
	headSize := 0
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
		result = head.baseInfo.TlvLength
	} else {
		var headExt TlvFormExtend
		buf.Reset()
		buf.Write(data[:binary.Size(headExt.baseInfo)])
		binary.Read(buf, binary.LittleEndian, &headExt.baseInfo)
		result = headExt.baseInfo.TlvLength
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

const (
	XdrType = iota
	XdrHttpType
	XdrFileType
)

func (this *DpiXdr) CheckType() int {
	//TODO
	return XdrType
}

func (this *DpiXdr) HashPartation() int32 {
	//init topic partition
	var err error
	AgentNum, err = mconfig.Conf.Int("kafka", "AgentNum")
	if err != nil {
		mlog.Error("app.conf AgentNum error")
		panic("app.conf AgentNum error")
	}
	//TODO
	return 0
}
