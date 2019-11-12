#include <cstdlib>

#include <iostream>
#include <string>
#include <vector>

#include "app/app.h"
#include "app/builder/builder.h"
#include "app/builder/currenter_impl.h"
#include "app/cleaner/cleaner.h"
#include "app/cleaner/remove_aller_impl.h"
#include "app/collector/cc/exe.h"
#include "app/collector/cc/inc.h"
#include "app/collector/cc/includes_parser_impl.h"
#include "app/collector/cc/obj.h"
#include "app/collector/cc/resolver_factory_impl.h"
#include "app/collector/collector.h"
#include "app/collector/fs/fs_collectini.h"
#include "app/collector/store.h"
#include "app/collector/trivial_collectini.h"
#include "app/lister/lister.h"
#include "app/runner/runner.h"
#include "core/err.h"
#include "core/log.h"
#include "ui/ui.h"
#include "util/flags.h"
#include "util/fs/fs.h"

// workaround for bug-02
#include "app/collector/base_collectini.h"

const static std::string version_string = "0.0.2";

#ifdef __linux__
static const char *compiler_c = "gcc";
static const char *compiler_cc = "g++";
static const char *archiver = "ar";
static const char *linker_c = "gcc";
static const char *linker_cc = "g++";
#elif __APPLE__
static const char *compiler_c = "clang";
static const char *compiler_cc = "clang++";
static const char *archiver = "ar";
static const char *linker_c = "clang";
static const char *linker_cc = "clang++";
#else
#error "unknown platform"
#endif

int main(int argc, const char *argv[]) {
  ::btool::util::Flags f;

  bool version = false;
  f.Bool("version", &version);

  bool debug = false;
  f.Bool("debug", &debug);

  std::string root = ".";
  f.String("root", &root);
  std::string target = "main";
  f.String("target", &target);

  bool clean = false;
  f.Bool("clean", &clean);
  bool list = false;
  f.Bool("list", &list);

  std::string err_s;
  bool success = f.Parse(argc, argv, &err_s);
  if (!success) {
    ERROR("parse flags: %s\n", err_s.c_str());
    return 1;
  }

  if (version) {
    INFO("version %s\n", version_string.c_str());
    return 0;
  }

  ::btool::ui::UI ui;

  ::btool::app::collector::cc::IncludesParserImpl ipi;
  ::btool::app::collector::cc::Inc i(&ipi);

  ::btool::app::collector::cc::ResolverFactoryImpl rfi(
      compiler_c, compiler_cc, archiver, linker_c, linker_cc,
      {"-Wall", "-Werror", "-g", "-O0", "--std=c17"},
      {"-Wall", "-Werror", "-g", "-O0", "--std=c++17"});
  ::btool::app::collector::cc::Obj o(&rfi);
  ::btool::app::collector::cc::Exe e(&rfi);

  ::btool::app::collector::fs::FSCollectini fsc(root.c_str());

  auto root_target = ::btool::util::fs::Join(root, target);
  ::btool::app::collector::TrivialCollectini tc(root_target);

  ::btool::app::collector::Store s;
  ::btool::app::collector::Collector collector(&s);
  collector.AddCollectini(&i);
  collector.AddCollectini(&o);
  collector.AddCollectini(&e);
  collector.AddCollectini(&fsc);
  collector.AddCollectini(&tc);

  ::btool::app::cleaner::RemoveAllerImpl rai;
  ::btool::app::cleaner::Cleaner cleaner(&rai);

  ::btool::app::lister::Lister lister(&std::cout);

  ::btool::app::builder::CurrenterImpl ci;
  ::btool::app::builder::Builder builder(&ci);

  ::btool::app::runner::Runner runner(&ui);

  ::btool::app::App a(&collector, &cleaner, &lister, &builder, &runner);
  auto err = a.Run(root_target, clean, list, false);
  if (err) {
    ERROR("%s\n", err.Msg());
    return 1;
  }

  return 0;
}
