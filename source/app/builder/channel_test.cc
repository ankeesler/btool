#include "app/builder/channel.h"

#include <chrono>
#include <future>
#include <thread>

#include "gtest/gtest.h"

TEST(Channel, OneAfterAnother) {
  ::btool::app::builder::Channel<int> c;
  for (std::size_t i = 0; i < 10; ++i) {
    EXPECT_TRUE(c.Tx(i));

    EXPECT_EQ(1UL, c.Size());

    int n;
    EXPECT_TRUE(c.Rx(&n));
    EXPECT_EQ(i, static_cast<std::size_t>(n));
  }

  EXPECT_EQ(0UL, c.Size());
}

TEST(Channel, AllInAllOut) {
  ::btool::app::builder::Channel<int> c;

  for (std::size_t i = 0; i < 10; ++i) {
    EXPECT_TRUE(c.Tx(i));
  }

  EXPECT_EQ(10UL, c.Size());

  for (std::size_t i = 0; i < 10; ++i) {
    int n;
    EXPECT_TRUE(c.Rx(&n));
    EXPECT_EQ(i, static_cast<std::size_t>(n));
  }

  EXPECT_EQ(0UL, c.Size());
}

TEST(Channel, MultiProducer) {
  ::btool::app::builder::Channel<int> c;

  std::future<int> frx = std::async([&c]() -> int {
    int all = 0;
    int count = 0;
    int n;
    while (c.Rx(&n)) {
      count++;
      all |= 1 << n;

      if (count == 3) {
        break;
      }
    }
    return all;
  });
  std::future<bool> f0 = std::async([&c]() -> bool { return c.Tx(0); });
  std::future<bool> f1 = std::async([&c]() -> bool { return c.Tx(1); });
  std::future<bool> f2 = std::async([&c]() -> bool { return c.Tx(2); });

  EXPECT_EQ(std::future_status::ready,
            f0.wait_for(std::chrono::milliseconds(500)));
  EXPECT_TRUE(f0.get());

  EXPECT_EQ(std::future_status::ready,
            f1.wait_for(std::chrono::milliseconds(500)));
  EXPECT_TRUE(f1.get());

  EXPECT_EQ(std::future_status::ready,
            f2.wait_for(std::chrono::milliseconds(500)));
  EXPECT_TRUE(f2.get());

  ASSERT_EQ(std::future_status::ready,
            frx.wait_for(std::chrono::milliseconds(500)));
  EXPECT_EQ(0x07, frx.get());
}

TEST(Channel, MultiConsumer) {
  ::btool::app::builder::Channel<int> c_in;
  ::btool::app::builder::Channel<int> c_out;

  auto rx = [&c_in, &c_out](int i) {
    int n;
    while (c_in.Rx(&n)) {
      c_out.Tx((1 << n) | (1 << (8 + i)));
    }
  };

  std::thread t0(rx, 0);
  std::thread t1(rx, 1);
  std::thread t2(rx, 2);

  for (std::size_t i = 0; i < 8; ++i) {
    EXPECT_TRUE(c_in.Tx(i));
  }

  int n, all = 0;
  for (std::size_t i = 0; i < 8; ++i) {
    EXPECT_TRUE(c_out.Rx(&n));
    EXPECT_EQ(0, ((n & 0x00FF) & all))
        << "consumer " << (n >> 8) << " received duplicate number "
        << (n & 0x00FF) << std::endl;
    ;
    all |= (n & 0x00FF);
  }
  EXPECT_EQ(0x00FF, all);

  c_in.Close();
  EXPECT_TRUE(c_in.IsClosed());

  t0.join();
  t1.join();
  t2.join();
}
