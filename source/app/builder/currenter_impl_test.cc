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

  EXPECT_FALSE(ci.Current(tmpd_));

  NodeFile nf(tmpd_);
  nf.Modify();

  EXPECT_TRUE(ci.Current(tmpd_));
}

TEST_F(CurrenterTest, AdjacentDep) {
  ::btool::app::builder::CurrenterImpl ci;

  EXPECT_FALSE(ci.Current(tmpc_));

  NodeFile nfd(tmpd_);
  EXPECT_FALSE(ci.Current(tmpc_));

  NodeFile nfc(tmpc_);
  EXPECT_TRUE(ci.Current(tmpc_));

  nfd.Modify();
  EXPECT_FALSE(ci.Current(tmpc_));

  nfc.Modify();
  EXPECT_TRUE(ci.Current(tmpc_));
}

TEST_F(CurrenterTest, AncestorDep) {
  ::btool::app::builder::CurrenterImpl ci;

  EXPECT_FALSE(ci.Current(tmpb_));

  NodeFile nfd(tmpd_);
  NodeFile nfc(tmpc_);
  EXPECT_FALSE(ci.Current(tmpb_));

  NodeFile nfb(tmpb_);
  EXPECT_TRUE(ci.Current(tmpb_));

  nfd.Modify();
  EXPECT_FALSE(ci.Current(tmpb_));

  nfc.Modify();
  EXPECT_FALSE(ci.Current(tmpb_));

  nfb.Modify();
  EXPECT_TRUE(ci.Current(tmpb_));
}
