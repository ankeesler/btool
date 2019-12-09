FROM debian:10.2

RUN apt-get update \
  && apt-get install -y \
    curl \
    g++ \
    gcc \
    python \
    libfindbin-libs-perl \
    make \
    jq \
    -y

RUN  mkdir /tmp/workspace
COPY ./script/install-btool.sh /tmp/workspace
RUN /tmp/workspace/install-btool.sh latest
COPY ./source /tmp/workspace/source
RUN cd /tmp/workspace && /tmp/workspace/install-btool.sh local
RUN rm -rf /tmp/workspace

RUN mkdir /etc/btool
WORKDIR /etc/btool
COPY ./example /etc/btool
