package main

import (
	"fmt"
	"log"
	"preproccess/modules/mconfig"
	"preproccess/modules/mlog"

	"github.com/howeyc/fsnotify"
)

var Watcher *fsnotify.Watcher

func init() {
	Watcher, err := fsnotify.NewWatcher()
	if err != nil {
		mlog.Error(err)
	}
}

func main() {
	go func() {
		err := Watcher.Watch("testDir")
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		select {
		case ev := <-Watcher.Event:
			if ev.IsCreate() {
				mlog.Debug(fmt.Println("Create file:", ev.Name))
			}
		case err := <-Watcher.Error:
			mlog.Error(err)
		}
	}

	/* ... do stuff ... */
	watcher.Close()
}
