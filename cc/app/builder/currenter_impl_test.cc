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

    ::FILE *f = ::fopen(node.Name().c_str(), "w");
    EXPECT_TRUE(f != nullptr) << "cannot open " << node.Name().c_str();
    fprintf(f, "open\n");
    fflush(f);
    EXPECT_EQ(0, fclose(f));
  }

  ~NodeFile() {
    int err = ::remove(node_.Name().c_str());
    EXPECT_EQ(0, err) << "remove: " << strerror(errno);
  }

  void Modify() {
    ::usleep(100000);

    ::FILE *f = ::fopen(node_.Name().c_str(), "a");
    fprintf(f, "modify\n");
    fflush(f);
    EXPECT_EQ(0, fclose(f));
  }

 private:
  ::btool::node::Node node_;
};

TEST_F(CurrenterTest, NoDeps) {
  ::btool::app::builder::CurrenterImpl ci;

  auto err = ci.Current(tmpd_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(false), err);

  NodeFile nf(tmpd_);
  nf.Modify();

  err = ci.Current(tmpd_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(true), err);
}

TEST_F(CurrenterTest, AdjacentDep) {
  ::btool::app::builder::CurrenterImpl ci;

  auto err = ci.Current(tmpc_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(false), err);

  NodeFile nfd(tmpd_);
  err = ci.Current(tmpc_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(false), err);

  NodeFile nfc(tmpc_);
  err = ci.Current(tmpc_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(true), err);

  nfd.Modify();
  err = ci.Current(tmpc_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(false), err);

  nfc.Modify();
  err = ci.Current(tmpc_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(true), err);
}

TEST_F(CurrenterTest, AncestorDep) {
  ::btool::app::builder::CurrenterImpl ci;

  auto err = ci.Current(tmpb_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(false), err);

  NodeFile nfd(tmpd_);
  NodeFile nfc(tmpc_);
  err = ci.Current(tmpb_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(false), err);

  NodeFile nfb(tmpb_);
  err = ci.Current(tmpb_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(true), err);

  nfd.Modify();
  err = ci.Current(tmpb_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(false), err);

  nfc.Modify();
  err = ci.Current(tmpb_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(false), err);

  nfb.Modify();
  err = ci.Current(tmpb_);
  EXPECT_EQ(::btool::core::Err<bool>::Success(true), err);
}
