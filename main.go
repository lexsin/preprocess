package main

import (
	"preprocess/modules/mconfig"
	//"preprocess/modules/mlog"
)

func main() {
	var block chan int

	DpiWatchDir, _ := mconfig.Conf.String("dir", "DpiXdrDir")
	CreateDir(DpiWatchDir)
	go RunNotify(DpiWatchDir, DpiHandle)

	VdsAlertWatchDir, _ := mconfig.Conf.String("dir", "VdsAlertDir")
	CreateDir(VdsAlertWatchDir)
	go RunNotify(VdsAlertWatchDir, VdsAlertHandler)

	IdsAlertWatchDir, _ := mconfig.Conf.String("dir", "IdsAlertDir")
	CreateDir(IdsAlertWatchDir)
	go RunNotify(IdsAlertWatchDir, IdsAlertHandler)

	//begin waf-alert http server
	RunWafServer()

	//block
	<-block
}
