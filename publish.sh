#!/bin/bash
proj_path=$GOPATH/src/github.com/hzxiao/ocss-srv
app_path=/home/jks/app
proj_name=ocss

cd $proj_path
git checkout .
git pull

sudo chmod +x ./pkg.sh
./pkg.sh

sudo ./$proj_name'd' stop

rm -rf $app_path/$proj_name

cp -r build/$proj_name $app_path/

cd /etc/rc.d/init.d/
sudo ./$proj_name'd' start