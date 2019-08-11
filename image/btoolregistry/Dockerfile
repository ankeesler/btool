FROM golang:buster
WORKDIR /app
COPY . /app
RUN go env && CGO_ENABLED=0 go build -o /tmp/registry ./cmd/registry

FROM alpine
WORKDIR /app
COPY --from=0 /tmp/registry .
COPY --from=0 /app/data data
EXPOSE 8080
CMD [ \
  "./registry", \
  "-loglevel", \
  "trace", \
  "-dir", \
  "data", \
  "-address", \
  ":8080" \
]
