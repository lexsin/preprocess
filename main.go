package main

import (
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"

	"github.com/howeyc/fsnotify"
)

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
			if ev.IsCreate() {
				go handle(ev)
			}
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
	go RunNotify(DpiWatchDir, DpiHandle)

	VdsAlertWatchDir, _ := mconfig.Conf.String("dir", "VdsAlertDir")
	go RunNotify(VdsAlertWatchDir, VdsAlertHandler)

	IdsAlertWatchDir, _ := mconfig.Conf.String("dir", "IdsAlertDir")
	go RunNotify(IdsAlertWatchDir, IdsAlertHandler)

	//begin waf-alert http server
	go RunWafServer()

	//block
	<-block
}
