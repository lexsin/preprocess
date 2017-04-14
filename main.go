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

	dpiDir, _ := mconfig.Conf.String("dir", "DpiXdrDir")
	go RunNotify(dpiDir, DpiHandle)

	vdsAlertDir, _ := mconfig.Conf.String("dir", "VdsAlertDir")
	go RunNotify(vdsAlertDir, VdsAlertHandler)

	idsAlertDir, _ := mconfig.Conf.String("dir", "IdsAlertDir")
	go RunNotify(idsAlertDir, IdsAlertHandler)

	//begin waf-alert http server
	wafAlertServer()

	//block
	<-block
}
