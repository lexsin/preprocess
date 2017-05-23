package pushkafka

import (
	"errors"
	"fmt"
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"
	"time"

	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
)

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
		mlog.Error("topic:", topic, "not exist")
		return errors.New("topic not exist")
	}
	data, err := info.PreProcessData(info.OriginalData())
	if err != nil {
		mlog.Error(fmt.Sprintf("PreProcessData error:%s", err.Error()))
		return errors.New(fmt.Sprintf("PreProcessData error:%s", err.Error()))
	}
	datainfo := dataInfo{
		data:      data,
		partition: info.Partition(),
	}
	for {
		select {
		case writer <- datainfo:
			mlog.Debug("push topic(", topic, ") partation(", datainfo.partition, ")success!")
			return nil
		default:
			mlog.Alert("topic:", topic, "wirte wait...")
			<-time.After(1 * time.Second)
		}
	}
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
		for {
			select {
			case data := <-ch:
				//mlog.Debug("get topic data:", string(data.data))
				producer := Broker.Producer(kafka.NewProducerConf())
				msg := &proto.Message{Value: data.data}

				for i := 0; i < 3; i++ {
					if _, err := producer.Produce(topicName, int32(data.partition), msg); err == nil {
						break
					} else {
						if err == proto.ErrRequestTimeout {
							if i == 2 {
								mlog.Error(fmt.Sprintf("Write topic %s paration %d error:%s,(%d)times",
									topicName, data.partition, err.Error()), 3)
							}
							continue
						}
						mlog.Error(fmt.Sprintf("Write topic %s paration %d error:%s,",
							topicName, data.partition, err.Error()))
						break
					}
				}
			case <-time.After(waitTimeOut):
				mlog.Debug("writer(", topicName, ") wait ", time.Duration(waitTimeOut).Seconds(), "s...")
			}
		}
	}()
	//save write chan name
	WriterMap[topicName] = ch

	mlog.Info("create topic:", topicName)
	return nil
}
