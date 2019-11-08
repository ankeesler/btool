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
#include "app/collector/collector.h"
#include "app/collector/fs/fs_collectini.h"
#include "app/collector/resolver_factory_impl.h"
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

int main(int argc, const char *argv[]) {
  ::btool::util::Flags f;

  bool debug = false;
  f.Bool("debug", &debug);

  std::string root = ".";
  f.String("root", &root);
  std::string target = "main";
  f.String("target", &target);

  bool list = false;
  f.Bool("list", &list);

  std::string err_s;
  bool success = f.Parse(argc, argv, &err_s);
  if (!success) {
    ERROR("parse flags: %s\n", err_s.c_str());
    exit(1);
  }

  ::btool::ui::UI ui;

  ::btool::app::collector::cc::IncludesParserImpl ipi;
  ::btool::app::collector::cc::Inc i(&ipi);

  ::btool::app::collector::ResolverFactoryImpl rfi;
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
  auto err = a.Run(root_target, false, list, false);
  if (err) {
    ERROR("%s\n", err.Msg());
    return 1;
  }

  return 0;
}
