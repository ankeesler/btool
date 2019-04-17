#include "log.h"

int main(int argc, char *argv[]) {
  Log log("main");
  log.Println("start");
  log.Println("end");
  return 0;
}
