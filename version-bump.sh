#!/bin/bash

cd /tmp/temp

mvn versions:use-latest-versions -Dincludes=$1

STATUS=$?
if [ $STATUS -eq 0 ]; then
   echo "Build Successful"
else
   echo "Build Failed"
fi
