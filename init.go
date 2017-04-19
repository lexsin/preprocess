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

const (
	XdrType = iota
	XdrHttpType
	XdrFileType
)

var DpiWatchDir string
var VdsAlertWatchDir string
var IdsAlertWatchDir string

var DoDelDpi bool
var DoDelIlegalDpi bool
var DoDeleteVdsAlert bool
var DoDelIlegalVdsAlert bool
var DoDeleteIdsAlert bool
var DoDelIlegalIdsAlert bool

var IdsAlertTopic string
var VdsAlertTopic string

func init() {
	//init log
	mlog.SetLogger("file", `{"filename":"logs/server.log"}`)
	mlog.SetLogger("console", "")
	mlog.SetLogLevel(mlog.LevelDebug)

	//init topic object
	topicInit()

	//init waf server
	wafServInit()

	//read config
	initVariate()
}
func initVariate() {
	var err error
	GetConfBool(&DoDelDpi, "resvfile", "resvLegalDpi", false)
	GetConfBool(&DoDelIlegalDpi, "resvfile", "resvIlegalDpi", true)
	GetConfBool(&DoDeleteVdsAlert, "resvfile", "resvLegalAltVds", false)
	GetConfBool(&DoDelIlegalVdsAlert, "resvfile", "resvIlegalAltVds", true)
	GetConfBool(&DoDeleteIdsAlert, "resvfile", "resvLegalAltIds", false)
	GetConfBool(&DoDelIlegalIdsAlert, "resvfile", "resvIlegalAltIds", true)

	IdsAlertTopic, err = mconfig.Conf.String("kafka", "IdsAlertTopicName")
	if err != nil {
		panic("[kafka]IdsAlertTopicName not config")
	}
	VdsAlertTopic, err = mconfig.Conf.String("kafka", "VdsAlertTopicName")
	if err != nil {
		panic("[kafka]VdsAlertTopicName not config")
	}
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
