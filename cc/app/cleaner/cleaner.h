#ifndef BTOOL_APP_CLEANER_CLEANER_H_
#define BTOOL_APP_CLEANER_CLEANER_H_

#include "node/node.h"

namespace btool::app::cleaner {

class Cleaner {
 public:
  class RemoveAller {
   public:
    virtual ~RemoveAller() {}
    virtual bool RemoveAll(const std::string &path, std::string *err) = 0;
  };

  Cleaner(RemoveAller *ra) : ra_(ra){};

  bool Clean(const ::btool::node::Node &node, std::string *err);

 private:
  RemoveAller *ra_;
};

};  // namespace btool::app::cleaner

#endif  // BTOOL_APP_CLEANER_CLEANER_H_
