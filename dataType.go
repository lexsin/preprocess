package main

//"preprocess/modules/xdrParse"

type DataType struct {
	topicName        string
	partition        int
	origiData        []byte
	handlePre        func(data []byte) ([]byte, error)
	partitionHandler func(data *BackendObj) int
}

func (this *DataType) TopicName() string {
	return this.topicName
}

func (this *DataType) OriginalData() []byte {
	return this.origiData
}

func (this *DataType) Partition() int {
	return this.partition
}

func (this *DataType) PreProcessData(data []byte) ([]byte, error) {
	return this.handlePre(data)
}

func XdrPreHandle(data []byte) ([]byte, error) {
	//TODO
	return data, nil
}

func XdrHttpPreHandle(data []byte) ([]byte, error) {
	//TODO
	return data, nil
}

func XdrFilePreHandle(data []byte) ([]byte, error) {
	//TODO
	return data, nil
}
