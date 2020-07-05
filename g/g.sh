#!/bin/bash

#mac os 修改getopt
OS=$(uname -s)
if [[ ${OS} -eq "Darwin" ]]; then
    export FLAGS_GETOPT_CMD="$(brew --prefix gnu-getopt)/bin/getopt"
fi

. shflags/shflags

#参数列表
DEFINE_boolean 'new' false 'new branch' 'n'

newbranch () {
    #获取当前分支
    nowbranch=$(git symbolic-ref --short -q HEAD)
    if [ $nowbranch != 'online' ] ; then
        echo -e '\033[31m不在online分支\033[m'
        read -p y:继续 continue
        if [[ $continue != 'y' ]]; then
            return -1
        fi
    fi
    git pull
    git checkout -b $(date +%Y-%m-%d-)$(openssl rand -hex 2) $(git last) && git branch --edit-description
}

main () {
    if [[ ${FLAGS_new} -eq ${FLAGS_TRUE} ]]; then
        newbranch
    fi
}

FLAGS "$@" || exit $?
eval set -- "${FLAGS_ARGV}"

main "$@"