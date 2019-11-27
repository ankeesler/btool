FROM debian:10.2

RUN apt-get update \
  && apt-get install curl emacs g++ gcc gdb git man manpages manpages-dev python valgrind -y

RUN curl \
  -L https://github.com/ankeesler/btool/releases/download/0.0/btool-0.0-linux-x86_64.gz \
  -o - \
  | gunzip \
  > /usr/local/bin/btool
RUN chmod +x /usr/local/bin/btool

RUN mkdir /etc/btool
WORKDIR /etc/btool
COPY ./example /etc/btool

RUN mkdir /etc/ci