#!/bin/bash
proj_path=$GOPATH/src/github.com/hzxiao/ocss-srv
app_path=/home/jks/app
proj_name=ocss

cd $proj_path
git checkout .
git pull

sudo chmod +x ./pkg.sh
./pkg.sh

cp build/$proj_name'.zip' $app_path
cd $app_path

unzip $proj_name'.zip'

sudo /etc/rc.d/init.d/$proj_name'd' restart
