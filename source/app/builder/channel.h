#ifndef BTOOL_APP_BUILDER_CHANNEL_H_
#define BTOOL_APP_BUILDER_CHANNEL_H_

#include <condition_variable>
#include <mutex>
#include <queue>

namespace btool::app::builder {

// A basic C++ implementation of the Golang "channel" concept.
template <typename T>
class Channel {
 public:
  Channel() : closed_(false) {}
  ~Channel() { Close(); }

  Channel<T>(const Channel<T> &) = delete;
  Channel<T> &operator=(const Channel<T> &) = delete;

  Channel<T>(Channel<T> &&) = delete;
  Channel<T> &operator=(Channel<T> &&) = delete;

  bool Tx(T t) {
    std::unique_lock<std::mutex> lock(mtx_);
    if (closed_) {
      return false;
    }

    q_.push(t);
    cv_.notify_one();
    return true;
  }

  bool Rx(T *t) {
    std::unique_lock<std::mutex> lock(mtx_);
    cv_.wait(lock, [this]() -> bool { return !q_.empty() || closed_; });
    if (closed_) {
      return false;
    }

    *t = q_.front();
    q_.pop();
    return true;
  }

  void Close() {
    std::unique_lock<std::mutex> lock(mtx_);
    closed_ = true;
    cv_.notify_all();
  }

  bool IsClosed() {
    std::unique_lock<std::mutex> lock(mtx_);
    return closed_;
  }

  std::size_t Size() {
    std::unique_lock<std::mutex> lock(mtx_);
    return q_.size();
  }

 private:
  bool closed_;
  std::queue<T> q_;
  std::condition_variable cv_;
  std::mutex mtx_;
};

};  // namespace btool::app::builder

#endif  // BTOOL_APP_BUILDER_CHANNEL_H_
