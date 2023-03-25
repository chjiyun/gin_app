#!/bin/bash

pwd=`pwd`
# appName=`basename $pwd`
# 从配置文件取name，默认是项目名，-r 代表过滤掉字符串的双引号
appName=`cat ./config/config.yml | yq -r .name`
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
# 应用启动端口
port="8000"

today=$(date "+%Y_%m_%d")
# 日志存放路径
logDir="/root/logs/${appName}"
info_log="${logDir}/info.${today}.log"
error_log="${logDir}/error.${today}.log"
path="main"

echo "----------------------"
echo "app name: ${appName}"
echo "app dir: ${pwd}"
echo "current branch: ${localBranch}"


echo "start clean log file"

if [ ! -d "$logDir" ]; then
# 创建多级目录
  mkdir -p "$logDir"
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


flags="-X '${path}.AppVersion=v1.0' -X '${path}.GoVersion=$(go version | awk '{print $3 " " $4}')' -X '${path}.BuildTime=$(date "+%Y.%m.%d %H:%M:%S")' -X '${path}.BuildUser=$(id -u -n)' -X '${path}.CommitId=$(git rev-parse --short HEAD)'"
buildResult=`go build -ldflags "$flags" -o "${targetFile}" "$buildPkg"`

# 编译成功才能杀旧进程
if [ $? -eq 0 ]; then 
  chmod 773 ${targetFile}
  echo "build success, filename: ${targetFile}"

  pid=`ps -ef |grep $targetFile | grep -v grep|awk '{print $2}'`
  echo "current pid is $pid"
  if [ "$pid" != "" ]; then
    echo "Prepare to kill the process: ${pid}"
    kill -9 $pid
    sleep 1
  fi
else
  echo "build error $buildResult"
  exit
fi

# nohup "./${targetFile}" 1>"${info_log}" 2>"${error_log}" & echo $! > "$pidFile"
nohup "./${targetFile}" 1>/dev/null 2>&1 &

echo "starting..."
# 监听端口是否启动
while :
do
    running=`lsof -i:$port | wc -l`
    if [ $running -gt "0" ]; then
        echo "server is running on $port"
        break
    fi
    sleep 1
done

# 打印版本信息
"./${targetFile}" -v
