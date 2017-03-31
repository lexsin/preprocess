package pushkafka

import (
	"errors"
	"fmt"
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"

	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
)

var kafkaAddrs = []string{"localhost:9092", "localhost:9093"}
var Broker *kafka.Broker
var WriterMap map[string]chan dataInfo

type PushKafkaer interface {
	TopicName() string
	Partition() int
	OriginalData() []byte
	PreProcessData(data []byte) ([]byte, error)
}

func PushKafka(info PushKafkaer) error {
	topic := info.TopicName()
	writer, ok := WriterMap[topic]
	if !ok {
		return errors.New("topic not exist")
	}
	data, err := info.PreProcessData(info.OriginalData())
	if err != nil {
		return errors.New(fmt.Sprintf("PreProcessData error:%s", err.Error()))
	}
	datainfo := dataInfo{
		data:      data,
		partition: info.Partition(),
	}
	writer <- datainfo
	return nil
}

type dataInfo struct {
	data      []byte
	partition int
}

func CreateTopicWriter(topicName string) error {
	//creat topic

	//run topic write
	chansize, _ := mconfig.Conf.Int("kafka", "ChannelBuffer")
	if chansize < 1 {
		mlog.Error("ChannelBuffer < 1")
		chansize = 1
	}
	ch := make(chan dataInfo, chansize)
	go func() {
		select {
		case data := <-ch:
			producer := Broker.Producer(kafka.NewProducerConf())
			msg := &proto.Message{Value: data.data}
			if _, err := producer.Produce(topicName, int32(data.partition), msg); err != nil {
				mlog.Error(fmt.Sprintf("Write topic %s paration %d error:%s",
					topicName, data.partition, err.Error()))
			}
		}
	}()
	//save write chan name
	WriterMap[topicName] = ch
	return nil
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

func init() {
	//init module first
	moduleInit()
	//init topic and write
	CreateTopicWriter("xdr")
	CreateTopicWriter("xdrHttp")
	CreateTopicWriter("xdrFile")
}
