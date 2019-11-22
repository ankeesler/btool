#ifndef BTOOL_ERR_H_
#define BTOOL_ERR_H_

#include <iostream>
#include <optional>
#include <string>

namespace btool {

inline std::string WrapErr(std::string wrapped, std::string wrapper) {
  return wrapper + std::string(": ") + wrapped;
}

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
  Err<T>(T ret) : ret_(ret), msg_(std::nullopt) {}

  static Err<T> Success(T ret) {
    Err<T> err;
    err.ret_ = ret;
    err.msg_ = std::nullopt;
    return err;
  }

  static Err<T> Failure(std::string msg) {
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
      return err && msg_ == err.msg_;
    } else {
      return !err && ret_ == err.ret_;
    }
  }

  bool operator!=(const Err<T> &err) const { return !(*this == err); }

  operator bool() const { return msg_.has_value(); }

  T Ret() const { return ret_; }
  std::string Msg() const { return msg_.value(); }

 private:
  Err<T>() {}

  T ret_;
  std::optional<std::string> msg_;
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
    err.msg_ = std::nullopt;
    return err;
  }

  static VoidErr Failure(std::string msg) {
    VoidErr err;
    err.msg_ = msg;
    return err;
  }

  VoidErr() : msg_(std::nullopt) {}

  VoidErr(VoidErr &err) = default;
  VoidErr &operator=(VoidErr &err) = default;

  VoidErr(VoidErr &&err) = default;
  VoidErr &operator=(VoidErr &&err) = default;

  bool operator==(const VoidErr &err) const {
    if (*this) {
      return err && msg_ == err.msg_;
    } else {
      return !err;
    }
  }

  bool operator!=(const VoidErr &err) const { return !(*this == err); }

  operator bool() const { return msg_.has_value(); }

  std::string Msg() const { return msg_.value(); }

 private:
  std::optional<std::string> msg_;
};

std::ostream &operator<<(std::ostream &os, const VoidErr &err);

};  // namespace btool

#endif  // BTOOL_ERR_H_
