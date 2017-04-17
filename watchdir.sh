#!/bin/bash

src_path=/mnt/ftp/DPI
dst_path=/mnt/DPI/XDR
log_path=/tmp/DPI/log

function main() {
/usr/local/bin/inotifywait -mrq --timefmt '%d-%m-%y-%H:%M:%S' --format '%T %w %f' -e create,moved_to $src_path | while read  time dir file event
do
	echo $time'_'$dir'_'$file >> $log_path/decompressed.log 
	lz4 -f -d $src_path/$file $dst_path/${file%.*}'.xdr'
done	
}

main $@



