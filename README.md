# btool

The simplest C/C++ build tool.

![btool](btool.png)

## To try out `btool`...

```
$ docker run -v $PWD:/src -w /src -it ankeesler/btool
$ /btool -root BasicC -target main
...
$ ./BasicC/main
hey!
$ cat example/README.md # for more information about the examples
```

## To install the latest `btool` build...

```
$ ./script/install-btool.sh latest
```

## To install a local `btool` build...

```
$ ./script/install-btool.sh local
```

## To run the tests...

```
$ ./script/test.sh -u # unit tests
$ ./script/test.sh -i # integration tests
```

## Guiding Principles

- Ease of use
- Speed
- Extensibility
