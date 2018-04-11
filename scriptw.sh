#!/usr/bin/env bash

if [ "$#" -ne 1 ]
then
    echo -e "   Please add a argument"
    echo -e "   Usage: ./scriptw [arg]\n"
    echo -e "\t pack\t packing app"
    echo -e "\t clean\t clean dist and build dir"
fi
if [ "$1" == 'pack' ]
then
    if [ -e ./build ];then "$0" clean ;fi
    mkdir build
    go build -o bin
    tar -zcvf sorry-gen.tar.gz bin ./resources
    mv sorry-gen.tar.gz build/
    rm bin
fi
if [ "$1" == 'clean' ]
then
    rm -rf ./build ./dist
fi