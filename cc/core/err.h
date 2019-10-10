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
// Err is movable and copyable.
//
// Declaring an Err with the default constructor means Success.
template <typename T>
class Err {
 public:
  Err<T>(T ret) : ret_(ret), msg_(nullptr) {}

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

  Err<T>(Err<T> &err) = default;
  Err<T> &operator=(Err<T> &err) = default;

  Err<T>(Err<T> &&err) = default;
  Err<T> &operator=(Err<T> &&err) = default;

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
};

template <typename T>
std::ostream &operator<<(std::ostream &os, const Err<T> &err) {
  if (err) {
    os << "err: failure: " << err.Msg();
  } else {
    os << "err: success: " << err.Ret();
  }
  return os;
}

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

  VoidErr() : msg_(nullptr) {}

  VoidErr(VoidErr &err) = default;
  VoidErr &operator=(VoidErr &err) = default;

  VoidErr(VoidErr &&err) = default;
  VoidErr &operator=(VoidErr &&err) = default;

  bool operator==(const VoidErr &err) const {
    if (*this) {
      return err && (::strcmp(msg_, err.msg_) == 0);
    } else {
      return !err;
    }
  }

  bool operator!=(const VoidErr &err) const { return !(*this == err); }

  operator bool() const { return msg_ != nullptr; }

  const char *Msg() const { return msg_; }

 private:
  const char *msg_;
};

std::ostream &operator<<(std::ostream &os, const VoidErr &err);

};  // namespace btool::core

#endif  // BTOOL_CORE_ERR_H_
