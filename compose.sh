#!/bin/bash
option=""
if [ "$1" = "prod" ]; then
    fname="-f docker-compose.prod.yaml"
fi

docker-compose $fname down  --remove-orphans; docker-compose $fname build && docker-compose $fname up -d --scale hnfaves=4
