# btool

The simplest C/C++ build tool.

[![CircleCI](https://circleci.com/gh/ankeesler/btool/tree/master.svg?style=svg)](https://circleci.com/gh/ankeesler/btool/tree/master)


```
$ ./btool -root /tmp/BasicC -target main
$ /tmp/BasicC/main
hey! i am running!
```

![btool](btool.png)

### To try out `btool`...

```
$ docker run -it ankeesler/btool btool -root example/BasicCC -target main -run -loglevel error
hey!
```

#### OR

```
$ go run ./cmd/btool -root example/BasicCC -target main -loglevel error
hey!
```

## To run the tests...

```
$ go test ./...
```

## Guiding Principles

- Ease of use
- Speed
- Extensibility
