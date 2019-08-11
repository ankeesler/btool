# btool

Zero-configuration C/C++ build tool.

```
$ ./btool -root /tmp/BasicC -target main
$ /tmp/BasicC/main
hey! i am running!
```

### To try out `btool`...

```
$ docker run -it ankeesler/btool
$ btool -help
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
