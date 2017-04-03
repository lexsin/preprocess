package xdrParse

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	XDR_HEAD = iota
	XDR_SESSION_STATUS
	XDR_APP_ID
	XDR_TUPLE
	XDR_SESSION_STAT
	XDR_SESSION_TIME
	XDR_BUSSI_STAT
	XDR_TCP_INFO
	XDR_RESPONSE_DELAY
	XDR_LEVEL_7_PROTO
	XDR_HTTP_BASE_INFO = 10
	XDR_HTTP_HOST
	XDR_HTTP_URL
	XDR_HTTP_ONLINE_HOST
	XDR_HTTP_USER_AGENT
	XDR_HTTP_CONTENT
	XDR_HTTP_REFER
	XDR_HTTP_COOKIE
	XDR_HTTP_LOCATION
	XDR_SIP
	XDR_SIP_CALLER = 20
	XDR_SIP_CALLED
	XDR_SIP_SESSION_ID
	XDR_RTSP
	XDR_RTSP_URL
	XDR_RTSP_USER_AGENT
	XDR_RTSP_SERVER_IP
	//FTP_INFO
	XDR_FTP_STATUS
	XDR_FTP_USER_NAME
	XDR_FTP_CUR_DIR
	XDR_FTP_TRANS_MODE = 30
	XDR_FTP_TRANS_TYPE
	XDR_FTP_FILE_NAME
	XDR_FTP_FILE_SIZE
	XDR_FTP_RSP_TM
	XDR_FTP_TRANS_TIME
	//MAIL_INFO
	XDR_MAIL_MSG_TYPE
	XDR_MAIL_RSP_STATUS
	XDR_MAIL_USER_NAME
	XDR_MAIL_RESV_INFO
	XDR_MAIL_LENGTH = 40
	XDR_MAIL_DOMAIN_INFO
	XDR_MAIL_RESV_ACNT
	XDR_MAIL_HDR
	XDR_MAIL_ACS_TYPE
	//DNS_INFO
	XDR_DNS_ZONES
	XDR_DNS_IP_NUM
	XDR_DNS_IPV4
	XDR_DNS_IPV6
	XDR_DNS_RSP_CODE
	XDR_DNS_RSQ_CNT = 50
	XDR_DNS_RSP_RECODE_CNT
	XDR_DNS_AUTH_CNTT_CNT
	XDR_DNS_EXTRA_RECORD_CNT
	XDR_DNS_RSP_DELAY  = 54
	XDR_HTTP_REQ_INFO  = 201
	XDR_HTTP_RESP_INFO = 202
	XDR_FILE_CONTENT   = 203
)

func normalDataDecode(data []byte, bin interface{}) error {
	buf := new(bytes.Buffer)
	buf.Write(data)
	if err := binary.Read(buf, binary.LittleEndian, bin); err != nil {
		return err
	}
	fmt.Println()
	return nil
}
func ParsSessionStatus(xdr *TlvValue, obj *DpiXdr) error {
	if xdr.ShortData != 0 && xdr.ShortData != 1 {
		return errors.New("XDR_SESSION_STATUS parse != 0/1")
	}
	obj.SessionStatus = xdr.ShortData
	return nil
}

func parsAppId(xdr *TlvValue, obj *DpiXdr) error {
	if xdr.ShortData < 100 && xdr.ShortData > 108 {
		return errors.New("XDR_APP_ID out range")
	}
	obj.AppId = xdr.ShortData
	return nil
}

func parsTuple(xdr *TlvValue, obj *DpiXdr) error {
	if err := normalDataDecode(xdr.Data, &obj.Tuple); err != nil {
		return errors.New(fmt.Sprintf("XDR_TUPLE error:%s", err.Error()))
	}
	return nil
}

func parsSessionStat(xdr *TlvValue, obj *DpiXdr) error {
	if err := normalDataDecode(xdr.Data, &obj.SesionStat); err != nil {
		return errors.New(fmt.Sprintf("XDR_TUPLE error:%s", err.Error()))
	}
	return nil
}

func parsSessionTime(xdr *TlvValue, obj *DpiXdr) error {
	if err := normalDataDecode(xdr.Data, &obj.SesionTime); err != nil {
		return errors.New(fmt.Sprintf("XDR_SESSION_TIME error:%s", err.Error()))
	}
	return nil
}

func parsBussiStat(xdr *TlvValue, obj *DpiXdr) error {
	if err := normalDataDecode(xdr.Data, &obj.SesionStat); err != nil {
		return errors.New(fmt.Sprintf("XDR_BUSSI_STAT error:%s", err.Error()))
	}
	return nil
}

func parsTcpInfo(xdr *TlvValue, obj *DpiXdr) error {
	if err := normalDataDecode(xdr.Data, &obj.TcpInfo); err != nil {
		return errors.New(fmt.Sprintf("XDR_TCP_INFO error:%s", err.Error()))
	}
	return nil
}

func parsResponseDelay(xdr *TlvValue, obj *DpiXdr) error {
	if err := normalDataDecode(xdr.Data, &obj.ResponseDelay); err != nil {
		return errors.New(fmt.Sprintf("XDR_RESPONSE_DELAY error:%s", err.Error()))
	}
	return nil
}

func parsLevel7Proto(xdr *TlvValue, obj *DpiXdr) error {
	//9
	if err := normalDataDecode(xdr.Data, &obj.Level7Proto); err != nil {
		return errors.New(fmt.Sprintf("XDR_LEVEL_7_PROTO error:%s", err.Error()))
	}
	return nil
}

func parsHttpBaseInfo(xdr *TlvValue, obj *DpiXdr) error {
	//10
	if err := normalDataDecode(xdr.Data, &obj.HttpBaseInfo); err != nil {
		return errors.New(fmt.Sprintf("XDR_HTTP_BASE_INFO error:%s", err.Error()))
	}
	return nil
}

func parsHttpHost(xdr *TlvValue, obj *DpiXdr) error {
	//11
	obj.HttpInfo.HttpHost = string(xdr.Data)
	return nil
}

func parsHttpUrl(xdr *TlvValue, obj *DpiXdr) error {
	//12
	obj.HttpInfo.HttpUrl = string(xdr.Data)
	return nil
}

func parsHttpOnlineHost(xdr *TlvValue, obj *DpiXdr) error {
	//13
	obj.HttpInfo.HttpOnlineHost = string(xdr.Data)
	return nil
}

func parsHttpUserAgent(xdr *TlvValue, obj *DpiXdr) error {
	//14
	obj.HttpInfo.HttpUserAgent = string(xdr.Data)
	return nil
}

func parsHttpContent(xdr *TlvValue, obj *DpiXdr) error {
	//15
	obj.HttpInfo.HttpContent = string(xdr.Data)
	return nil
}

func parsHttpRefer(xdr *TlvValue, obj *DpiXdr) error {
	//16
	obj.HttpInfo.HttpRefer = string(xdr.Data)
	return nil
}

func parsHttpCookie(xdr *TlvValue, obj *DpiXdr) error {
	//17
	obj.HttpInfo.HttpCookie = string(xdr.Data)
	return nil
}

func parsHttpLocation(xdr *TlvValue, obj *DpiXdr) error {
	//18
	obj.HttpInfo.HttpLocation = string(xdr.Data)
	return nil
}

func parsSip(xdr *TlvValue, obj *DpiXdr) error {
	//19
	if err := normalDataDecode(xdr.Data, &obj.Sip); err != nil {
		return errors.New(fmt.Sprintf("XDR_SIP error:%s", err.Error()))
	}
	return nil
}

func parsSipCaller(xdr *TlvValue, obj *DpiXdr) error {
	//20
	obj.SipCaller = string(xdr.Data)
	return nil
}

func parsSipCalled(xdr *TlvValue, obj *DpiXdr) error {
	//21
	obj.SipCalled = string(xdr.Data)
	return nil
}

func parsSipSessionId(xdr *TlvValue, obj *DpiXdr) error {
	//22
	obj.SipSessionId = string(xdr.Data)
	return nil
}

func parsRtsp(xdr *TlvValue, obj *DpiXdr) error {
	//23
	if err := normalDataDecode(xdr.Data, &obj.RTSP); err != nil {
		return errors.New(fmt.Sprintf("XDR_RTSP error:%s", err.Error()))
	}
	return nil
}

func parsRtspUrl(xdr *TlvValue, obj *DpiXdr) error {
	//24
	obj.RtspUrl = string(xdr.Data)
	return nil
}

func parsRtspUserAgent(xdr *TlvValue, obj *DpiXdr) error {
	//25
	obj.RtspUserAgent = string(xdr.Data)
	return nil
}

func parsRtspServerIp(xdr *TlvValue, obj *DpiXdr) error {
	//26
	obj.RtspServerIp = string(xdr.Data)
	return nil
}

func parsFtpStatus(xdr *TlvValue, obj *DpiXdr) error {
	//27
	if err := normalDataDecode(xdr.Data, &obj.FtpInfo.FtpStatus); err != nil {
		return errors.New(fmt.Sprintf("XDR_FTP_STATUS error:%s", err.Error()))
	}
	return nil
}

func parsFtpUserName(xdr *TlvValue, obj *DpiXdr) error {
	//28
	obj.FtpInfo.FtpUserName = string(xdr.Data)
	return nil
}

func parsFtpCurDir(xdr *TlvValue, obj *DpiXdr) error {
	//29
	obj.FtpInfo.FtpCurDir = string(xdr.Data)
	return nil
}

func parsFtpTransMode(xdr *TlvValue, obj *DpiXdr) error {
	//30
	if err := normalDataDecode(xdr.Data, &obj.FtpInfo.FtpTransMode); err != nil {
		return errors.New(fmt.Sprintf("XDR_FTP_TRANS_MODE error:%s", err.Error()))
	}
	return nil
}

func parsFtpTransType(xdr *TlvValue, obj *DpiXdr) error {
	//31
	if err := normalDataDecode(xdr.Data, &obj.FtpInfo.FtpTransType); err != nil {
		return errors.New(fmt.Sprintf("XDR_FTP_TRANS_TYPE error:%s", err.Error()))
	}
	return nil
}

func parsFtpFileName(xdr *TlvValue, obj *DpiXdr) error {
	//32
	obj.FtpInfo.FtpFileName = string(xdr.Data)
	return nil
}

func parsFtpFileSize(xdr *TlvValue, obj *DpiXdr) error {
	//33
	if err := normalDataDecode(xdr.Data, &obj.FtpInfo.FtpFileSize); err != nil {
		return errors.New(fmt.Sprintf("XDR_FTP_FILE_SIZE error:%s", err.Error()))
	}
	return nil
}

func parsFtpRspTm(xdr *TlvValue, obj *DpiXdr) error {
	//34
	if err := normalDataDecode(xdr.Data, &obj.FtpInfo.FtpRspTm); err != nil {
		return errors.New(fmt.Sprintf("XDR_FTP_RSP_TM error:%s", err.Error()))
	}
	return nil
}

func parsFtpTransTime(xdr *TlvValue, obj *DpiXdr) error {
	//35
	if err := normalDataDecode(xdr.Data, &obj.FtpInfo.FtpTransTime); err != nil {
		return errors.New(fmt.Sprintf("XDR_FTP_TRANS_TIME error:%s", err.Error()))
	}
	return nil
}

func parsMailMsgType(xdr *TlvValue, obj *DpiXdr) error {
	//36
	if err := normalDataDecode(xdr.Data, &obj.MailInfo.MailMsgType); err != nil {
		return errors.New(fmt.Sprintf("XDR_MAIL_MSG_TYPE error:%s", err.Error()))
	}
	return nil
}

func parsMailRspStatus(xdr *TlvValue, obj *DpiXdr) error {
	//37
	if err := normalDataDecode(xdr.Data, &obj.MailInfo.MailRspStatus); err != nil {
		return errors.New(fmt.Sprintf("XDR_MAIL_RSP_STATUS error:%s", err.Error()))
	}
	return nil
}

func parsMailUserName(xdr *TlvValue, obj *DpiXdr) error {
	//38
	obj.MailInfo.MailUserName = string(xdr.Data)
	return nil
}

func parsMailResvInfo(xdr *TlvValue, obj *DpiXdr) error {
	//39
	obj.MailInfo.MailResvInfo = string(xdr.Data)
	return nil
}

func parsMailLength(xdr *TlvValue, obj *DpiXdr) error {
	//40
	if err := normalDataDecode(xdr.Data, &obj.MailInfo.MailLength); err != nil {
		return errors.New(fmt.Sprintf("XDR_MAIL_RSP_STATUS error:%s", err.Error()))
	}
	return nil
}

func parsMailDomainInfo(xdr *TlvValue, obj *DpiXdr) error {
	//41
	obj.MailInfo.MailDomainInfo = string(xdr.Data)
	return nil
}

func parsMailRcevAcnt(xdr *TlvValue, obj *DpiXdr) error {
	//42
	obj.MailInfo.MailResvInfo = string(xdr.Data)
	return nil
}

func parsMailHdr(xdr *TlvValue, obj *DpiXdr) error {
	//43
	obj.MailInfo.MailHdr = string(xdr.Data)
	return nil
}

func parsMailAcsType(xdr *TlvValue, obj *DpiXdr) error {
	//44
	if err := normalDataDecode(xdr.Data, &obj.MailInfo.MailAcsType); err != nil {
		return errors.New(fmt.Sprintf("XDR_MAIL_ACS_TYPE error:%s", err.Error()))
	}
	return nil
}

func parsDnsZones(xdr *TlvValue, obj *DpiXdr) error {
	//45
	obj.DnsInfo.DnsZones = string(xdr.Data)
	return nil
}

func parsDnsIpNum(xdr *TlvValue, obj *DpiXdr) error {
	//46
	if err := normalDataDecode(xdr.Data, &obj.DnsInfo.DnsIpNum); err != nil {
		return errors.New(fmt.Sprintf("XDR_DNS_IP_NUM error:%s", err.Error()))
	}
	return nil
}

func parsDnsIpv4(xdr *TlvValue, obj *DpiXdr) error {
	//47
	obj.DnsInfo.DnsIpv4 = string(xdr.Data)
	return nil
}

func parsDnsIpv6(xdr *TlvValue, obj *DpiXdr) error {
	//48
	obj.DnsInfo.DnsIpv6 = string(xdr.Data)
	return nil
}

func parsDnsRspCode(xdr *TlvValue, obj *DpiXdr) error {
	//49
	if err := normalDataDecode(xdr.Data, &obj.DnsInfo.DnsRspCode); err != nil {
		return errors.New(fmt.Sprintf("XDR_DNS_RSP_CODE error:%s", err.Error()))
	}
	return nil
}

func parsDnsRsqCnt(xdr *TlvValue, obj *DpiXdr) error {
	//50
	if err := normalDataDecode(xdr.Data, &obj.DnsInfo.DnsRsqCnt); err != nil {
		return errors.New(fmt.Sprintf("XDR_DNS_RSQ_CNT error:%s", err.Error()))
	}
	return nil
}

func parsDnsRspRecordCnt(xdr *TlvValue, obj *DpiXdr) error {
	//51
	if err := normalDataDecode(xdr.Data, &obj.DnsInfo.DnsRspRecordCnt); err != nil {
		return errors.New(fmt.Sprintf("XDR_DNS_RSP_RECODE_CNT error:%s", err.Error()))
	}
	return nil
}

func parsDNSAuthCnttCnt(xdr *TlvValue, obj *DpiXdr) error {
	//52
	if err := normalDataDecode(xdr.Data, &obj.DnsInfo.DnsAuthCnttCnt); err != nil {
		return errors.New(fmt.Sprintf("XDR_DNS_AUTH_CNTT_CNT error:%s", err.Error()))
	}
	return nil
}

func parsDnsExtraRecordCnt(xdr *TlvValue, obj *DpiXdr) error {
	//53
	if err := normalDataDecode(xdr.Data, &obj.DnsInfo.DnsExtraRecordCnt); err != nil {
		return errors.New(fmt.Sprintf("XDR_DNS_EXTRA_RECORD_CNT error:%s", err.Error()))
	}
	return nil
}

func parsDnsRspDelay(xdr *TlvValue, obj *DpiXdr) error {
	//54
	if err := normalDataDecode(xdr.Data, &obj.DnsInfo.DnsRspDelay); err != nil {
		return errors.New(fmt.Sprintf("XDR_DNS_RSP_DELAY error:%s", err.Error()))
	}
	return nil
}

func parsHttpRespInfo(xdr *TlvValue, obj *DpiXdr) error {
	//202
	obj.HttpRespInfo = xdr.Data
	return nil
}

func parsHttpReqInfo(xdr *TlvValue, obj *DpiXdr) error {
	//203
	obj.HttpReqInfo = xdr.Data
	return nil
}

func parsFileContent(xdr *TlvValue, obj *DpiXdr) error {
	//204
	obj.FileContent = xdr.Data
	return nil
}
