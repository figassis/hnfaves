#!/bin/bash
clear
tag=`./version.sh`
docker build --rm -t figassis/hnfaves:$tag . && docker push figassis/hnfaves:$tag