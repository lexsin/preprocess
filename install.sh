#!/bin/bash

#mkdir 
#mv source to go src
#go get package
#go build
#./preprocess
#./watchdir.sh


main() {
	$(create_dir)	
	$(go_get_package)
	go build
	nohup ./preprocess &
	./watchdir.sh &
}

create_dir() {
	mkdir -p /home/ftp/DPI
	mkdir -p /home
}

go_get_package() {
	go get github.com/howeyc/fsnotify
	go get github.com/astaxie/beego
	go get github.com/optiopay/kafka
}

check_key() {
	local kv = $1
	local stand_key = $2
	local key = `echo "${kv}" | cut -d \= -f 1`
	if [[ "${stand_key}" != "${key}" ]];then
		return false
		
	fi
}

get_value() {
	local kv = $1	
	local key = $2
	if [ ! $(check_key "${kv}" "${stand_key}") ];then
		return 1	
	fi
	
}

main $@
