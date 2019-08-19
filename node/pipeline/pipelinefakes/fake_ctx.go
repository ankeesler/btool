// Code generated by counterfeiter. DO NOT EDIT.
package pipelinefakes

import (
	"sync"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
)

type FakeCtx struct {
	AddStub        func(*node.Node)
	addMutex       sync.RWMutex
	addArgsForCall []struct {
		arg1 *node.Node
	}
	AllStub        func() []*node.Node
	allMutex       sync.RWMutex
	allArgsForCall []struct {
	}
	allReturns struct {
		result1 []*node.Node
	}
	allReturnsOnCall map[int]struct {
		result1 []*node.Node
	}
	FindStub        func(string) *node.Node
	findMutex       sync.RWMutex
	findArgsForCall []struct {
		arg1 string
	}
	findReturns struct {
		result1 *node.Node
	}
	findReturnsOnCall map[int]struct {
		result1 *node.Node
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCtx) Add(arg1 *node.Node) {
	fake.addMutex.Lock()
	fake.addArgsForCall = append(fake.addArgsForCall, struct {
		arg1 *node.Node
	}{arg1})
	fake.recordInvocation("Add", []interface{}{arg1})
	fake.addMutex.Unlock()
	if fake.AddStub != nil {
		fake.AddStub(arg1)
	}
}

func (fake *FakeCtx) AddCallCount() int {
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	return len(fake.addArgsForCall)
}

func (fake *FakeCtx) AddCalls(stub func(*node.Node)) {
	fake.addMutex.Lock()
	defer fake.addMutex.Unlock()
	fake.AddStub = stub
}

func (fake *FakeCtx) AddArgsForCall(i int) *node.Node {
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	argsForCall := fake.addArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCtx) All() []*node.Node {
	fake.allMutex.Lock()
	ret, specificReturn := fake.allReturnsOnCall[len(fake.allArgsForCall)]
	fake.allArgsForCall = append(fake.allArgsForCall, struct {
	}{})
	fake.recordInvocation("All", []interface{}{})
	fake.allMutex.Unlock()
	if fake.AllStub != nil {
		return fake.AllStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.allReturns
	return fakeReturns.result1
}

func (fake *FakeCtx) AllCallCount() int {
	fake.allMutex.RLock()
	defer fake.allMutex.RUnlock()
	return len(fake.allArgsForCall)
}

func (fake *FakeCtx) AllCalls(stub func() []*node.Node) {
	fake.allMutex.Lock()
	defer fake.allMutex.Unlock()
	fake.AllStub = stub
}

func (fake *FakeCtx) AllReturns(result1 []*node.Node) {
	fake.allMutex.Lock()
	defer fake.allMutex.Unlock()
	fake.AllStub = nil
	fake.allReturns = struct {
		result1 []*node.Node
	}{result1}
}

func (fake *FakeCtx) AllReturnsOnCall(i int, result1 []*node.Node) {
	fake.allMutex.Lock()
	defer fake.allMutex.Unlock()
	fake.AllStub = nil
	if fake.allReturnsOnCall == nil {
		fake.allReturnsOnCall = make(map[int]struct {
			result1 []*node.Node
		})
	}
	fake.allReturnsOnCall[i] = struct {
		result1 []*node.Node
	}{result1}
}

func (fake *FakeCtx) Find(arg1 string) *node.Node {
	fake.findMutex.Lock()
	ret, specificReturn := fake.findReturnsOnCall[len(fake.findArgsForCall)]
	fake.findArgsForCall = append(fake.findArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Find", []interface{}{arg1})
	fake.findMutex.Unlock()
	if fake.FindStub != nil {
		return fake.FindStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.findReturns
	return fakeReturns.result1
}

func (fake *FakeCtx) FindCallCount() int {
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	return len(fake.findArgsForCall)
}

func (fake *FakeCtx) FindCalls(stub func(string) *node.Node) {
	fake.findMutex.Lock()
	defer fake.findMutex.Unlock()
	fake.FindStub = stub
}

func (fake *FakeCtx) FindArgsForCall(i int) string {
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	argsForCall := fake.findArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCtx) FindReturns(result1 *node.Node) {
	fake.findMutex.Lock()
	defer fake.findMutex.Unlock()
	fake.FindStub = nil
	fake.findReturns = struct {
		result1 *node.Node
	}{result1}
}

func (fake *FakeCtx) FindReturnsOnCall(i int, result1 *node.Node) {
	fake.findMutex.Lock()
	defer fake.findMutex.Unlock()
	fake.FindStub = nil
	if fake.findReturnsOnCall == nil {
		fake.findReturnsOnCall = make(map[int]struct {
			result1 *node.Node
		})
	}
	fake.findReturnsOnCall[i] = struct {
		result1 *node.Node
	}{result1}
}

func (fake *FakeCtx) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	fake.allMutex.RLock()
	defer fake.allMutex.RUnlock()
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCtx) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ pipeline.Ctx = new(FakeCtx)
