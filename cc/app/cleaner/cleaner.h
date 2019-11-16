#ifndef BTOOL_APP_CLEANER_CLEANER_H_
#define BTOOL_APP_CLEANER_CLEANER_H_

#include <string>

#include "app/app.h"
#include "err.h"
#include "node/node.h"

namespace btool::app::cleaner {

class Cleaner : public ::btool::app::App::Cleaner {
 public:
  class RemoveAller {
   public:
    virtual ~RemoveAller() {}
    virtual ::btool::VoidErr RemoveAll(const std::string &path) = 0;
  };

  Cleaner(RemoveAller *ra) : ra_(ra){};

  ::btool::VoidErr Clean(const ::btool::node::Node &node) override;

 private:
  RemoveAller *ra_;
};

};  // namespace btool::app::cleaner

#endif  // BTOOL_APP_CLEANER_CLEANER_H_
