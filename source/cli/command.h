#include <string>

namespace btool::cli {

class Command {
public:
  virtual const std::string& Name() const = 0;
};

};
