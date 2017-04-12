package main

import (
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"

	"github.com/howeyc/fsnotify"
)

func RunNotify(dir string, handle func(ev *fsnotify.FileEvent) error) {
	//var Watcher *fsnotify.Watcher
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
	dpiDir, _ := mconfig.Conf.String("dir", "DpiDir")
	go RunNotify(dpiDir, DpiHandle)

	vdsAlertDir, _ := mconfig.Conf.String("dir", "vdsAlertDir")
	go RunNotify(vdsAlertDir, VdsAlertHandler)

	idsAlertDir, _ := mconfig.Conf.String("dir", "IdsAlertDir")
	go RunNotify(idsAlertDir, IdsAlertHandler)
	//block
	for {
	}
}
