FROM debian:10.2 as builder

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

WORKDIR /workspace
COPY script script
COPY source source
RUN script/install-btool.sh latest
RUN script/install-btool.sh local

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/source/btool .
USER nonroot:nonroot
ENTRYPOINT ["/btool"]
