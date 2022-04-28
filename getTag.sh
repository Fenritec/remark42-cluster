#!/bin/bash
TAG=`git describe --exact-match --tags --abbrev=0 $(git log -n1 --pretty='%h') 2> /dev/null`
if [ $? -gt 0 ]
then
  TAG="devel"
fi
echo $TAG
