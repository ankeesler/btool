FROM debian:10.2

RUN apt-get update \
  && apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    emacs \
    g++ \
    gcc \
    gdb \
    git \
    gnupg2 \
    man \
    manpages \
    manpages-dev \
    python \
    software-properties-common \
    valgrind \
    wget \
    -y
RUN curl -fsSL https://download.docker.com/linux/debian/gpg \
  | apt-key add -
RUN add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/debian \
   $(lsb_release -cs) \
   stable"
RUN apt-get update && apt-get install docker-ce -y

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