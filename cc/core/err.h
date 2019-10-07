#ifndef BTOOL_CORE_ERR_H_
#define BTOOL_CORE_ERR_H_

#include <cstring>

#include <iostream>
#include <string>

namespace btool::core {

// Err
//
// Err is an error message + return value pair. Only one of these is active at
// once.
//
// Err is movable, but not copyable.
template <typename T>
class Err {
 public:
  static Err<T> Success(T ret) {
    Err<T> err;
    err.ret_ = ret;
    err.msg_ = nullptr;
    return err;
  }

  static Err<T> Failure(const char *msg) {
    Err<T> err;
    err.msg_ = msg;
    return err;
  }

  Err<T>(Err<T> &&err) : ret_(std::move(err.ret_)), msg_(err.msg_) {}

  Err<T> &operator=(Err<T> &&err) {
    ret_ = std::move(err.ret_);
    msg_ = err.msg_;
    return *this;
  }

  bool operator==(const Err<T> &err) const {
    if (*this) {
      return err && (::strcmp(msg_, err.msg_) == 0);
    } else {
      return !err && ret_ == err.ret_;
    }
  }

  bool operator!=(const Err<T> &err) const { return !(*this == err); }

  operator bool() const { return msg_ != nullptr; }

  T Ret() const { return ret_; }
  const char *Msg() const { return msg_; }

 private:
  Err<T>() {}

  T ret_;
  const char *msg_;
};  // namespace btool::core

// VoidErr
//
// This type behaves exactly as an Err<void> would from above.
class VoidErr {
 public:
  static VoidErr Success() {
    VoidErr err;
    err.msg_ = nullptr;
    return err;
  }

  static VoidErr Failure(const char *msg) {
    VoidErr err;
    err.msg_ = msg;
    return err;
  }

  VoidErr(VoidErr &&err) : msg_(err.msg_) {}

  VoidErr &operator=(VoidErr &&err) {
    msg_ = err.msg_;
    return *this;
  }

  bool operator==(const VoidErr &err) const {
    if (*this) {
      return err && (::strcmp(msg_, err.msg_) == 0);
    } else {
      return true;
    }
  }

  operator bool() const { return msg_ != nullptr; }

  const char *Msg() const { return msg_; }

 private:
  VoidErr() {}

  const char *msg_;
};
};  // namespace btool::core

#endif  // BTOOL_CORE_ERR_H_
