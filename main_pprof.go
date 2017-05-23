package main

/*
import (
	"flag"
	"log"
	"os"
	"os/signal"
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"
	"runtime/pprof"
	//"syscall"
	//"runtime"
)


func main() {
	startPProf()

	DpiWatchDir, _ := mconfig.Conf.String("dir", "DpiXdrDir")
	CreateDir(DpiWatchDir)
	go notify_ftp_mv(DpiWatchDir, DpiHandle)

	VdsAlertWatchDir, _ := mconfig.Conf.String("dir", "VdsAlertDir")
	CreateDir(VdsAlertWatchDir)
	go notify_ftp_mv(VdsAlertWatchDir, VdsAlertHandler)

	IdsAlertWatchDir, _ := mconfig.Conf.String("dir", "IdsAlertDir")
	CreateDir(IdsAlertWatchDir)
	//go RunNotify(IdsAlertWatchDir, IdsAlertHandler)
	go notify_ftp_mv(IdsAlertWatchDir, IdsAlertHandler)

	//begin waf-alert http server
	go RunWafServer()

	c := make(chan os.Signal)
	signal.Notify(c)
	s := <-c
	mlog.Info("exit by signal:", s.String())
	startMempro()
}

func startPProf() {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

}

func startMempro() {
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	flag.Parse()
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		//runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
*/
