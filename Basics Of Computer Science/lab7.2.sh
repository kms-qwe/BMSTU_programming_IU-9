#!/bin/bash

if ! [[ $# -eq 1 && -d "$1" ]]; then
    echo "Invalid arguments"
    exit 1
fi

export cnt=0

rec() {
    if [[ -d "$1" ]]; then
        for name in "$1"/*; do
            rec "$name"
        done
    elif [[ -f "$1" && "$1" =~ \.(c|h)$ ]]; then
    current_cnt=$(cat "$1" | sed '/^\s*$/d' | grep -c '')
    cnt=$(($cnt + $current_cnt))
    fi
}

rec "$1"
echo "$cnt"