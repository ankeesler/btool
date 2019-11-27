# btool

The simplest C/C++ build tool.

[![CircleCI](https://circleci.com/gh/ankeesler/btool/tree/master.svg?style=svg)](https://circleci.com/gh/ankeesler/btool/tree/master)

```
$ ./btool -root ./BasicC -target main
$ ./BasicC/main
hey! i am running!
```

![btool](btool.png)

### To try out `btool`...

```
$ docker run -it ankeesler/btool
$ btool -root BasicC -target btool
$ ./cc/btool
hey! i am running!
```

## To build `btool`...

```
$ btool -root cc -target btool
```

## To run the tests...

```
$ ./script/test.sh
```

## Guiding Principles

- Ease of use
- Speed
- Extensibility
