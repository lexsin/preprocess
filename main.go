package main

import (
	"preprocess/modules/mconfig"
	//"preprocess/modules/mlog"
)

func main() {
	var block chan int

	DpiWatchDir, _ := mconfig.Conf.String("dir", "DpiXdrDir")
	CreateDir(DpiWatchDir)
	//go RunNotify(DpiWatchDir, DpiHandle)
	go notify_ftp_mv(DpiWatchDir, DpiHandle)

	VdsAlertWatchDir, _ := mconfig.Conf.String("dir", "VdsAlertDir")
	CreateDir(VdsAlertWatchDir)
	//go RunNotify(VdsAlertWatchDir, VdsAlertHandler)
	go notify_ftp_mv(VdsAlertWatchDir, VdsAlertHandler)

	IdsAlertWatchDir, _ := mconfig.Conf.String("dir", "IdsAlertDir")
	CreateDir(IdsAlertWatchDir)
	//go RunNotify(IdsAlertWatchDir, IdsAlertHandler)
	go notify_ftp_mv(IdsAlertWatchDir, IdsAlertHandler)

	//begin waf-alert http server
	RunWafServer()

	//block
	<-block
}
