#include "util/sha256.h"

#include <ios>
#include <iostream>
#include <istream>
#include <string>
#include <vector>

#include "openssl/sha.h"

#include "err.h"
#include "log.h"
#include "util/util.h"

namespace btool::util {

static std::size_t Size(std::istream *is);

std::string SHA256(std::istream *is) {
  std::size_t data_size = Size(is);
  unsigned char *data = new unsigned char[data_size];
  is->read((char *)data, data_size);
  if (is->gcount() != data_size) {
    delete[] data;
    THROW_ERR("failed to read " + std::to_string(data_size) +
              " bytes from stream");
  }

  unsigned char md[SHA256_DIGEST_LENGTH];
  ::SHA256(data, data_size, md);

  delete[] data;

  return Hex(md, sizeof(md) / sizeof(md[0]));
}

static std::size_t Size(std::istream *is) {
  is->seekg(0, std::ios_base::end);
  std::size_t stream_size = is->tellg();
  is->seekg(0, std::ios_base::beg);
  return stream_size;
}

};  // namespace btool::util
