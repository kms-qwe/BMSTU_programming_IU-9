#!/bin/bash


if ! [[ $# -eq 2 && -f $1 && $2 =~ ^[0-9]+$ ]]; then
    echo "Invalid arguments"
    exit 1
fi
curTime=$(date +'%Y-%m-%d-%H-%M-%S-%3N')
outFile="out_lab7_$curTime.log"
# echo $outFile
errorFile="error_lab7_$curTime.log"
path=$1
delay=$2
counter=0
cnt=1
bash $path 1>>$outFile 2>>$errorFile &
curPID=$!
echo "C: $cnt RN: $counter curPID =$curPID"
cnt=$(($cnt + 1))
while true; do
    if [[ $cnt -eq 11 ]]; then
        echo "THE END"
        exit 0
    fi
    if ! ps -p $curPID 1>/dev/null; then
        bash $path 1>>$outFile 2>>$errorFile &
        curPID=$!
        counter=$(($counter + 1))
        echo "C: $cnt RN: $counter: curPID = $curPID"
    else
        echo "C: $cnt RN: $counter: The program in progress"
    fi
    sleep $(($delay))
    cnt=$(($cnt + 1))
done

