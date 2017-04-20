#!/bin/bash

src_path=/mnt/ftp/DPI
dst_path=/mnt/DPI/XDR
log_path=/tmp/DPI/log
temp_path=/tmp/DPI

function main() {
/usr/local/bin/inotifywait -mrq --timefmt '%d-%m-%y-%H:%M:%S' --format '%T %w %f' -e create,moved_to $src_path | while read  time dir file event
do
	echo $time'_'$dir'_'$file >> $log_path/decompressed.log 
	
	if [[ "${file##*.}" = "z4" ]] then
	dstfile=${file%.*}'.xdr'
		lz4 -f -d $src_path/$file $temp_path/${dstfile}
		mv $temp_path/${dstfile}   $dst_path/${dstfile}
	fi 
	
	rm $src_path/$file
done	
}

main $@



