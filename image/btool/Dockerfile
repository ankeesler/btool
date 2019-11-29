FROM debian:10.2

RUN apt-get update \
  && apt-get install -y \
    curl \
    g++ \
    gcc \
    python \
    -y

RUN curl \
  -L https://github.com/ankeesler/btool/releases/download/0.0/btool-0.0-linux-x86_64.gz \
  -o - \
  | gunzip \
  > /usr/local/bin/btool \
  && chmod +x /usr/local/bin/btool

RUN mkdir /etc/btool
WORKDIR /etc/btool
COPY ./example /etc/btool
