#include "currenter_impl.h"

#include <errno.h>
#include <cassert>
#include <cstdio>
#include <cstring>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "node/testing/node.h"

class CurrenterTest : public ::btool::node::testing::NodeTest {};

class NodeFile {
 public:
  NodeFile(const ::btool::node::Node &node) : node_(node) {
    ::usleep(100000);

    ::FILE *f = ::fopen(node.name().c_str(), "w");
    EXPECT_TRUE(f != nullptr) << "cannot open " << node.name().c_str();
    fprintf(f, "open\n");
    fflush(f);
    EXPECT_EQ(0, fclose(f));
  }

  ~NodeFile() {
    int err = ::remove(node_.name().c_str());
    EXPECT_EQ(0, err) << "remove: " << strerror(errno);
  }

  void Modify() {
    ::usleep(100000);

    ::FILE *f = ::fopen(node_.name().c_str(), "a");
    fprintf(f, "modify\n");
    fflush(f);
    EXPECT_EQ(0, fclose(f));
  }

 private:
  ::btool::node::Node node_;
};

TEST_F(CurrenterTest, NoDeps) {
  ::btool::app::builder::CurrenterImpl ci;

  bool current;
  std::string err;
  EXPECT_TRUE(ci.Current(tmpd_, &current, &err)) << "error: " << err;
  EXPECT_FALSE(current);

  NodeFile nf(tmpd_);
  nf.Modify();

  EXPECT_TRUE(ci.Current(tmpd_, &current, &err)) << "error: " << err;
  EXPECT_TRUE(current);
}

TEST_F(CurrenterTest, AdjacentDep) {
  ::btool::app::builder::CurrenterImpl ci;

  bool current;
  std::string err;
  EXPECT_TRUE(ci.Current(tmpc_, &current, &err)) << "error: " << err;
  EXPECT_FALSE(current);

  NodeFile nfd(tmpd_);
  EXPECT_TRUE(ci.Current(tmpc_, &current, &err)) << "error: " << err;
  EXPECT_FALSE(current);

  NodeFile nfc(tmpc_);
  EXPECT_TRUE(ci.Current(tmpc_, &current, &err)) << "error: " << err;
  EXPECT_TRUE(current);

  nfd.Modify();
  EXPECT_TRUE(ci.Current(tmpc_, &current, &err)) << "error: " << err;
  EXPECT_FALSE(current);

  nfc.Modify();
  EXPECT_TRUE(ci.Current(tmpc_, &current, &err)) << "error: " << err;
  EXPECT_TRUE(current);
}

TEST_F(CurrenterTest, AncestorDep) {
  ::btool::app::builder::CurrenterImpl ci;

  bool current;
  std::string err;
  EXPECT_TRUE(ci.Current(tmpb_, &current, &err)) << "error: " << err;
  EXPECT_FALSE(current);

  NodeFile nfd(tmpd_);
  NodeFile nfc(tmpc_);
  EXPECT_TRUE(ci.Current(tmpb_, &current, &err)) << "error: " << err;
  EXPECT_FALSE(current);

  NodeFile nfb(tmpb_);
  EXPECT_TRUE(ci.Current(tmpb_, &current, &err)) << "error: " << err;
  EXPECT_TRUE(current);

  nfd.Modify();
  EXPECT_TRUE(ci.Current(tmpb_, &current, &err)) << "error: " << err;
  EXPECT_FALSE(current);

  nfc.Modify();
  EXPECT_TRUE(ci.Current(tmpb_, &current, &err)) << "error: " << err;
  EXPECT_FALSE(current);

  nfb.Modify();
  EXPECT_TRUE(ci.Current(tmpb_, &current, &err)) << "error: " << err;
  EXPECT_TRUE(current);
}
