#include "ui/ui.h"

#include "gtest/gtest.h"

TEST(UI, MakeNamePretty) {
  // Cache files
  EXPECT_EQ(
      "$CACHE/abc.../tuna.o",
      ::btool::ui::MakeNamePretty("/Users/marlin/.btool/abc123/some-lib/tuna.o",
                                  "/Users/marlin/.btool"));

  // Long files
  EXPECT_EQ("source/app.../tuna.o",
            ::btool::ui::MakeNamePretty("source/app/some/long/path/that/should/"
                                        "be/shortened/yes/some-lib/tuna.o",
                                        "/Users/marlin/.btool"));
}
