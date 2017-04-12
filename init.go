package main

import (
	//"fmt"
	"preprocess/modules/mconfig"
	//"preprocess/modules/mlog"

	//"github.com/howeyc/fsnotify"
)

var TopicMap map[int]DataType

//var Watcher *fsnotify.Watcher
var AgentNum int

func init() {
	/*
		//init fsnotify
		var err error
		Watcher, err = fsnotify.NewWatcher()
		if err != nil {
			mlog.Error(err)
			panic(fmt.Sprintf("fsnotify.NewWatcher() error:%s", err))
		}
	*/

	//init topic object
	TopicMap = make(map[int]DataType, 0)

	XdrTopic, _ := mconfig.Conf.String("kafka", "XdrTopicName")
	HttpTopic, _ := mconfig.Conf.String("kafka", "HttpTopicName")
	FileTopic, _ := mconfig.Conf.String("kafka", "FileTopicName")

	TopicMap[XdrType] = DataType{
		topicName: XdrTopic,
		handlePre: XdrPreHandle,
	}
	TopicMap[XdrHttpType] = DataType{
		topicName: HttpTopic,
		handlePre: XdrHttpPreHandle,
	}
	TopicMap[XdrFileType] = DataType{
		topicName: FileTopic,
		handlePre: XdrFilePreHandle,
	}
}

const (
	XdrType = iota
	XdrHttpType
	XdrFileType
)
