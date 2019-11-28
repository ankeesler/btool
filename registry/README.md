# btool registry data

This directory contains the btool registry data that is served from the btool
registry at btoolregistry.cfapps.io.

To see an index of this data, run `curl btoolregistry.cfapps.io`.

## Format

The format of these files can be found in `btool::app::collector::registry`
package. Specifically, `registry.h` and `yaml_serializer.cc` will be helpful.
