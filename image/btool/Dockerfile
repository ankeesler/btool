FROM debian:10.2

RUN apt-get update \
  && apt-get install -y \
    curl \
    g++ \
    gcc \
    python \
    libfindbin-libs-perl \
    make \
    -y

COPY ./script/install-btool.sh /tmp/install-btool.sh
RUN /tmp/install-btool.sh

RUN mkdir /etc/btool
WORKDIR /etc/btool
COPY ./example /etc/btool
