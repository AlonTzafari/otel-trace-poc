#!/bin/bash
docker build -f dockerfile-go --build-arg SERVICE=$1 -t $1 .