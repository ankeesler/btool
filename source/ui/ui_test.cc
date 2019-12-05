#include "ui/ui.h"

#include "gtest/gtest.h"

TEST(UI, MakeNamePretty) {
  EXPECT_EQ(
      "$CACHE/abc.../tuna.o",
      ::btool::ui::MakeNamePretty("/Users/marlin/.btool/abc123/some-lib/tuna.o",
                                  "/Users/marlin/.btool"));
}
