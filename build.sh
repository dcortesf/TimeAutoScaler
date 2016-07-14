#!/bin/bash

echo "=== Build Time-autoscaler for OpenShift v3 ==="

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./image/time-autoscaler ./go/*.go

if [ $? = 0 ]
then
  echo "[INFO] - Compile success ==="
  echo "[INFO] - Building image ... ==="
	docker build --tag=time-autoscaler -f ./image/Dockerfile ./image
fi
