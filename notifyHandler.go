package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	//"preprocess/modules/mconfig"
	"preprocess/modules/mlog"
	"preprocess/modules/pushkafka"
	"preprocess/modules/xdrParse"

	//"time"

	"github.com/howeyc/fsnotify"
)

func RunNotify(dir string, handle func(filename string) error) {
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
			mlog.Debug("event:", ev)

		case ev := <-Watcher.Event:
			if ev.IsCreate() {
				go handle(ev.Name)
			}

		case err := <-Watcher.Error:
			mlog.Error(err)
		}
	}

	/* ... do stuff ... */
	Watcher.Close()
}
func IdsAlertHandler(filename string) error {
	defer func() {
		err := recover()
		dealFileByErr_idsAlert(err, filename)
		if err != nil {
			mlog.Error(err)
		}
	}()

	checkForm := func(line string) error {
		alert := IdsAlert{}
		if err := json.Unmarshal([]byte(line), &alert); err != nil {
			return err
		}
		return nil
	}

	return AlertHandler(filename, IdsAlertTopic, "xdr", checkForm)
}

func VdsAlertHandler(filename string) error {
	defer func() {
		err := recover()
		dealFileByErr_vdsAlert(err, filename)
		if err != nil {
			mlog.Error(err)
		}
	}()

	checkForm := func(line string) error {
		alert := VdsFullAlert{}
		if err := json.Unmarshal([]byte(line), &alert); err != nil {
			return err
		}
		return nil
	}

	return AlertHandler(filename, VdsAlertTopic, "xdr", checkForm)
}

func DpiHandle(filename string) error {
	defer func() {
		err := recover()
		dealFileByErr_dpi(err, filename)
		if err != nil {
			mlog.Error(err)
		}
	}()

	mlog.Debug(fmt.Println("Create file:", filename))

	//check file suffix
	if ok := CheckSuffix(filename, []string{"xdr", "XDR"}...); !ok {
		mlog.Error("file: %s suffix error!", filename)
		panic(ErrSuffixErr)
	}

	//read file
	content, err := ReadFile(filename)
	if err != nil {
		mlog.Error("read file %s error:%s", filename, err.Error())
		panic(ErrReadFileErr)
	}

	//XDR==>pre object
	datalist, err := xdrParse.ParseXdr(content)
	if err != nil {
		mlog.Error("parse file %s Xdr error:%s", filename, err.Error())
		panic(ErrXdrParseErr)
	}

	//deal objects foreach
	count, success := dealForeachXdrs(datalist)
	if success < count {
		mlog.Error("file push kafka error")
		mlog.Error("file:", filename, "count=", count, "success=", success)
		panic(ErrPushKafkaErr)
	}
	return nil
}

func deleteFileSwitch(filename string, reserv bool, dir string) error {
	if reserv {
		if dir != "" {
			RenameFile(filename, dir)
		}
		return nil
	} else {
		DeleteFile(filename)
	}
	return nil
}

/*
 * 逐条处理，1.转换成backend 2.入kafka，
 */
func dealForeachXdrs(xdrs []*xdrParse.DpiXdr) (count int, success int) {
	var success1 = 0
	var fail1 = 0
	count = 0
	stops := make(chan error, len(xdrs))
	realsize := 0
	for _, xdr := range xdrs {
		count++
		//check type
		xtype := xdr.CheckType()
		switch xtype {
		case xdrParse.XdrType:
			if err := dealXdrCommon(xdr); err != nil {
				fail1++
				break
			}
			success1++
		case xdrParse.XdrHttpType:
			go dealHttpFileXdr(xdr, stops)
			realsize++
		case xdrParse.XdrFileType:
			go dealHttpFileXdr(xdr, stops)
			realsize++
		default:
			mlog.Error("CheckType error! return ", xtype)
		}
	}
	var success2 = 0
	var fail2 = 0
	for i := 0; i < realsize; i++ {
		if err := <-stops; err != nil {
			fail2++
		} else {
			success2++
		}
	}
	success = success1 + success2
	mlog.Debug("xdr count=", count, "recieve success=", success1+success2)
	return
}

/*
 * deal xdr which is http or file,need save to ceph
 */
func dealHttpFileXdr(xdr *xdrParse.DpiXdr, stop chan error) error {
	//1.save to ceph
	if err := saveToCephPerXdr(xdr); err != nil {
		stop <- err
		return err
	}
	//2.pre obj --> backend obj
	backend := PerTransToBackendObj(xdr)

	//3.push to kafka
	if err := DoPushTopic(backend); err != nil {
		stop <- err
		return err
	}
	stop <- nil
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

type perAlertFuncs struct {
	pushkafka func(line string) error
	checkForm func(line string) error
}

func AlertHandler(fileName string, topicName string, suffix string, checkForm func(line string) error) error {
	mlog.Debug(fmt.Println("Create file:", fileName))
	//check file suffix
	if ok := CheckSuffix(fileName, suffix); !ok {
		mlog.Debug(fmt.Sprintf("file: %s suffix error!", fileName))
		panic(ErrSuffixErr)
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
	inter := perAlertFuncs{
		pushkafka: pushkafkaFunc,
		checkForm: checkForm,
	}
	if _, err := DealFilePerline(fileName, inter); err != nil {
		return err
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

func dealFileByErr_dpi(errer interface{}, filename string) error {
	if errer == nil {
		deleteFileSwitch(filename, DoDelDpi, "")
	} else {
		if isXdrPkgErr(errer.(error)) {
			deleteFileSwitch(filename, DoDelIlegalDpi, DpiWatchDir+"/illegal/")
		} else {
			deleteFileSwitch(filename, true, DpiWatchDir+"/bak/")
		}
	}
	return nil
}
func dealFileByErr_vdsAlert(errer interface{}, filename string) error {
	if errer == nil {
		deleteFileSwitch(filename, DoDeleteVdsAlert, "")
	} else {
		deleteFileSwitch(filename, true, VdsAlertWatchDir+"/bak/")
	}
	return nil
}
func dealFileByErr_idsAlert(errer interface{}, filename string) error {
	if errer == nil {
		deleteFileSwitch(filename, DoDeleteIdsAlert, "")
	} else {
		deleteFileSwitch(filename, true, IdsAlertWatchDir+"/bak/")
	}
	return nil
}
