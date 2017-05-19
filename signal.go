package main

import (
	"os"
	"os/signal"
	"preprocess/modules/mlog"
	"syscall"
)

func safeExit() {
	signalListen()
}

func signalListen() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-c
	sigExitHandler(s)
}

func sigExitHandler(s os.Signal) {
	mlog.Info("get signal:", s.String())
	mlog.Info("pregress will exit!")
	//os.Exit(1)
}
