#!/bin/bash
echo "start..."
basepath=$(cd `dirname $0`; pwd)
cd $basepath
sleep 1
nohup ./getip >run.log 2>&1 &
echo "started..."