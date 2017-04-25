#!/bin/bash

#mv source to go src
#go get package
#go build
#./preprocess
#./watchdir.sh

main() {
			
	$(go_get_package)
	go build
	nohup ./preprocess &
	./watchdir.sh &
}

go_get_package() {
	go get github.com/howeyc/fsnotify
	go get github.com/astaxie/beego
	go get github.com/optiopay/kafka
}

check_key() {
	kv = $1
	stand_key = $2
	key = `echo "${kv}" | cut -d \= -f 1`
	if [[ "${stand_key}" != "${key}" ]];then
		return false
		
	fi
}

get_value() {
	kv = $1	
	key = $2
	if [ ! $(check_key "${kv}" "${stand_key}") ];then
		return 1	
	fi
	
}

main $@
