#!/bin/bash

echo "=== Build Time-autoscaler for OpenShift v3 ==="

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./image/time-scaler-engine ./go/*.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o ./release/time-scaler-engine ./go/*.go


if [ $? = 0 ]
then
  echo "[INFO] - Compile success ==="
  echo "[INFO] - Building image ... ==="
	docker build --tag=time-scaler-engine -f ./image/Dockerfile ./image
fi
