# btool registry data

This directory contains the btool registry data that is served from the btool
registry at `ankeesler.github.io/btool/index.yml`.

To see an index of this data, run `curl ankeesler.github.io/btool/index.yml`.

## Format

The format of these files can be found in `btool::app::collector::registry`
package. Specifically, `registry.h` and `yaml_serializer.cc` will be helpful.
