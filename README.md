# btool

The simplest C/C++ build tool.

![btool](btool.png)

```
$ ./btool -root /tmp/BasicC -target main
$ /tmp/BasicC/main
hey! i am running!
```

### To try out `btool`...

```
$ docker run -it ankeesler/btool
$ btool -root example/BasicCC -target main
```

#### OR

```
$ go build -o btool ./cmd/btool
$ ./btool -root example/BasicCC -target main
```

## To run the tests...

```
$ go test ./...
```

## Guiding Principles

- Ease of use
- Speed
- Extensibility