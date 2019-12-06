#include "node/property_store.h"

#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

using ::testing::Contains;

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

  std::vector<std::pair<std::string, ::btool::node::PropertyStore::Type>> stuff;
  ps.ForEach([&stuff](const std::string &name,
                      ::btool::node::PropertyStore::Type type) {
    stuff.push_back(
        std::pair<std::string, ::btool::node::PropertyStore::Type>(name, type));
  });
  EXPECT_THAT(
      stuff,
      Contains(std::pair<std::string, ::btool::node::PropertyStore::Type>(
          "some-property", ::btool::node::PropertyStore::kBool)));
  EXPECT_THAT(
      stuff,
      Contains(std::pair<std::string, ::btool::node::PropertyStore::Type>(
          "some-other-property", ::btool::node::PropertyStore::kBool)));
}

TEST(PropertyStore, String) {
  ::btool::node::PropertyStore ps;

  const std::string *s;
  ps.Read("some-property", &s);
  EXPECT_EQ(nullptr, s);

  ps.Write("some-property", "tuna");
  ps.Read("some-property", &s);
  EXPECT_EQ("tuna", *s);

  ps.Write("some-other-property", "fish");
  ps.Read("some-property", &s);
  EXPECT_EQ("tuna", *s);
  ps.Read("some-other-property", &s);
  EXPECT_EQ("fish", *s);

  ::btool::node::PropertyStore copy = ps;
  copy.Read("some-property", &s);
  EXPECT_EQ("tuna", *s);
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

  ps.ForEach("some-property", [](std::string *v) { v->insert(0, "a-"); });
  EXPECT_EQ(3UL, s->size());
  EXPECT_EQ("a-tuna", s->at(0));
  EXPECT_EQ("a-fish", s->at(1));
  EXPECT_EQ("a-marlin", s->at(2));

  ::btool::node::PropertyStore copy = ps;
  copy.Read("some-property", &s);
  ASSERT_EQ(3UL, s->size());
  EXPECT_EQ("a-tuna", s->at(0));
  EXPECT_EQ("a-fish", s->at(1));
  EXPECT_EQ("a-marlin", s->at(2));
}
