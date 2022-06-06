FROM ubuntu:latest

WORKDIR /work

RUN apt-get update && \
  apt-get install -y yasm nasm gcc && \
  apt-get install -y golang-go && \
  apt-get install -y make && \
  apt-get install -y lldb