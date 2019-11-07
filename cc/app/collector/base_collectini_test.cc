#include "app/collector/base_collectini.h"

#include <vector>
#include <string>

#include "gtest/gtest.h"

#include "app/collector/store.h"

class Collectini : public ::btool::app::collector::BaseCollectini {
  public:
    std::vector<::btool::app::collector::Store *> collect_calls_;
    std::vector<std::pair<::btool::app::collector::Store *, std::string>> on_notify_calls_;

    void Collect(::btool::app::collector::Store *s) override {
      collect_calls_.push_back(s);
      Notify(s, "some-other-name");
    }

  protected:
    void OnNotify(::btool::app::collector::Store *s, const std::string &name) override {
      on_notify_calls_.push_back({s, name});
    }
};

TEST(BaseCollectini, A) {
  Collectini a;
  Collectini b;
  Collectini c;

  ::btool::app::collector::Store s;
  a.Collect(&s);
  EXPECT_EQ(&s, a.collect_calls_[0]);

  EXPECT_EQ(0UL, a.on_notify_calls_.size());
  EXPECT_EQ(&s, b.on_notify_calls_[0].first);
  EXPECT_EQ("some-other-name", b.on_notify_calls_[0].second);
  EXPECT_EQ(&s, c.on_notify_calls_[0].first);
  EXPECT_EQ("some-other-name", c.on_notify_calls_[0].second);
}
