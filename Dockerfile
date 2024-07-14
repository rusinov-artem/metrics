FROM golang:1.22

RUN apt update ; apt install -y git make jq curl vim htop tmux ncat iputils-ping net-tools;

RUN go install github.com/go-delve/delve/cmd/dlv@latest &&\
    go install github.com/amobe/gocov-merger@latest &&\
    go install github.com/nikolaydubina/go-cover-treemap@v1.4.2

WORKDIR /app