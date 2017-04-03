package xdrParse

var DecodeFuncMap map[int]func(xdr *TlvValue, obj *DpiXdr) error

func init() {
	DecodeFuncMap = make(map[int]func(xdr *TlvValue, obj *DpiXdr) error, 0)

	DecodeFuncMap[XDR_SESSION_STATUS] = ParsSessionStatus
	DecodeFuncMap[XDR_APP_ID] = parsAppId
	DecodeFuncMap[XDR_TUPLE] = parsTuple
	DecodeFuncMap[XDR_SESSION_STAT] = parsSessionStat
	DecodeFuncMap[XDR_SESSION_TIME] = parsSessionTime
	DecodeFuncMap[XDR_BUSSI_STAT] = parsBussiStat
	DecodeFuncMap[XDR_TCP_INFO] = parsTcpInfo
	DecodeFuncMap[XDR_RESPONSE_DELAY] = parsResponseDelay
	DecodeFuncMap[XDR_LEVEL_7_PROTO] = parsLevel7Proto
	DecodeFuncMap[XDR_HTTP_BASE_INFO] = parsHttpBaseInfo
	DecodeFuncMap[XDR_HTTP_HOST] = parsHttpHost
	DecodeFuncMap[XDR_HTTP_URL] = parsHttpUrl
	DecodeFuncMap[XDR_HTTP_ONLINE_HOST] = parsHttpOnlineHost
	DecodeFuncMap[XDR_HTTP_USER_AGENT] = parsHttpUserAgent
	DecodeFuncMap[XDR_HTTP_CONTENT] = parsHttpContent
	DecodeFuncMap[XDR_HTTP_REFER] = parsHttpRefer
	DecodeFuncMap[XDR_HTTP_COOKIE] = parsHttpCookie
	DecodeFuncMap[XDR_HTTP_LOCATION] = parsHttpLocation
	DecodeFuncMap[XDR_SIP] = parsSip
	DecodeFuncMap[XDR_SIP_CALLER] = parsSipCaller
	DecodeFuncMap[XDR_SIP_CALLED] = parsSipCalled
	DecodeFuncMap[XDR_SIP_SESSION_ID] = parsSipSessionId
	DecodeFuncMap[XDR_RTSP] = parsRtsp
	DecodeFuncMap[XDR_RTSP_URL] = parsRtspUrl
	DecodeFuncMap[XDR_RTSP_USER_AGENT] = parsRtspUserAgent
	DecodeFuncMap[XDR_RTSP_SERVER_IP] = parsRtspServerIp
	DecodeFuncMap[XDR_FTP_STATUS] = parsFtpStatus
	DecodeFuncMap[XDR_FTP_USER_NAME] = parsFtpUserName
	DecodeFuncMap[XDR_FTP_CUR_DIR] = parsFtpCurDir
	DecodeFuncMap[XDR_FTP_TRANS_MODE] = parsFtpTransMode
	DecodeFuncMap[XDR_FTP_TRANS_TYPE] = parsFtpTransType
	DecodeFuncMap[XDR_FTP_FILE_NAME] = parsFtpFileName
	DecodeFuncMap[XDR_FTP_FILE_SIZE] = parsFtpFileSize
	DecodeFuncMap[XDR_FTP_RSP_TM] = parsFtpRspTm
	DecodeFuncMap[XDR_FTP_TRANS_TIME] = parsFtpTransTime
	DecodeFuncMap[XDR_MAIL_MSG_TYPE] = parsMailMsgType
	DecodeFuncMap[XDR_MAIL_RSP_STATUS] = parsMailRspStatus
	DecodeFuncMap[XDR_MAIL_USER_NAME] = parsMailUserName
	DecodeFuncMap[XDR_MAIL_RESV_INFO] = parsMailResvInfo
	DecodeFuncMap[XDR_MAIL_LENGTH] = parsMailLength
	DecodeFuncMap[XDR_MAIL_DOMAIN_INFO] = parsMailDomainInfo
	DecodeFuncMap[XDR_MAIL_RESV_ACNT] = parsMailRcevAcnt
	DecodeFuncMap[XDR_MAIL_HDR] = parsMailHdr
	DecodeFuncMap[XDR_MAIL_ACS_TYPE] = parsMailAcsType
	DecodeFuncMap[XDR_DNS_ZONES] = parsDnsZones
	DecodeFuncMap[XDR_DNS_IP_NUM] = parsDnsIpNum
	DecodeFuncMap[XDR_DNS_IPV4] = parsDnsIpv4
	DecodeFuncMap[XDR_DNS_IPV6] = parsDnsIpv6
	DecodeFuncMap[XDR_DNS_RSP_CODE] = parsDnsRspCode
	DecodeFuncMap[XDR_DNS_RSQ_CNT] = parsDnsRsqCnt
	DecodeFuncMap[XDR_DNS_RSP_RECODE_CNT] = parsDnsRspRecordCnt
	DecodeFuncMap[XDR_DNS_AUTH_CNTT_CNT] = parsDNSAuthCnttCnt
	DecodeFuncMap[XDR_DNS_EXTRA_RECORD_CNT] = parsDnsExtraRecordCnt
	DecodeFuncMap[XDR_DNS_RSP_DELAY] = parsDnsRspDelay
	DecodeFuncMap[XDR_HTTP_REQ_INFO] = parsHttpRespInfo
	DecodeFuncMap[XDR_HTTP_RESP_INFO] = parsHttpReqInfo
	DecodeFuncMap[XDR_FILE_CONTENT] = parsFileContent

	errInit()
}
