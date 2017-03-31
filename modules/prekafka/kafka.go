package prekafka

import (
	"fmt"
	"preprocess/modules/mlog"

	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
)

var kafkaAddrs = []string{"localhost:9092", "localhost:9093"}
var Broker *kafka.Broker

func init() {
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

type PushTopicer interface {
	OriginalData() []byte
	TopicName() string
	PreProcessData(data []byte) ([]byte, error)
	Partition() int
}

func PushTopic(inter PushTopicer) error {
	origidata := inter.OriginalData()
	data, err := inter.PreProcessData(origidata)
	if err != nil {
		return err
	}
	topic := inter.TopicName()
	partition := inter.Partition()
	return produce(topic, partition, data)
}

func produce(topic string, partition int, data []byte) error {
	producer := Broker.Producer(kafka.NewProducerConf())
	msg := &proto.Message{Value: data}
	if _, err := producer.Produce(topic, int32(partition), msg); err != nil {
		return err
	}
	return nil
}
