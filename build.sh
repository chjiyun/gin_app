#!/bin/bash

pwd=`pwd`
projectName=`basename $pwd`
# 编译的分支
branch="master"
# 编译后的输出文件名称，赋值当前项目文件名
targetFile="$projectName"
# 编译的包名
buildPkg="main.go"
# 编译结果
buildResult=""

pidFile="pid.txt"

today="$(date "+%Y_%m_%d")"
# 日志存放路径
logDir="/root/logs/${projectName}"
info_log="${logDir}/info_${today}.log"
error_log="${logDir}/error_${today}.log"


echo "项目路径：${pwd}"

if [ -n "$1" ]; then
  branch="$1"
  echo "Switch branch to ${branch}"
else
  echo "Building Branch: ${branch}"
fi

git checkout "$branch"
git pull

if [ -f $pidFile ]; then
  # pid=$(cat $pidFile)
  pid=$(<$pidFile)
  echo "Prepare to kill the process: ${pid}"
  kill -9 $pid
fi

buildResult=`go build -o "${targetFile}" "$buildPkg"`

if [ -z "$buildResult" ]; then
chmod 773 ${targetFile}
echo "build success, filename：${targetFile}"
else
echo "build error $buildResult"
exit
fi

if [ ! -d "$logDir" ]; then
  mkdir "$logDir"
fi

if [ ! -f "$info_log" ]; then
  touch "$info_log"
fi

if [ ! -f "$error_log" ]; then
  touch "$error_log"
fi

nohup "./${targetFile}" 1>"${info_log}" 2>"${error_log}" & echo $! > "$pidFile"

echo "------new pid: $(<$pidFile)"

echo "deploy success..."
