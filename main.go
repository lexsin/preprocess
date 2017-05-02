package main

import (
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"

	"github.com/howeyc/fsnotify"
)

type FileEvent *fsnotify.FileEvent

func RunNotify(dir string, handle func(ev *fsnotify.FileEvent) error) {
	Watcher, err := fsnotify.NewWatcher()
	if err != nil {
		mlog.Error(err)
		panic("fsnotify.NewWatcher() error:" + err.Error())
	}

	go func() {
		err := Watcher.Watch(dir)
		if err != nil {
			mlog.Error(err)
		}
	}()
	mlog.Info("begin watch dir:", dir)
	for {
		select {
		case ev := <-Watcher.Event:
			mlog.Debug("event:", ev)
			/*
				case ev := <-Watcher.Event:
					if ev.IsCreate() {
						go handle(ev)
					}
			*/
		case err := <-Watcher.Error:
			mlog.Error(err)
		}
	}

	/* ... do stuff ... */
	Watcher.Close()
}

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
