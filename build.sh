#!/bin/bash

pwd=`pwd`
appName=`basename $pwd`
# 编译的分支
branch="master"
# 当前git分支名
localBranch=`git rev-parse --abbrev-ref HEAD`
# 编译后的输出文件名称，赋值当前项目文件名
targetFile="$appName"
# 编译的包名
buildPkg="main.go"
# 编译结果
buildResult=""

pidFile="pid.txt"

today=$(date "+%Y_%m_%d")
# 日志存放路径
logDir="/root/logs/${appName}"
info_log="${logDir}/info.${today}.log"
error_log="${logDir}/error.${today}.log"
path=targetFile

echo "当前分支: ${localBranch}"

if [ -f $pidFile ]; then
  # pid=$(cat $pidFile)
  pid=$(<$pidFile)
  echo "Prepare to kill the process: ${pid}"
  kill -9 $pid
fi

echo "项目路径: ${pwd}"
echo "start clean log file"

if [ ! -d "$logDir" ]; then
  mkdir "$logDir"
fi

# 半年前的日期
m=$(date -d "-6 months" "+%s")

echo "半年前的日期: $(date -d @${m} "+%Y-%m-%d")"

index=1
f=`ls ${logDir} -1 -c`

# 清理旧日志文件
for name in $f
do
  # echo "日志${index}：$name"
  dateStr=$(echo ${name} | grep -Eo "[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}")
  # echo $dateStr
  # 判断是否有效
  if date -d ${dateStr} > /dev/null 2>&1; then
    t1=`date -d "$dateStr" +%s`

    if [ $t1 -lt $m ]; then
      echo ">>> delete file: ${name}"
      rm -rf "${logDir}/${name}"
    fi

  fi
  let index++
done

echo "complete the clean"

# if [ -n "$1" ]; then
#   branch="$1"
#   echo "Switch branch to ${branch}"
# else
#   echo "Building Branch: ${branch}"
# fi

# git checkout "$branch"
# git pull


flags="-X '${path}.version=v1.0' -X '${path}.go_version=$(go version)' -X '${path}.build_time=$(date "+%Y.%m.%d %H:%M:%S")' -X '${path}.build_user=$(id -u -n)' -X '${path}.commit_id=$(git rev-parse --short HEAD)'"
buildResult=`go build -ldflags "$flags" -o "${targetFile}" "$buildPkg"`
echo "$flags"

if [ -z "$buildResult" ]; then
  chmod 773 ${targetFile}
  echo "result: $buildResult"
  echo "build success, filename: ${targetFile}"
else
  echo "build error $buildResult"
  exit
fi


# if [ ! -f "$info_log" ]; then
#   touch "$info_log"
# fi

# if [ ! -f "$error_log" ]; then
#   touch "$error_log"
# fi

# nohup "./${targetFile}" 1>"${info_log}" 2>"${error_log}" & echo $! > "$pidFile"
nohup "./${targetFile}" 1>/dev/null 2>&1 & echo $! > "$pidFile"

echo "------new pid: $(<$pidFile)"

echo "deploy success..."
