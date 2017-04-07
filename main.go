package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"
	"preprocess/modules/xdrParse"
	"strings"

	"preprocess/modules/pushkafka"

	"github.com/howeyc/fsnotify"
)

func main() {
	dpiDir, _ := mconfig.Conf.String("dir", "WatchDir")
	RunNotify(dpiDir, DpiHandle)
}

func RunNotify(dir string, handle func(ev *fsnotify.FileEvent) error) {
	go func() {
		err := Watcher.Watch(dir)
		if err != nil {
			log.Fatal(err)
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

func ReadFile(FileName string) ([]byte, error) {
	fileContent, err := ioutil.ReadFile(FileName)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

func CheckSuffix(FileName string) bool {
	f := func(c rune) bool {
		return c == '.'
	}
	ss := strings.FieldsFunc(FileName, f)
	suffix := ss[len(ss)-1]
	if suffix != "xdr" && suffix != "XDR" {
		return false
	}
	return true
}

func DpiHandle(ev *fsnotify.FileEvent) error {
	defer func() {
		if err := recover(); err != nil {
			mlog.Error(err)
		}
	}()
	mlog.Debug(fmt.Println("Create file:", ev.Name))

	//check file suffix
	if ok := CheckSuffix(ev.Name); !ok {
		panic(fmt.Sprintf("file: %s suffix error!", ev.Name))
	}

	//read file
	content, err := ReadFile(ev.Name)
	if err != nil {
		panic(fmt.Sprintf("read file %s error:%s", ev.Name, err.Error()))
	}

	//XDR==>struct
	datalist, err := xdrParse.ParseXdr(content)
	if err != nil {
		panic(fmt.Sprintf("parse file %s Xdr error:%s", ev.Name, err.Error()))
	}

	//trans to backend obj
	backlist := TransToBackendObj(datalist)

	//pre process file into ceph
	//TODO

	//push to kafka
	for _, datap := range backlist {
		//mlog.Debug("data=", datap)
		go DoPushTopic(datap)
	}

	return nil
}

func DoPushTopic(datap *BackendObj) error {
	//struct==>json
	jsonstr, _ := json.Marshal(*datap)
	mlog.Debug("DoPushTopic json:", string(jsonstr))

	//json==>topic
	dtype := datap.CheckType()

	mm, ok := TopicMap[dtype]
	if !ok {
		mlog.Error("TopicMap %d not exist", dtype)
		return errors.New(fmt.Sprintf("topicMap[%s] not exist", dtype))
	}
	topic := &TopicType{
		topicName: mm.topicName,
		handlePre: mm.handlePre,
		origiData: jsonstr,
		partition: int(datap.HashPartation()),
	}
	if err := pushkafka.PushKafka(topic); err != nil {
		return err
	}
	return nil
}
