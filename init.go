package main

import (
	//"fmt"
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"

	//"github.com/julienschmidt/httprouter"
	//"preprocess/modules/mlog"
)

var TopicMap map[int]DataType

var AgentNum int

func init() {
	//init log
	mlog.SetLogger("file", `{"filename":"logs/server.log"}`)
	mlog.SetLogger("console", "")
	mlog.SetLogLevel(mlog.LevelDebug)

	//init topic object
	topicInit()

	//init waf server
	wafServInit()
}

func topicInit() {
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
