#ifndef __BLAH_GRAPH_H__
#define __BLAH_GRAPH_H__

#include <list>
#include <map>
#include <string>

#include "blah.h"

class Blah {
public:
  class Resolver {
  public:
    void Resolve(const std::string& path) = 0;
  };

  Blah(const std::string& path, ) : path_(path) { }

  
private:
  std::string path_;
};

class BlahGraph {
public:
  Blah *AddNode(const std::string& path);
  Blah *AddEdge(const std::string& from, const std::string& to);

private:
  std::map<std::string, Blah*> blahs_;
  std::map<std::string*, std::list<Blah*>> dependencies__;
};

#endif // __BLAH_GRAPH_H__
