package main

type TopicType struct {
	topicName string
	partition int
	origiData []byte
	handlePre func(data []byte) ([]byte, error)
}

func (this *TopicType) TopicName() string {
	return this.topicName
}

func (this *TopicType) OriginalData() []byte {
	return this.origiData
}

func (this *TopicType) Partition() int {
	return this.partition
}

func (this *TopicType) PreProcessData(data []byte) ([]byte, error) {
	return this.handlePre(data)
}

func XdrPreHandle(data []byte) ([]byte, error) {
	//TODO
	return data, nil
}

func XdrHttpPreHandle(data []byte) ([]byte, error) {
	//TODO
	return nil, nil
}

func XdrFilePreHandle(data []byte) ([]byte, error) {
	//TODO
	return nil, nil
}
