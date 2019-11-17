#include "err.h"

#include <iostream>
#include <ostream>

#include "gtest/gtest.h"

static ::btool::Err<int> IntSuccess() { return ::btool::Err<int>::Success(5); }

static ::btool::Err<int> IntFailure() {
  return ::btool::Err<int>::Failure("oh no");
}

static ::btool::VoidErr VoidSuccess() { return ::btool::VoidErr::Success(); }

static ::btool::VoidErr VoidFailure() {
  return ::btool::VoidErr::Failure("bummer");
}

class Tuna {
 public:
  int Fish() const { return 5; }

  bool operator==(const Tuna &t) const { return true; }
};

std::ostream &operator<<(std::ostream &os, const Tuna &tuna) {
  os << "tuna: " << tuna.Fish();
  return os;
}

static ::btool::Err<Tuna> TunaSuccess() {
  Tuna t;
  return ::btool::Err<Tuna>::Success(t);
}

static ::btool::Err<Tuna> TunaFailure() {
  return ::btool::Err<Tuna>::Failure("no bueno");
}

TEST(Err, Int) {
  auto s = IntSuccess();
  EXPECT_FALSE(s);
  EXPECT_EQ(5, s.Ret());
  EXPECT_EQ(s, IntSuccess());

  auto f = IntFailure();
  EXPECT_TRUE(f);
  EXPECT_EQ("oh no", f.Msg());
  EXPECT_EQ(f, IntFailure());

  EXPECT_NE(s, f);

  auto fa0 = ::btool::Err<int>::Failure("nope-0");
  auto fb0 = ::btool::Err<int>::Failure("nope-0");
  auto fa1 = ::btool::Err<int>::Failure("nope-1");
  EXPECT_EQ(fa0, fb0);
  EXPECT_NE(fa0, fa1);

  auto sa0 = ::btool::Err<int>::Success(0);
  auto sb0 = ::btool::Err<int>::Success(0);
  auto sa1 = ::btool::Err<int>::Success(1);
  EXPECT_EQ(sa0, sb0);
  EXPECT_NE(sa0, sa1);
}

TEST(Err, Void) {
  auto s = VoidSuccess();
  EXPECT_FALSE(s);
  EXPECT_EQ(s, VoidSuccess());

  auto f = VoidFailure();
  EXPECT_TRUE(f);
  EXPECT_EQ("bummer", f.Msg());
  EXPECT_EQ(f, VoidFailure());

  EXPECT_NE(s, f);

  ::btool::VoidErr err;
  EXPECT_EQ(err, ::btool::VoidErr::Success());
  EXPECT_EQ(err, s);
}

TEST(Err, Tuna) {
  auto s = TunaSuccess();
  Tuna t;
  EXPECT_FALSE(s);
  EXPECT_EQ(t, s.Ret());
  EXPECT_EQ(s, TunaSuccess());

  auto f = TunaFailure();
  EXPECT_TRUE(f);
  EXPECT_EQ("no bueno", f.Msg());
  EXPECT_EQ(f, TunaFailure());

  EXPECT_NE(s, f);
}
