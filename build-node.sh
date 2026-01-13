#!/bin/bash
docker build -f dockerfile-node --build-arg SERVICE=$1 -t $1 .