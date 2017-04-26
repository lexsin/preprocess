#!/bin/bash

#mkdir 
#mv source to go src
#go get package
#go build
#./preprocess
#./watchdir.sh

RealRootDir="/root"
LnsRootDir="/mnt"
DpiZ4Dir="${RealRootDir}/ftp/Dpi"
DpiXdrDir="${RealRootDir}/Dpi/xdr"
IdsAlertDir="${RealRootDir}/ftp/alert/ids"
VdsAlertDir="${RealRootDir}/ftp/alert/vds"
temp_path="/tmp/DPI"


main() {
	$(create_dir)	
	$(go_get_package)
	go build
	nohup ./preprocess &
	./preWatchdir.sh &
}

create_dir() {
	mkdir -p  ${DpiZ4Dir}
	mkdir -p  ${DpiXdrDir}
	mkdir -p  ${IdsAlertDir}
	mkdir -p  ${VdsAlertDir}
	mkdir -p  ${temp_path}
	ln -s "${RealRootDir}/ftp" "${LnsRootDir}/ftp"
	ln -s "${RealRootDir}/Dpi" "${LnsRootDir}/Dpi"
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
