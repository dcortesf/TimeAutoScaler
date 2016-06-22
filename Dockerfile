FROM golang:1.7-onbuild
COPY . /go/src/mytime
RUN go get -d -v
RUN go install -v
