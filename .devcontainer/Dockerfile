FROM golang:latest

RUN apt update -y && apt upgrade -y
RUN apt install -y git wget curl unzip

RUN wget dl.elv.sh/linux-arm64/elvish-HEAD.tar.gz
RUN tar -C /usr/local/bin -xvf elvish-HEAD.tar.gz

WORKDIR /workspace
COPY .. .
