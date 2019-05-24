# btool

Push to start C/C++ building!

```bash
$ ./btool -root fixture/ComplexCC/ -target fixture/ComplexCC/main.cc
btool |  info | scanning from file fixture/ComplexCC/main.cc
btool |  info | building graph
btool |  info | compiling node fixture/ComplexCC/dep-0/dep-0a.cc
btool |  info | compiling node fixture/ComplexCC/dep-1/dep-1a.cc
btool |  info | compiling node fixture/ComplexCC/dep-2/dep-2a.cc
btool |  info | compiling node fixture/ComplexCC/main.cc
btool |  info | linking
$ ./.btool/binaries/out
hey! i am running!
```