#include "util/util.h"

#include <chrono>
#include <thread>

#include "gtest/gtest.h"

TEST(Util, Hex) {
  unsigned char bytes[] = {
      0xAB,
      0x12,
      0xDE,
      0x34,
  };
  EXPECT_EQ("ab12de34",
            ::btool::util::Hex(bytes, sizeof(bytes) / sizeof(bytes[0])));
}

TEST(Util, CommaSeparatedNumber) {
  EXPECT_EQ("5", ::btool::util::CommaSeparatedNumber(5));
  EXPECT_EQ("500", ::btool::util::CommaSeparatedNumber(500));

  EXPECT_EQ("5,004", ::btool::util::CommaSeparatedNumber(5004));
  EXPECT_EQ("500,004", ::btool::util::CommaSeparatedNumber(500004));

  EXPECT_EQ("5,000,004", ::btool::util::CommaSeparatedNumber(5000004));
  EXPECT_EQ("50,000,000,000", ::btool::util::CommaSeparatedNumber(50000000000));
}

TEST(Util, Time) {
  std::chrono::milliseconds sleep_dur(500);
  auto dur = ::btool::util::Time(
      [&sleep_dur] { std::this_thread::sleep_for(sleep_dur); });
  EXPECT_TRUE(dur.count() >= sleep_dur.count());
}
