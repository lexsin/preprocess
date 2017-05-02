package main

import (
	"bytes"
	"encoding/binary"
	//"preprocess/modules/mlog"
	"preprocess/modules/xdrParse"
)

type BackendInfo struct {
	Type int
	Data *BackendObj
}

/*
 * object array slice 不能为空
 */
type BackendObj struct {
	Vendor string `json:Verdor`
	Id     uint64 `json:id`
	Ipv4   bool   `json:Ipv4`
	Class  uint8  `json:Class`
	Type   uint32 `json:Type`
	Time   uint64 `json:Time`
	Conn   struct {
		Proto uint8  `json:Proto`
		Sport uint16 `json:Sport`
		Dport uint16 `json:Dport`
		Sip   string `json:Sip`
		Dip   string `json:Dip`
	} `json:Conn`
	ConnEx struct {
		Over bool `json:Over`
		Dir  bool `json:Dir`
	} `json:ConnEx`
	ConnSt struct {
		FlowUp     uint64 `json:FlowUp`
		FlowDown   uint64 `json:FlowDown`
		PktUp      uint64 `json:PktUp`
		PktDown    uint64 `json:PktDown`
		IpFragUp   uint64 `json:IpFragUp`
		IpFragDown uint64 `json:IpFragDown`
	} `json:ConnEt`
	ConnTime struct {
		Start uint64 `json:Start`
		End   uint64 `json:End`
	} `json:ConnTime`
	ServSt struct {
		FlowUp          uint64 `json:FlowUp`
		FlowDown        uint64 `json:FlowDown`
		PktUp           uint64 `json:PktUp`
		PktDown         uint64 `json:PktDown`
		IpFragUp        uint64 `json:IpFragUp`
		IpFragDown      uint64 `json:IpFragDown`
		TcpDisorderUp   uint64 `json:TcpDisorderUp`
		TcpDisorderDown uint64 `json:TcpDisorderDown`
		TcpRetranUp     uint64 `json:TcpRetranUp`
		TcpRetranDown   uint64 `json:TcpRetranDown`
	} `json:ServSt`
	Tcp struct {
		DisorderUp        uint64 `json:DisorderUp`
		DisorderDown      uint64 `json:DisorderDown`
		RetranUp          uint64 `json:RetranUp`
		RetranDown        uint64 `json:RetranDown`
		SynAckDelay       uint16 `json:SynAckDelay`
		AckDelay          uint16 `json:AckDelay`
		ReportFlag        uint8  `json:ReportFlag`
		CloseReason       uint8  `json:CloseReason`
		FirstRequestDelay uint32 `json:FirstRequestDelay`
		FirstResponseDely uint32 `json:FirstResponseDely`
		Window            uint32 `json:Window`
		Mss               uint16 `json:Mss`
		SynCount          uint64 `json:SynCount`
		SynAckCount       uint64 `json:SynAckCount`
		AckCount          uint8  `json:AckCount`
		SessionOK         bool   `json:SessionOK`
		Handshake12       bool   `json:Handshake12`
		Handshake23       bool   `json:Handshake23`
		Open              int32  `json:Open`
		Close             int32  `json:Close`
	} `json:Tcp`
	Http struct {
		Host              string `json:Host`
		Url               string `json:Url`
		XonlineHost       string `json:XonlineHost`
		UserAgent         string `json:UserAgent`
		ContentType       string `json:ContentType`
		Refer             string `json:Refer`
		Cookie            string `json:Cookie`
		Location          string `json:Location`
		Request           string `json:Request`
		Response          string `json:Response`
		RequestTime       uint64 `json:RequestTime`
		FirstResponseTime uint64 `json:FirstResponseTime`
		LastContentTime   uint64 `json:LastContentTime`
		ServTime          uint64 `json:ServTime`
		ContentLen        uint32 `json:ContentLen`
		StateCode         uint16 `json:StateCode`
		Method            uint8  `json:Method`
		Version           uint8  `json:Version`
		HeadFlag          bool   `json:HeadFlag`
		ServFlag          uint8  `json:ServFlag`
		RequestFlag       bool   `json:RequestFlag`
		Browser           uint8  `json:Browser`
		Portal            uint8  `json:Portal`
	} `json:Http`
	Sip struct {
		CallingNo    string `json:CallingNo`
		CalledNo     string `json:CalledNo`
		SessionId    string `json:SessionId`
		CallDir      uint8  `json:CallDir`
		CallType     uint8  `json:CallType`
		HangupReason uint8  `json:HangupReason`
		SignalType   uint8  `json:SignalType`
		StreamCount  uint16 `json:StreamCount`
		Malloc       bool   `json:Malloc`
		Bye          bool   `json:Bye`
		Invite       bool   `json:Invite`
	} `json:Sip`
	Rtsp struct {
		Url              string `json:Url`
		UserAgent        string `json:UserAgent`
		ServerIp         string `json:ServerIp`
		ClientBeginPort  uint16 `json:ClientBeginPort`
		ClientEndPort    uint16 `json:ClientEndPort`
		ServerBeginPort  uint16 `json:ServerBeginPort`
		ServerEndPort    uint16 `json:ServerEndPort`
		VideoStreamCount uint16 `json:VideoStreamCount`
		AudeoStreamCount uint16 `json:AudeoStreamCount`
		ResDelay         uint32 `json:ResDelay`
	} `json:Rtsp`
	Ftp struct {
		State      uint16 `json:State`
		UserCount  uint64 `json:UserCount`
		CurrentDir string `json:CurrentDir`
		TransMode  uint8  `json:TransMode`
		TransType  uint8  `json:TransType`
		FileCount  uint64 `json:FileCount`
		FileSize   uint32 `json:FileSize`
		RspTm      uint64 `json:RspTm`
		TransTm    uint64 `json:TransTm`
	} `json:Ftp`
	Mail struct {
		MsgType     uint16 `json:MsgType`
		RspState    uint16 `json:RspState`
		UserName    string `json:UserName`
		RecverInfo  string `json:RecverInfo`
		Len         uint32 `json:Len`
		DomainInfo  string `json:DomainInfo`
		RecvAccount string `json:RecvAccount`
		Hdr         string `json:Hdr`
		AcsType     uint8  `json:AcsType`
	} `json:Mail`
	Dns struct {
		Domain           string `json:Domain`
		IpCount          uint8  `json:IpCount`
		Ipv4             string `json:Ipv4`
		Ipv6             string `json:Ipv6`
		RspCode          uint8  `json:RspCode`
		ReqCount         uint8  `json:ReqCount`
		RspRecordCount   uint8  `json:RspRecordCount`
		AuthCnttCount    uint8  `json:AuthCnttCount`
		ExtraRecordCount uint8  `json:ExtraRecordCount`
		RspDelay         uint32 `json:RspDelay`
		PktValid         bool   `json:PktValid`
	} `json:Dns`
	Vpn struct {
		Type uint64 `json:Type`
	} `json:Vpn`
	Proxy struct {
		Type uint64 `json:Type`
	} `json:Proxy`
	QQ struct {
		Number string `json:Number`
	} `json:QQ`
	App struct {
		ProtoInfo uint64 `json:ProtoInfo`
		Status    uint64 `json:Status`
		ClassId   uint64 `json:ClassId`
		Proto     uint64 `json:Proto`
		File      string `json:File`
	} `json:App`
	/*
		Alert struct {
		} `json:"-"`
	*/
}

func PerTransToBackendObj(src *xdrParse.DpiXdr) *BackendInfo {
	obj := &BackendObj{}
	//obj.Vendor =
	//Id
	if src.Tuple.Version == 0 {
		obj.Ipv4 = false
	} else {
		obj.Ipv4 = true
	}
	obj.Class = src.AppId
	if len(src.HttpReqInfo) != 0 {
		obj.Type = XdrHttpType
	} else if len(src.FileContent) != 0 {
		obj.Type = XdrFileType
	} else {
		obj.Type = XdrType
	}
	obj.Time = src.SesionTime.StartTime
	obj.Conn.Proto = src.Tuple.L4Proto
	obj.Conn.Sport = src.Tuple.SrcPort
	obj.Conn.Dport = src.Tuple.DstPort
	obj.Conn.Sip = Ipv4IntToString(src.Tuple.SrcIpv4)
	obj.Conn.Dip = Ipv4IntToString(src.Tuple.DstIpv4)
	//obj.ConnEx.Over =
	//obj.ConnEx.Dir
	obj.ConnSt.FlowUp = uint64(src.SesionStat.UpFlow)
	obj.ConnSt.FlowDown = uint64(src.SesionStat.DownFlow)
	obj.ConnSt.PktUp = uint64(src.SesionStat.UpPkgNum)
	obj.ConnSt.PktDown = uint64(src.SesionStat.DownPkgNum)
	//obj.ConnSt.IpFragUp =
	//obj.ConnSt.IpFragDown
	obj.ConnTime.Start = src.SesionTime.StartTime
	obj.ConnTime.End = src.SesionTime.EndTIme
	//obj.ServSt.FlowUp
	//obj.ServSt.FlowDown
	//obj.ServSt.PktUp
	//obj.ServSt.PktDown
	//obj.ServSt.IpFragUp
	//obj.ServSt.IpFragDown
	//obj.ServSt.TcpDisorderUp
	//obj.ServSt.TcpDisorderDown
	//obj.ServSt.TcpRetranUp
	//obj.ServSt.TcpRetranDown

	//obj.Tcp.DisorderUp =
	//obj.Tcp.DisorderDown
	obj.Tcp.SynAckDelay = src.TcpInfo.SynackToSynTime
	obj.Tcp.AckDelay = src.TcpInfo.AckToSynTime
	obj.Tcp.ReportFlag = src.TcpInfo.UbReportFlag
	obj.Tcp.CloseReason = src.TcpInfo.CloseSsnReason
	obj.Tcp.FirstRequestDelay = src.TcpInfo.FirstRequestDelay
	obj.Tcp.FirstResponseDely = src.TcpInfo.FirstResponseDely
	obj.Tcp.Window = src.TcpInfo.TcpWindow
	obj.Tcp.Mss = src.TcpInfo.Mss
	//obj.Tcp.SynCount = src.TcpInfo.
	//obj.Tcp.SynAckCount
	obj.Tcp.AckCount = src.TcpInfo.TcpAckCount
	//obj.Tcp.SessionOK =
	if src.TcpInfo.TcpStatusFirst == 0 {
		obj.Tcp.Handshake12 = true
	} else {
		obj.Tcp.Handshake12 = false
	}
	if src.TcpInfo.TcpStatusSecond == 0 {
		obj.Tcp.Handshake23 = true
	} else {
		obj.Tcp.Handshake23 = false
	}

	//obj.Tcp.Open =
	//obj.Tcp.Close

	obj.Http.Host = src.HttpInfo.HttpHost
	obj.Http.Url = src.HttpInfo.HttpUrl
	obj.Http.XonlineHost = src.HttpInfo.HttpOnlineHost
	obj.Http.UserAgent = src.HttpInfo.HttpUserAgent
	obj.Http.ContentType = src.HttpInfo.HttpContent
	obj.Http.Refer = src.HttpInfo.HttpRefer
	obj.Http.Cookie = src.HttpInfo.HttpCookie
	obj.Http.Location = src.HttpInfo.HttpLocation
	obj.Http.Request = string(src.HttpReqInfo)
	obj.Http.Response = string(src.HttpRespInfo)
	obj.Http.RequestTime = src.HttpBaseInfo.Ulactiontime
	obj.Http.FirstResponseTime = src.HttpBaseInfo.Ulfirst_packet_time
	obj.Http.LastContentTime = src.HttpBaseInfo.UlLast_Packet_Time
	obj.Http.ServTime = src.HttpBaseInfo.UlServiceTime
	obj.Http.ContentLen = src.HttpBaseInfo.UlContentLenth
	obj.Http.StateCode = src.HttpBaseInfo.UsHttpStatus
	obj.Http.Method = src.HttpBaseInfo.UcHttpMethod
	obj.Http.Version = src.HttpBaseInfo.Uchttpversion
	obj.Http.HeadFlag = IntToBool(uint32(src.HttpBaseInfo.UcUnionFlag & 0x04))
	obj.Http.ServFlag = uint8(src.HttpBaseInfo.UcUnionFlag & 0x38)
	obj.Http.RequestFlag = IntToBool(uint32(src.HttpBaseInfo.UcUnionFlag & 0xc0))
	obj.Http.Browser = src.HttpBaseInfo.UcIE
	obj.Http.Portal = src.HttpBaseInfo.UcPortal
	obj.Sip.CallingNo = src.SipCaller
	obj.Sip.CalledNo = src.SipCalled
	obj.Sip.SessionId = src.SipSessionId
	obj.Sip.CallDir = src.Sip.UcCallDirection
	obj.Sip.CallType = src.Sip.UcCallType
	obj.Sip.HangupReason = src.Sip.UcHookReason
	obj.Sip.SignalType = src.Sip.UcSignalType
	obj.Sip.StreamCount = src.Sip.UsDataflowNum
	obj.Sip.Malloc = IntToBool(uint32(src.Sip.UbSipIVMR & 0x2000))
	obj.Sip.Bye = IntToBool(uint32(src.Sip.UbSipIVMR & 0x4000))
	obj.Sip.Invite = IntToBool(uint32(src.Sip.UbSipIVMR & 0x8000))
	obj.Rtsp.Url = src.RtspUrl
	obj.Rtsp.UserAgent = src.RtspUserAgent
	obj.Rtsp.ServerIp = src.RtspServerIp
	obj.Rtsp.ClientBeginPort = src.RTSP.UsCStartPort
	obj.Rtsp.ClientEndPort = src.RTSP.UsCEndPort
	obj.Rtsp.ServerBeginPort = src.RTSP.UsSStartPort
	obj.Rtsp.ServerEndPort = src.RTSP.UsSEndPort
	obj.Rtsp.VideoStreamCount = src.RTSP.UsSsnVideoCount
	obj.Rtsp.AudeoStreamCount = src.RTSP.UsSsnAudioCount
	obj.Rtsp.ResDelay = src.RTSP.UlResDelay
	obj.Ftp.State = src.FtpInfo.FtpStatus
	//obj.Ftp.UserCount = src.FtpInfo.FtpUserName//??
	obj.Ftp.CurrentDir = src.FtpInfo.FtpCurDir
	obj.Ftp.TransMode = src.FtpInfo.FtpTransMode
	obj.Ftp.TransType = src.FtpInfo.FtpTransType
	//obj.Ftp.FileCount =
	obj.Ftp.FileSize = src.FtpInfo.FtpFileSize
	obj.Ftp.RspTm = src.FtpInfo.FtpRspTm
	obj.Ftp.TransTm = src.FtpInfo.FtpTransTime
	obj.Mail.MsgType = src.MailInfo.MailMsgType
	obj.Mail.RspState = src.MailInfo.MailRspStatus
	obj.Mail.UserName = src.MailInfo.MailUserName
	obj.Mail.RecverInfo = src.MailInfo.MailResvInfo
	obj.Mail.Len = src.MailInfo.MailLength
	obj.Mail.DomainInfo = src.MailInfo.MailDomainInfo
	obj.Mail.RecvAccount = src.MailInfo.MailRcevAcnt
	obj.Mail.Hdr = src.MailInfo.MailHdr
	obj.Mail.AcsType = src.MailInfo.MailAcsType
	obj.Dns.Domain = src.DnsInfo.DnsZones
	obj.Dns.IpCount = src.DnsInfo.DnsIpNum
	obj.Dns.Ipv4 = src.DnsInfo.DnsIpv4
	obj.Dns.Ipv6 = src.DnsInfo.DnsIpv6
	obj.Dns.RspCode = src.DnsInfo.DnsRspCode
	obj.Dns.ReqCount = src.DnsInfo.DnsRsqCnt
	obj.Dns.RspRecordCount = src.DnsInfo.DnsRspRecordCnt
	obj.Dns.AuthCnttCount = src.DnsInfo.DnsAuthCnttCnt
	obj.Dns.ExtraRecordCount = src.DnsInfo.DnsExtraRecordCnt
	obj.Dns.RspDelay = src.DnsInfo.DnsRspDelay
	//obj.Dns.PktValid =
	//obj.Vpn.Type =
	//obj.Proxy
	//obj.QQ

	//obj.App.ProtoInfo
	//obj.App.Status
	//obj.App.ClassId
	//obj.App.Proto
	obj.App.File = string(src.FileContent)
	info := &BackendInfo{
		Type: src.CheckType(),
		Data: obj,
	}
	return info
}

func TransToBackendObj(origiList []*xdrParse.DpiXdr) []*BackendInfo {
	list := []*BackendInfo{}
	for _, src := range origiList {
		info := PerTransToBackendObj(src)
		list = append(list, info)
	}
	return list
}

/*
func (this *BackendObj) CheckType() int {
	//TODO

	if len(this.App.File) != 0 {
		return XdrFileType
	} else if len(this.Http.Request) != 0 {
		return XdrHttpType
	} else {
		return XdrType
	}
	return -1
}
*/
func (this *BackendObj) HashPartation() uint32 {
	var divisor uint32 = 3
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, []byte(this.Conn.Sip))
	binary.Write(buf, binary.LittleEndian, []byte(this.Conn.Dip))
	binary.Write(buf, binary.LittleEndian, this.Conn.Dport)
	binary.Write(buf, binary.LittleEndian, this.Conn.Sport)
	binary.Write(buf, binary.LittleEndian, this.Conn.Proto)
	//bufbytes := buf.Bytes()
	//mlog.Debug("bufbytes len=", len(bufbytes))
	var n uint8
	var sum uint32 = 0
	length := buf.Len()
	for i := 0; i < length; i++ {
		binary.Read(buf, binary.LittleEndian, &n)
		sum += uint32(n) * divisor
	}
	/*
		//init topic partition
		sum := Ipv4StringToInt(this.Conn.Sip) +
			Ipv4StringToInt(this.Conn.Dip) +
			uint32(this.Conn.Dport) +
			uint32(this.Conn.Sport) +
			uint32(this.Conn.Proto)
	*/
	return uint32(uint32(sum) % uint32(PartitionNum))
}
