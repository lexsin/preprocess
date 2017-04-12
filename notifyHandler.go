package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	//"log"
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"
	"preprocess/modules/xdrParse"
	"strings"

	"preprocess/modules/pushkafka"

	"github.com/howeyc/fsnotify"
)

func IdsAlertHandler(ev *fsnotify.FileEvent) error {
	topicname, _ := mconfig.Conf.String("kafka", "IdsAlertTopicName")
	return AlertHandler(ev.Name, topicname)
}

func VdsAlertHandler(ev *fsnotify.FileEvent) error {
	topicname, _ := mconfig.Conf.String("kafka", "VdsAlertTopicName")
	return AlertHandler(ev.Name, topicname)
}

func AlertHandler(fileName string, topicName string) error {
	defer func() {
		if err := recover(); err != nil {
			mlog.Error(err)
		}
	}()
	mlog.Debug(fmt.Println("Create file:", fileName))

	//check file suffix
	if ok := CheckSuffix(fileName, "alert"); !ok {
		panic(fmt.Sprintf("file: %s suffix error!", fileName))
	}

	//push kafka

	pushkafkaFunc := func(line string) error {
		data := &DataType{
			topicName: topicName,
			//handlePre: mm.handlePre,
			origiData: []byte(line),
			partition: 0,
		}
		if err := pushkafka.PushKafka(data); err != nil {
			return err
		}
		return nil
	}
	//read file
	DealFilePerline(fileName, pushkafkaFunc)
	return nil
}

func DpiHandle(ev *fsnotify.FileEvent) error {
	defer func() {
		if err := recover(); err != nil {
			mlog.Error(err)
		}
	}()
	mlog.Debug(fmt.Println("Create file:", ev.Name))

	//check file suffix
	if ok := CheckSuffix(ev.Name, []string{"xdr", "XDR"}...); !ok {
		panic(fmt.Sprintf("file: %s suffix error!", ev.Name))
	}

	//read file
	content, err := ReadFile(ev.Name)
	if err != nil {
		panic(fmt.Sprintf("read file %s error:%s", ev.Name, err.Error()))
	}

	//XDR==>pre object
	datalist, err := xdrParse.ParseXdr(content)
	if err != nil {
		panic(fmt.Sprintf("parse file %s Xdr error:%s", ev.Name, err.Error()))
	}
	//file to ceph
	saveToCeph(datalist)
	//pre object ==> backend object
	backlist := TransToBackendObj(datalist)

	//push to kafka
	for _, backObj := range backlist {
		//mlog.Debug("data=", datap)
		go DoPushTopic(backObj)
	}

	return nil
}

func ReadFile(FileName string) ([]byte, error) {
	fileContent, err := ioutil.ReadFile(FileName)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

func CheckSuffix(FileName string, suffixs ...string) bool {
	f := func(c rune) bool {
		return c == '.'
	}
	ss := strings.FieldsFunc(FileName, f)
	suffix := ss[len(ss)-1]
	for _, stand := range suffixs {
		if suffix == stand {
			return true
		}
	}
	return false
	/*
		if suffix != "xdr" && suffix != "XDR" {
			return false
		}
		return true
	*/
}

func DoPushTopic(backObj *BackendInfo) error {
	//struct==>json
	jsonstr, _ := json.Marshal(backObj.Data)
	//mlog.Debug("DoPushTopic json:", string(jsonstr))

	//json==>topic
	mm, ok := TopicMap[backObj.Type]
	mlog.Debug("TopicMap[", backObj.Type, "]=", TopicMap[backObj.Type])
	if !ok {
		mlog.Error("TopicMap %d not exist", backObj.Type)
		return errors.New(fmt.Sprintf("topicMap[%s] not exist", backObj.Type))
	}
	topic := &DataType{
		topicName: mm.topicName,
		handlePre: mm.handlePre,
		origiData: jsonstr,
		partition: int(backObj.Data.HashPartation()),
	}
	if err := pushkafka.PushKafka(topic); err != nil {
		return err
	}

	return nil
}
