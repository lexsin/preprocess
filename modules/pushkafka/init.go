package pushkafka

import (
	"errors"
	"fmt"
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"
	"time"

	"github.com/optiopay/kafka"
)

//var kafkaAddrs = []string{"10.80.6.9:9092", "10.80.6.9:9093"}
var kafkaAddrs []string
var Broker *kafka.Broker
var WriterMap map[string]chan dataInfo
var waitTimeOut = 600 * time.Second

func init() {
	//init module first
	moduleInit()
	//init topic and write
	CreateTopicWriter("xdr")
	CreateTopicWriter("xdrHttp")
	CreateTopicWriter("xdrFile")

	CreateTopicWriter("waf-alert")
	CreateTopicWriter("ids-alert")
	CreateTopicWriter("vds-alert")
	//init debug
	debugInit()
}

func debugInit() {
	//debugOpen, _ := mconfig.Conf.Int("debug", "pushKafka")
}

func moduleInit() {
	//init map
	WriterMap = make(map[string]chan dataInfo, 0)
	//init broker
	conf := kafka.NewBrokerConf("apt-pre")
	conf.AllowTopicCreation = true
	// connect to kafka cluster
	var err error
	host, err := mconfig.Conf.String("kafka", "KafkaHost")
	if err != nil {
		mlog.Error("[kafka]KafkaHost not configure")
		panic(ErrNotConfigErr)
	}
	addr1 := fmt.Sprintf("%s:9092", host)
	addr2 := fmt.Sprintf("%s:9093", host)
	mlog.Debug("kafka addr1=", addr1)
	mlog.Debug("kafka addr2=", addr2)
	kafkaAddrs = make([]string, 0)
	kafkaAddrs = append(kafkaAddrs, addr1)
	kafkaAddrs = append(kafkaAddrs, addr2)

	Broker, err = kafka.Dial(kafkaAddrs, conf)
	if err != nil {
		mlog.Error(err)
		panic(fmt.Sprintf("kafka.Dial %s", err.Error()))
	}
	mlog.Info("connect kafka success!")
	//defer broker.Close()
}

var ErrNotConfigErr = errors.New("not configure")
