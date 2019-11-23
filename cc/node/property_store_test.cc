#include "node/property_store.h"

#include "gtest/gtest.h"

TEST(PropertyStore, Bool) {
  ::btool::node::PropertyStore ps;

  const bool *b;
  ps.Read("some-property", &b);
  EXPECT_EQ(nullptr, b);

  ps.Write("some-property", true);
  ps.Read("some-property", &b);
  EXPECT_EQ(true, *b);

  ps.Write("some-property", false);
  ps.Read("some-property", &b);
  EXPECT_EQ(false, *b);

  ps.Write("some-other-property", true);
  ps.Read("some-other-property", &b);
  EXPECT_EQ(true, *b);

  ::btool::node::PropertyStore copy = ps;
  copy.Read("some-property", &b);
  ASSERT_TRUE(b != nullptr);
  EXPECT_EQ(false, *b);
}

TEST(PropertyStore, Strings) {
  ::btool::node::PropertyStore ps;

  const std::vector<std::string> *s;
  ps.Read("some-property", &s);
  EXPECT_EQ(nullptr, s);

  ps.Append("some-property", "tuna");
  ps.Read("some-property", &s);
  EXPECT_EQ(1UL, s->size());
  EXPECT_EQ("tuna", s->at(0));

  ps.Append("some-property", "fish");
  ps.Append("some-property", "marlin");
  ps.Read("some-property", &s);
  EXPECT_EQ(3UL, s->size());
  EXPECT_EQ("tuna", s->at(0));
  EXPECT_EQ("fish", s->at(1));
  EXPECT_EQ("marlin", s->at(2));

  ::btool::node::PropertyStore copy = ps;
  copy.Read("some-property", &s);
  ASSERT_EQ(3UL, s->size());
  EXPECT_EQ("tuna", s->at(0));
  EXPECT_EQ("fish", s->at(1));
  EXPECT_EQ("marlin", s->at(2));
}
