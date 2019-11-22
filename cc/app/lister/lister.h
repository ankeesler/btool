#ifndef BTOOL_APP_LISTER_LISTER_H_
#define BTOOL_APP_LISTER_LISTER_H_

#include <ostream>

#include "app/app.h"
#include "node/node.h"

namespace btool::app::lister {

class Lister : public ::btool::app::App::Lister {
 public:
  Lister(std::ostream *os) : os_(os) {}

  bool List(const ::btool::node::Node &node, std::string *ret_err) override;

 private:
  std::ostream *os_;
};

};  // namespace btool::app::lister

#endif  // BTOOL_APP_LISTER_LISTER_H_
