#!/bin/sh

for i in `seq 0 10000`
do
	curl "http://localhost:6873/pub?id=1a60e483d304ff356a991e490e58ecad&msg=$i"
	echo
	curl "http://localhost:6873/pub?id=b963511e400577bc8dfe5bcc4b6e677d&msg=$i"
	echo

	sleep 1
done
