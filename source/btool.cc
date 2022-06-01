#include <cstdlib>
#include <iostream>
#include <stdexcept>
#include <string>
#include <vector>

#include "app/app.h"
#include "app/builder/currenter_impl.h"
#include "app/builder/parallel_builder.h"
#include "app/builder/work_pool_impl.h"
#include "app/cleaner/cleaner.h"
#include "app/cleaner/remove_aller_impl.h"
#include "app/collector/cc/exe.h"
#include "app/collector/cc/inc.h"
#include "app/collector/cc/includes_parser_impl.h"
#include "app/collector/cc/obj.h"
#include "app/collector/cc/resolver_factory_delegate.h"
#include "app/collector/cc/resolver_factory_impl.h"
#include "app/collector/collector.h"
#include "app/collector/fs/fs_collectini.h"
#include "app/collector/registry/fs_registry.h"
#include "app/collector/registry/gaggle_collector_impl.h"
#include "app/collector/registry/http_registry.h"
#include "app/collector/registry/registry.h"
#include "app/collector/registry/registry_collectini.h"
#include "app/collector/registry/resolver_factory_delegate.h"
#include "app/collector/registry/resolver_factory_impl.h"
#include "app/collector/registry/yaml_file_cache.h"
#include "app/collector/registry/yaml_serializer.h"
#include "app/collector/store.h"
#include "app/collector/trivial_collectini.h"
#include "app/lister/lister.h"
#include "app/runner/runner.h"
#include "err.h"
#include "log.h"
#include "ui/ui.h"
#include "util/flags.h"
#include "util/fs/fs.h"
#include "util/string/string.h"

const static std::string version_string = "0.8";

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

static std::string GetDefaultCache();

int main(int argc, const char *argv[]) {
  ::btool::Log::SetCurrentLevel(::btool::Log::kInfo);

  ::btool::util::Flags f;

  bool version = false;
  f.Bool("version", "Print btool version and exit", &version);

  bool help = false;
  f.Bool("help", "Print btool usage", &help);

  std::string loglevel = "info";
  f.String("loglevel",
           "Select logging verbosity (debug, info, error) (default: info)",
           &loglevel);

  std::string root = ".";
  f.String("root", "Specify project root", &root);
  std::string target = "";
  f.String("target", "Specify build target", &target);
  std::string cache = GetDefaultCache();
  f.String("cache", "Specify build cache (default: " + cache + ")", &cache);

  std::string registry = "https://btoolregistry.cfapps.io";
  f.String("registry", "Specify registry URI (default: " + registry + ")",
           &registry);
  int registry_cache_timeout_s = 1 * 60 * 60;  // 1 hr
  f.Int("registrycachetimeout",
        "Specify registry cache timeout (seconds); set to 0 to disable caching",
        &registry_cache_timeout_s);

  int threads = 1;
  f.Int("threads",
        "Specify number of threads to use when building (default: 1)",
        &threads);

  bool clean = false;
  f.Bool("clean", "Perform clean of target graph", &clean);
  bool list = false;
  f.Bool("list", "Show target graph", &list);
  bool run = false;
  f.Bool("run", "Execute target after building", &run);

  std::string err_s;
  bool success = f.Parse(argc, argv, &err_s);
  if (!success) {
    ERROR("parse flags: %s\n", err_s.c_str());
    return 1;
  }

  ::btool::Log::Level l = ::btool::Log::ParseLevel(loglevel);
  if (l == ::btool::Log::kUnknown) {
    ERRORS() << "couldn't parse log level: " << loglevel << std::endl;
    return 1;
  }
  ::btool::Log::SetCurrentLevel(l);

  if (version) {
    INFO("version %s\n", version_string.c_str());
    return 0;
  }

  if (help || target.empty()) {
    std::cout << "btool" << std::endl;
    f.Usage(&std::cout);
    return 1;
  }

  ::btool::ui::UI ui(cache);

  ::btool::app::collector::Store s;

  ::btool::app::collector::registry::YamlSerializer<
      ::btool::app::collector::registry::Index>
      ys_i;
  ::btool::app::collector::registry::YamlSerializer<
      ::btool::app::collector::registry::Gaggle>
      ys_g;
  ::btool::app::collector::registry::FsRegistry fr(registry, &ys_i, &ys_g);
  ::btool::app::collector::registry::HttpRegistry hr(registry, &ys_i, &ys_g);
  ::btool::app::collector::registry::Registry *r;
  if (::btool::util::string::HasPrefix(registry, "http")) {
    r = &hr;
  } else {
    r = &fr;
  }

  ::btool::app::collector::registry::ResolverFactoryImpl r_rfi;
  ::btool::app::collector::registry::ResolverFactoryDelegate r_rfd(&r_rfi);

  ::btool::app::collector::cc::ResolverFactoryImpl c_rfi(
      compiler_c, compiler_cc, archiver, linker_c, linker_cc,
      {"-Wall", "-Werror", "-g", "-O0", "--std=c17"},
      {"-Wall", "-Werror", "-g", "-O0", "--std=c++17"}, {}, {});
  ::btool::app::collector::cc::ResolverFactoryDelegate c_rfd(&c_rfi);

  ::btool::app::collector::registry::GaggleCollectorImpl gci;
  gci.AddResolverFactoryDelegate(&r_rfd);
  gci.AddResolverFactoryDelegate(&c_rfd);

  ::btool::app::collector::registry::YamlFileCache<
      ::btool::app::collector::registry::Index>
      yfc_i(&ys_i, cache, std::chrono::seconds(registry_cache_timeout_s));
  ::btool::app::collector::registry::YamlFileCache<
      ::btool::app::collector::registry::Gaggle>
      yfc_g(&ys_g, cache, std::chrono::seconds(registry_cache_timeout_s));
  ::btool::app::collector::registry::RegistryCollectini rc(r, cache, &yfc_i,
                                                           &yfc_g, &gci);

  ::btool::app::collector::cc::IncludesParserImpl ipi;
  ::btool::app::collector::cc::Inc i(&ipi);

  ::btool::app::collector::cc::Obj o(&c_rfi);
  ::btool::app::collector::cc::Exe e(&c_rfi);

  ::btool::app::collector::fs::FSCollectini fsc(root.c_str());

  auto root_target = ::btool::util::fs::Join(root, target);
  ::btool::app::collector::TrivialCollectini tc(root_target);

  ::btool::app::collector::Collector collector(&s);
  collector.AddCollectini(&rc);
  collector.AddCollectini(&i);
  collector.AddCollectini(&o);
  collector.AddCollectini(&e);
  collector.AddCollectini(&fsc);
  collector.AddCollectini(&tc);

  ::btool::app::cleaner::RemoveAllerImpl rai;
  ::btool::app::cleaner::Cleaner cleaner(&rai);

  ::btool::app::lister::Lister lister(&std::cout);

  ::btool::app::builder::CurrenterImpl ci;
  ::btool::app::builder::WorkPoolImpl wpi(threads);
  ::btool::app::builder::ParallelBuilder parallel_builder(&wpi, &ci, &ui);

  ::btool::app::runner::Runner runner(&ui);

  ::btool::app::App a(&collector, &cleaner, &lister, &parallel_builder,
                      &runner);
  try {
    a.Run(root_target, clean, list, run);
  } catch (const std::exception &e) {
    ERRORS() << e.what() << std::endl;
    return 1;
  }

  return 0;
}

static std::string GetDefaultCache() {
  auto home = ::getenv("HOME");
  return (home == nullptr ? ".btool" : ::btool::util::fs::Join(home, ".btool"));
}
