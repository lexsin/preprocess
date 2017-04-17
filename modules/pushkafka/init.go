package pushkafka

import (
	//"preprocess/modules/mconfig"
	"fmt"
	"preprocess/modules/mlog"

	"github.com/optiopay/kafka"
)

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
	Broker, err = kafka.Dial(kafkaAddrs, conf)
	if err != nil {
		mlog.Error(err)
		panic(fmt.Sprintf("kafka.Dial %s", err.Error()))
	}
	mlog.Info("connect kafka success!")
	//defer broker.Close()
}
