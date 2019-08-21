// Code generated by counterfeiter. DO NOT EDIT.
package cleanerfakes

import (
	"sync"

	"github.com/ankeesler/btool/cleaner"
	"github.com/ankeesler/btool/node"
)

type FakeCallback struct {
	OnCleanStub        func(*node.Node)
	onCleanMutex       sync.RWMutex
	onCleanArgsForCall []struct {
		arg1 *node.Node
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCallback) OnClean(arg1 *node.Node) {
	fake.onCleanMutex.Lock()
	fake.onCleanArgsForCall = append(fake.onCleanArgsForCall, struct {
		arg1 *node.Node
	}{arg1})
	fake.recordInvocation("OnClean", []interface{}{arg1})
	fake.onCleanMutex.Unlock()
	if fake.OnCleanStub != nil {
		fake.OnCleanStub(arg1)
	}
}

func (fake *FakeCallback) OnCleanCallCount() int {
	fake.onCleanMutex.RLock()
	defer fake.onCleanMutex.RUnlock()
	return len(fake.onCleanArgsForCall)
}

func (fake *FakeCallback) OnCleanCalls(stub func(*node.Node)) {
	fake.onCleanMutex.Lock()
	defer fake.onCleanMutex.Unlock()
	fake.OnCleanStub = stub
}

func (fake *FakeCallback) OnCleanArgsForCall(i int) *node.Node {
	fake.onCleanMutex.RLock()
	defer fake.onCleanMutex.RUnlock()
	argsForCall := fake.onCleanArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCallback) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.onCleanMutex.RLock()
	defer fake.onCleanMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCallback) recordInvocation(key string, args []interface{}) {
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

var _ cleaner.Callback = new(FakeCallback)
