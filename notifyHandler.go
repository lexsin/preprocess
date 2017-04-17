package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"
	"preprocess/modules/pushkafka"
	"preprocess/modules/xdrParse"
	"time"

	"github.com/howeyc/fsnotify"
)

func IdsAlertHandler(ev *fsnotify.FileEvent) error {
	topicname, _ := mconfig.Conf.String("kafka", "IdsAlertTopicName")
	return AlertHandler(ev.Name, topicname, "xdr")
}

func VdsAlertHandler(ev *fsnotify.FileEvent) error {
	topicname, _ := mconfig.Conf.String("kafka", "VdsAlertTopicName")
	return AlertHandler(ev.Name, topicname, "xdr")
}

func DpiHandle(ev *fsnotify.FileEvent) error {
	defer func() {
		if err := recover(); err != nil {
			mlog.Error(err)
		}
	}()
	var filename = ev.Name
	mlog.Alert("time1=", time.Now().Unix())
	mlog.Debug(fmt.Println("Create file:", filename))

	//check file suffix
	if ok := CheckSuffix(filename, []string{"xdr", "XDR"}...); !ok {
		panic(fmt.Sprintf("file: %s suffix error!", filename))
	}

	//read file
	content, err := ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("read file %s error:%s", filename, err.Error()))
	}
	mlog.Alert("time2=", time.Now().Unix())
	//XDR==>pre object
	datalist, err := xdrParse.ParseXdr(content)
	if err != nil {
		panic(fmt.Sprintf("parse file %s Xdr error:%s", filename, err.Error()))
	}
	mlog.Alert("time3=", time.Now().Unix())
	//file to ceph
	count, _ := dealForeachXdrs(datalist)
	mlog.Debug("file", filename, "has xdr", count)
	return nil
}

func dealForeachXdrs(xdrs []*xdrParse.DpiXdr) (int, error) {
	var count = 0
	for _, xdr := range xdrs {
		//check type
		xtype := xdr.CheckType()
		switch xtype {
		case xdrParse.XdrType:
			dealXdrCommon(xdr)
		case xdrParse.XdrHttpType:
			go dealHttpFileXdr(xdr)
		case xdrParse.XdrFileType:
			go dealHttpFileXdr(xdr)
		default:
			mlog.Error("CheckType error! return ", xtype)
		}
		count++
	}
	return count, nil
}

/*
 * deal xdr which is http or file,need save to ceph
 */
func dealHttpFileXdr(xdr *xdrParse.DpiXdr) error {
	//1.save to ceph
	if err := saveToCephPerXdr(xdr); err != nil {
		return err
	}
	//2.pre obj --> backend obj
	backend := PerTransToBackendObj(xdr)

	//3.push to kafka
	if err := DoPushTopic(backend); err != nil {
		return err
	}
	return nil
}

/*
 * deal xdr not http and file,not need save to ceph
 */
func dealXdrCommon(xdr *xdrParse.DpiXdr) error {
	//0.save to ceph
	//no need
	//1.pre obj --> backend obj
	backend := PerTransToBackendObj(xdr)
	//2.push to kafka
	if err := DoPushTopic(backend); err != nil {
		return err
	}
	return nil
}

func AlertHandler(fileName string, topicName string, suffix string) error {
	defer func() {
		if err := recover(); err != nil {
			mlog.Error(err)
		}
	}()
	mlog.Debug(fmt.Println("Create file:", fileName))
	mlog.Alert("time1=", time.Now())
	//check file suffix
	if ok := CheckSuffix(fileName, suffix); !ok {
		panic(fmt.Sprintf("file: %s suffix error!", fileName))
	}

	//push kafka

	pushkafkaFunc := func(line string) error {
		data := &DataType{
			topicName: topicName,
			handlePre: func(data []byte) ([]byte, error) {
				return data, nil
			},
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

func ReadFile(FileName string) ([]byte, error) {
	fileContent, err := ioutil.ReadFile(FileName)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
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
