// +build linux

package main

import (
	"errors"
	"fmt"
	"preprocess/modules/mlog"

	"github.com/rjeczalik/notify"
)

func notify_ftp_mv(dir string, handle func(filename string) error) {
	var err error
	var ei notify.EventInfo

	c := make(chan notify.EventInfo, 1)
	defer notify.Stop(c)
	mlog.Info("begin watch dir:", dir)

	for {
		if err = notify.Watch(dir, c, notify.InCloseWrite, notify.InMovedTo); err != nil {
			panic(errors.New(fmt.Sprintln("notify watch dir", dir, err.Error())))
		}

		switch ei = <-c; ei.Event() {
		case notify.InCloseWrite, notify.InMovedTo:
			handle(ei.Path())
		default:
			mlog.Error(fmt.Println("notify get event:", ei))
		}
	}
}
