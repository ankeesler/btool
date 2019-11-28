#ifndef BTOOL_APP_APP_H_
#define BTOOL_APP_APP_H_

#include <string>

#include "err.h"
#include "node/node.h"

namespace btool::app {

class App {
 public:
  class Collector {
   public:
    virtual ~Collector() {}
    virtual ::btool::node::Node *Collect(const std::string &name) = 0;
  };

  class Cleaner {
   public:
    virtual ~Cleaner() {}
    virtual void Clean(const ::btool::node::Node &) = 0;
  };

  class Lister {
   public:
    virtual ~Lister() {}
    virtual void List(const ::btool::node::Node &) = 0;
  };

  class Builder {
   public:
    virtual ~Builder() {}
    virtual void Build(const ::btool::node::Node &) = 0;
  };

  class Runner {
   public:
    virtual ~Runner() {}
    virtual void Run(const ::btool::node::Node &) = 0;
  };

  App(Collector *collector, Cleaner *cleaner, Lister *lister, Builder *builder,
      Runner *runner)
      : collector_(collector),
        cleaner_(cleaner),
        lister_(lister),
        builder_(builder),
        runner_(runner) {}

  void Run(const std::string &target, bool clean, bool list, bool run);

 private:
  Collector *collector_;
  Cleaner *cleaner_;
  Lister *lister_;
  Builder *builder_;
  Runner *runner_;
};

};  // namespace btool::app

#endif  // BTOOL_APP_APP_H_
