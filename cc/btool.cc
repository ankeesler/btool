#include <cstdlib>

#include <iostream>
#include <string>
#include <vector>

#include "app/app.h"
#include "app/builder/builder.h"
#include "app/builder/currenter_impl.h"
#include "app/cleaner/cleaner.h"
#include "app/cleaner/remove_aller_impl.h"
#include "app/collector/collector.h"
#include "app/lister/lister.h"
#include "app/runner/runner.h"
#include "core/err.h"
#include "core/flags.h"
#include "core/log.h"
#include "ui/ui.h"

int main(int argc, const char *argv[]) {
  ::btool::core::Flags f;

  bool debug = false;
  f.Bool("debug", &debug);

  std::string err_s;
  bool success = f.Parse(argc, argv, &err_s);
  if (!success) {
    ERROR("parse flags: %s\n", err_s.c_str());
    exit(1);
  }
  DEBUG("debug is %d\n", debug);

  ::btool::ui::UI ui;

  ::btool::app::collector::Collector collector;

  ::btool::app::cleaner::RemoveAllerImpl rai;
  ::btool::app::cleaner::Cleaner cleaner(&rai);

  ::btool::app::lister::Lister lister(&std::cout);

  ::btool::app::builder::CurrenterImpl ci;
  ::btool::app::builder::Builder builder(&ci);

  ::btool::app::runner::Runner runner(&ui);

  ::btool::app::App a(&collector, &cleaner, &lister, &builder, &runner);
  auto err = a.Run(false, false, false);
  if (err) {
    ERROR("%s\n", err.Msg());
    return 1;
  }

  return 0;
}
