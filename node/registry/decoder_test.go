package registry_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/ankeesler/btool/node/registry"
	"github.com/ankeesler/btool/node/registry/registryfakes"
	"github.com/go-test/deep"
)

func TestDecoderDecode(t *testing.T) {
	aResolver := &nodefakes.FakeResolver{}
	bResolver := &nodefakes.FakeResolver{}
	cResolver := &nodefakes.FakeResolver{}

	newResolverMapper := func() *registryfakes.FakeResolverMapper {
		resolverMapper := &registryfakes.FakeResolverMapper{}
		resolverMapper.MapStub = func(
			name string,
			config map[string]interface{},
		) (node.Resolver, error) {
			switch name {
			case "a":
				return aResolver, nil
			case "b":
				return bResolver, nil
			case "c":
				return cResolver, nil
			default:
				return nil, errors.New("unknown resolver: " + name)
			}
		}
		return resolverMapper
	}

	aN := node.New("a")
	bN := node.New("b")
	cN := node.New("c")
	nN := node.New("n").Dependency(aN, bN, cN)
	nN.Resolver = aResolver

	newNodeMapper := func() *registryfakes.FakeNodeMapper {
		nodeMapper := &registryfakes.FakeNodeMapper{}
		nodeMapper.MapStub = func(name string) (*node.Node, error) {
			switch name {
			case "a":
				return aN, nil
			case "b":
				return bN, nil
			case "c":
				return cN, nil
			default:
				return nil, errors.New("unknown node: " + name)
			}
		}
		return nodeMapper
	}

	config := map[string]interface{}{
		"tuna":   "fish",
		"marlin": []string{"a", "b", "c"},
	}

	data := []struct {
		name    string
		n       *registry.Node
		exN     *node.Node
		exError string
	}{
		{
			name: "Success",
			n: &registry.Node{
				Name: "n",
				Dependencies: []string{
					"a",
					"b",
					"c",
				},
				Resolver: registry.Resolver{
					Name:   "a",
					Config: config,
				},
			},
			exN: nN,
		},
		{
			name: "UnknownNode",
			n: &registry.Node{
				Name: "n",
				Dependencies: []string{
					"a",
					"z",
					"c",
				},
				Resolver: registry.Resolver{
					Name:   "a",
					Config: config,
				},
			},
			exError: "unknown node: z",
		},
		{
			name: "UnknownResolver",
			n: &registry.Node{
				Name: "n",
				Dependencies: []string{
					"a",
					"b",
					"c",
				},
				Resolver: registry.Resolver{
					Name:   "z",
					Config: config,
				},
			},
			exError: "unknown resolver: z",
		},
	}

	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			nodeMapper := newNodeMapper()
			resolverMapper := newResolverMapper()
			d := registry.NewDecoder(nodeMapper, resolverMapper)

			acN, acErr := d.Decode(datum.n)
			if acErr != nil {
				if datum.exError == "" {
					t.Fatal(acErr)
				} else if !strings.HasSuffix(acErr.Error(), datum.exError) {
					t.Fatalf(
						"'%s' does not end with '%s'",
						acErr.Error(),
						datum.exError,
					)
				}
				return
			} else {
				if datum.exError != "" {
					t.Fatalf(
						"expected error ending in '%s'",
						datum.exError,
					)
				}
			}

			if diff := deep.Equal(datum.exN, acN); diff != nil {
				t.Error(diff)
			}

			if ex, ac := 3, nodeMapper.MapCallCount(); ex != ac {
				t.Fatal(ex, "!=", ac)
			}
			if ex, ac := "a", nodeMapper.MapArgsForCall(0); ex != ac {
				t.Error(ex, "!=", ac)
			}
			if ex, ac := "b", nodeMapper.MapArgsForCall(1); ex != ac {
				t.Error(ex, "!=", ac)
			}
			if ex, ac := "c", nodeMapper.MapArgsForCall(2); ex != ac {
				t.Error(ex, "!=", ac)
			}

			if ex, ac := 1, resolverMapper.MapCallCount(); ex != ac {
				t.Fatal(ex, "!=", ac)
			}
		})
	}
}
