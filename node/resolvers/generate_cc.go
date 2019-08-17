package resolvers

//import (
//	"bytes"
//	"fmt"
//	"path/filepath"
//	"strings"
//	"unicode"
//
//	"github.com/ankeesler/btool/config"
//	"github.com/pkg/errors"
//	"github.com/sirupsen/logrus"
//	"github.com/spf13/afero"
//)
//
//type Generator struct {
//	fs     afero.Fs
//	config *config.Config
//}
//
//func New(fs afero.Fs, config *config.Config) *Generator {
//	return &Generator{
//		fs:     fs,
//		config: config,
//	}
//}
//
//func (g *Generator) Class(path string) error {
//	if filepath.Ext(path) != "" {
//		return errors.New("do not put extension onto generate path (e.g., some/path/to/file")
//	}
//
//	hFile := filepath.Join(g.config.Root, path+".h")
//	hContent := g.hContent(path)
//	if err := afero.WriteFile(g.fs, hFile, hContent, 0600); err != nil {
//		return errors.Wrap(err, "write "+hFile)
//	}
//	log.Infof("generated %s", hFile)
//
//	cFile := filepath.Join(g.config.Root, path+".cc")
//	cContent := g.cContent(path)
//	if err := afero.WriteFile(g.fs, cFile, cContent, 0600); err != nil {
//		return errors.Wrap(err, "write "+cFile)
//	}
//	log.Infof("generated %s", cFile)
//
//	return nil
//}
//
//func (g *Generator) hContent(path string) []byte {
//	buf := bytes.NewBuffer([]byte{})
//
//	ifndef := g.ifndef(path)
//	namespace := g.namespace(path)
//	class := g.class(path)
//
//	fmt.Fprintf(
//		buf,
//		`#ifndef %s
//#define %s
//
//namespace %s {
//
//class %s {
//
//};
//
//}; // namespace %s
//
//#endif // %s
//`,
//		ifndef,
//		ifndef,
//		namespace,
//		class,
//		namespace,
//		ifndef,
//	)
//
//	return buf.Bytes()
//}
//
//func (g *Generator) cContent(path string) []byte {
//	buf := bytes.NewBuffer([]byte{})
//
//	namespace := g.namespace(path)
//
//	fmt.Fprintf(
//		buf,
//		`#include "%s"
//
//namespace %s {
//
//}; // namespace %s
//`,
//		path+".h",
//		namespace,
//		namespace,
//	)
//
//	return buf.Bytes()
//}
//
//func (g *Generator) ifndef(path string) string {
//	buf := bytes.NewBuffer([]byte{})
//
//	fmt.Fprintf(buf, "%s_", strings.ToUpper(g.cleanProjectName()))
//	fmt.Fprintf(
//		buf,
//		"%s_H_",
//		strings.NewReplacer(
//			"_", "",
//			"/", "_",
//		).Replace(strings.ToUpper(path)),
//	)
//
//	return buf.String()
//}
//
//func (g *Generator) namespace(path string) string {
//	buf := bytes.NewBuffer([]byte{})
//
//	buf.WriteString(g.cleanProjectName())
//
//	dir := filepath.Dir(path)
//	if dir != "." {
//		fmt.Fprintf(buf, "::%s", strings.ReplaceAll(dir, "/", "::"))
//	}
//
//	return buf.String()
//}
//
//func (g *Generator) class(path string) string {
//	base := []byte(filepath.Base(path))
//	class := make([]rune, 0)
//
//	upper := true
//	for i := range base {
//		r := rune(base[i])
//		if r == '_' {
//			upper = true
//		} else if upper {
//			class = append(class, unicode.ToUpper(r))
//			upper = false
//		} else {
//			class = append(class, r)
//		}
//	}
//
//	//from := 0
//	//for {
//	//	to := bytes.IndexByte(base[from:], '_')
//	//	if to == -1 {
//	//		break
//	//	}
//	//
//	//	copy(class, base[from:to])
//	//	from = to
//	//}
//
//	return string(class)
//}
//
//func (g *Generator) cleanProjectName() string {
//	return strings.NewReplacer(
//		"-", "",
//	).Replace(g.config.Name)
//}
