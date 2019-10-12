#include "fs_collectini.h"

#include "gtest/gtest.h"

#include "node/store.h"

// TODO: these director paths should be more reliable. Test utility?

TEST(FSCollectini, Yeah) {
  ::btool::app::collector::fs::FSCollectini fsc("cc/testdata/app/collector/fs");
  ::btool::node::Store s;
  EXPECT_EQ(::btool::core::VoidErr::Success(), fsc.Collect(&s));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/a.cc"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/b.h"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/c.c"));
  EXPECT_FALSE(s.Get("cc/testdata/app/collector/fs/d.go"));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/a.cc"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/b.h"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/c.c"));
  EXPECT_FALSE(s.Get("cc/testdata/app/collector/fs/dir0/d.go"));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/dir0/a.cc"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/dir0/b.h"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/dir0/c.c"));
  EXPECT_FALSE(s.Get("cc/testdata/app/collector/fs/dir0/dir0/d.go"));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/dir1/a.cc"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/dir1/b.h"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir0/dir1/c.c"));
  EXPECT_FALSE(s.Get("cc/testdata/app/collector/fs/dir0/dir1/d.go"));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/a.cc"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/b.h"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/c.c"));
  EXPECT_FALSE(s.Get("cc/testdata/app/collector/fs/dir1/d.go"));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/dir0/a.cc"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/dir0/b.h"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/dir0/c.c"));
  EXPECT_FALSE(s.Get("cc/testdata/app/collector/fs/dir1/dir0/d.go"));

  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/dir1/a.cc"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/dir1/b.h"));
  EXPECT_TRUE(s.Get("cc/testdata/app/collector/fs/dir1/dir1/c.c"));
  EXPECT_FALSE(s.Get("cc/testdata/app/collector/fs/dir1/dir1/d.go"));
}
