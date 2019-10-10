#include "fs_collectini.h"

#include "gtest/gtest.h"

#include "node/store.h"

// TODO: these director paths should be more reliable. Test utility?

TEST(FSCollectini, Yeah) {
  ::btool::app::collector::fs::FSCollectini fsc("cc/testdata/app/collector/fs");
  ::btool::node::Store s;
  EXPECT_EQ(::btool::core::VoidErr::Success(), fsc.Collect(&s));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/a")) << s;
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/b"));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/a"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/b"));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/dir0/a"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/dir0/b"));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/dir1/a"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/dir1/b"));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/a"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/b"));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/dir0/a"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/dir0/b"));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/dir1/a"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/dir1/b"));
}
