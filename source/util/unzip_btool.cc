#include "util/unzip_btool.h"

#include <errno.h>

#include <cstdlib>
#include <cstring>
#include <fstream>
#include <iostream>

#include "err.h"
#include "log.h"
#include "util/fs/fs.h"
#include "zip.h"

namespace btool::util {

static void Handle(zip_file_t *f, const char *name, zip_uint64_t size,
                   const std::string &dest_dir);
static bool ErrorIsInvalidArgument(zip_t *z);

void Unzip(const std::string &zipfile, const std::string &dest_dir) {
  int e = 0;
  zip_t *z = zip_open(zipfile.c_str(), ZIP_CHECKCONS | ZIP_RDONLY, &e);
  if (z == nullptr) {
    zip_error_t error;
    zip_error_init_with_code(&error, e);
    std::string message = "could not open zipfile " + zipfile + ": " +
                          std::string(zip_error_strerror(&error));
    zip_error_fini(&error);

    THROW_ERR(message);
  }
  DEBUGS() << "opened zipfile " << zipfile << std::endl;

  const char *name = nullptr;
  zip_uint64_t index = 0;
  while ((name = zip_get_name(z, index, 0)) != nullptr) {
    DEBUGS() << "found file named " << name << std::endl;

    zip_stat_t s;
    if (zip_stat(z, name, 0, &s) != 0) {
      std::string message = "failed to stat file " + std::string(name) +
                            " in zipfile " + zipfile + ": " +
                            std::string(zip_strerror(z));
      zip_close(z);

      THROW_ERR(message);
    }

    if ((s.valid & ZIP_STAT_SIZE) == 0) {
      std::string message = "cannot get size of file " + std::string(name) +
                            " in zipfile " + zipfile + ": " +
                            std::string(zip_strerror(z));
      zip_close(z);

      THROW_ERR(message);
    }
    zip_uint64_t size = s.size;

    zip_file_t *f = zip_fopen(z, name, 0);
    if (f == nullptr) {
      std::string message = "cannot open file " + std::string(name) +
                            " in zipfile " + zipfile + ": " +
                            std::string(zip_strerror(z));
      zip_close(z);

      THROW_ERR(message);
    }

    try {
      Handle(f, name, size, dest_dir);
    } catch (const ::btool::Err &err) {
      zip_close(z);
      zip_fclose(f);
      THROW_ERR("failed to handle file " + std::string(name) + " in zipfile " +
                zipfile + ": " + err.what());
    }

    zip_fclose(f);

    index++;
  }

  if (!ErrorIsInvalidArgument(z)) {
    std::string message =
        "failed to get name of file: " + std::string(zip_strerror(z));
    zip_close(z);

    THROW_ERR(message);
  }

  zip_close(z);
}

static void Handle(zip_file_t *f, const char *name, zip_uint64_t size,
                   const std::string &dest_dir) {
  std::string dest_path = ::btool::util::fs::Join(dest_dir, name);
  bool is_dir = dest_path[dest_path.size() - 1] == '/';
  if (is_dir) {
    ::btool::util::fs::MkdirAll(dest_path);
  } else {
    std::ofstream ofs(dest_path);
    zip_uint64_t size_read = 0;
    const std::size_t buf_size = 1024;
    char buf[buf_size];
    while (size_read < size) {
      if (!ofs) {
        THROW_ERR("something is wrong with unzip destination " + dest_path);
      }

      zip_int64_t count = zip_fread(f, buf, buf_size);
      switch (count) {
        case -1:
          THROW_ERR("failed to read uncompressed data: " +
                    std::string(zip_file_strerror(f)));
          break;
        case 0:
          break;
        default:
          size_read += count;
          DEBUGS() << "read " << count << " bytes for " << dest_path
                   << ", total read = " << size_read << std::endl;
          ofs.write(buf, count);
      }
    }
  }
}

static bool ErrorIsInvalidArgument(zip_t *z) {
  zip_error_t *error = zip_get_error(z);
  return (error == nullptr ? false : zip_error_code_zip(error) == ZIP_ER_INVAL);
}

};  // namespace btool::util
