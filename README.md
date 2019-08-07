# btool

Zero-configuration C/C++ build tool.

```
$ ./btool -root /tmp/BasicC -target main
$ /tmp/BasicC/main
hey! i am running!
```

### To build `btool`...

```
$ go build -o btool ./cmd/btool
$ ./btool -help
```

## To run the tests...

```
$ go test ./...
```
