package xdrParse

import (
	"errors"
	//"preprocess/modules/mlog"
)

type DpiXdr struct {
	SessionStatus uint8
	AppId         uint8
	Tuple         FiveTupleInfo
	SesionStat    StandStat
	SesionTime    SessionTime
	BussiStat     StandStat
	TcpInfo       TcpInfo
	ResponseDelay uint32
	Level7Proto   Level7Proto
	HttpBaseInfo  HttpBaseInfo
	HttpInfo      HttpInfo //unfold
	Sip           SipInfo
	SipCaller     string
	SipCalled     string
	SipSessionId  string
	RTSP          RtspBase
	RtspUrl       string
	RtspUserAgent string
	RtspServerIp  string
	FtpInfo       FtpInfo  //unfold
	MailInfo      MailInfo //unfold
	DnsInfo       DnsInfo  //unfold
	//55-57
	//58-61
	HttpReqInfo  []byte
	HttpRespInfo []byte
	FileContent  []byte
}

type DnsInfo struct {
	DnsZones          string
	DnsIpNum          uint8
	DnsIpv4           string //UINT32??
	DnsIpv6           string //STRUCT??
	DnsRspCode        uint8
	DnsRsqCnt         uint8
	DnsRspRecordCnt   uint8
	DnsAuthCnttCnt    uint8
	DnsExtraRecordCnt uint8
	DnsRspDelay       uint32
}

type FtpInfo struct {
	FtpStatus    uint16
	FtpUserName  string
	FtpCurDir    string
	FtpTransMode uint8
	FtpTransType uint8
	FtpFileName  string
	FtpFileSize  uint32
	FtpRspTm     uint64
	FtpTransTime uint64
}

type MailInfo struct {
	MailMsgType    uint16
	MailRspStatus  uint16 //??
	MailUserName   string
	MailResvInfo   string
	MailLength     uint32
	MailDomainInfo string
	MailRcevAcnt   string
	MailHdr        string
	MailAcsType    uint8
}

type RtspBase struct {
	UsCStartPort    uint16
	UsCEndPort      uint16
	UsSStartPort    uint16
	UsSEndPort      uint16
	UsSsnVideoCount uint16
	UsSsnAudioCount uint16
	UlResDelay      uint32
}

type SipInfo struct {
	UcCallDirection uint8
	UcCallType      uint8
	UcHookReason    uint8
	UcSignalType    uint8
	UsDataflowNum   uint16
	UbSipIVMR       uint16 //ubSipInvite+ubSipBye+malloc+resv
}
type HttpInfo struct {
	HttpHost       string
	HttpUrl        string
	HttpOnlineHost string
	HttpUserAgent  string
	HttpContent    string
	HttpRefer      string
	HttpCookie     string
	HttpLocation   string
}
type HttpBaseInfo struct {
	Ulactiontime        uint64
	Ulfirst_packet_time uint64
	UlLast_Packet_Time  uint64
	UlServiceTime       uint64
	UlContentLenth      uint32
	UsHttpStatus        uint16
	UcHttpMethod        uint8
	Uchttpversion       uint8
	UcUnionFlag         uint8 //ucFirstRequestFlag+ucSerFlag+ucHeadFlag+res
	UcIE                uint8
	UcPortal            uint8
	Resv                uint8
}
type Level7Proto struct {
	AppStatus uint8
	ClassID   uint8
	Protocol  uint16
}

type TcpInfo struct {
	SynackToSynTime    uint16
	AckToSynTime       uint16
	UbReportFlag       uint8
	CloseSsnReason     uint8
	Reserve            uint16
	FirstRequestDelay  uint32
	FirstResponseDely  uint32
	TcpWindow          uint32
	Mss                uint16
	TcpRetryCount      uint8
	TcpRetryAckCount   uint8
	TcpAckCount        uint8
	TcpConnectStatusNo uint8
	TcpStatusFirst     uint8
	TcpStatusSecond    uint8
}

type SessionTime struct {
	CreatTime uint64
	StartTime uint64
	EndTIme   uint64
}

type StandStat struct {
	UpFlow            uint32
	DownFlow          uint32
	UpPkgNum          uint32
	DownPkgNum        uint32
	UpTcpOutPkgNum    uint32
	DownTcpOutPkgNum  uint32
	UpTcpRetrPkgNum   uint32
	DownTcpRetrPkgNum uint32
	UpIpBurstNum      uint32
	DownIpBurstNum    uint32
	TcpUdpStartTime   uint16
	TcpUdpEndTime     uint16
}
type FiveTupleInfo struct {
	Version  uint8
	Dir      uint8
	L4Proto  uint8
	Resv     uint8
	SrcPort  uint16
	DstPort  uint16
	SrcIpv4  uint32
	SrcIpv6v Ipv6_s
	DstIpv4  uint32
	DstIpv6v Ipv6_s
}
type Ipv6_s struct {
	A uint32
	B uint32
	C uint32
}

func (this *DpiXdr) CheckType() int {
	/*
		if len(this.FileContent) != 0 {
			return XdrFileType
		} else if len(this.HttpReqInfo) != 0 {
			return XdrHttpType
		} else {
			return XdrType
		}
	*/
	if len(this.HttpReqInfo) != 0 {
		return XdrHttpType
	} else if len(this.FileContent) != 0 {
		return XdrFileType
	} else {
		return XdrType
	}
	return -1
}

var ErrXdrHeadErr error
var ErrXdrNotEnoughLenErr error

func errInit() {
	ErrXdrHeadErr = errors.New("head check error")
	ErrXdrNotEnoughLenErr = errors.New("xdr data isn't enough long")
}
