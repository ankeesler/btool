// Code generated by counterfeiter. DO NOT EDIT.
package collector0fakes

import (
	"sync"

	collector "github.com/ankeesler/btool/collector0"
)

type FakeProducer struct {
	ProduceStub        func(collector.Store) error
	produceMutex       sync.RWMutex
	produceArgsForCall []struct {
		arg1 collector.Store
	}
	produceReturns struct {
		result1 error
	}
	produceReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeProducer) Produce(arg1 collector.Store) error {
	fake.produceMutex.Lock()
	ret, specificReturn := fake.produceReturnsOnCall[len(fake.produceArgsForCall)]
	fake.produceArgsForCall = append(fake.produceArgsForCall, struct {
		arg1 collector.Store
	}{arg1})
	fake.recordInvocation("Produce", []interface{}{arg1})
	fake.produceMutex.Unlock()
	if fake.ProduceStub != nil {
		return fake.ProduceStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.produceReturns
	return fakeReturns.result1
}

func (fake *FakeProducer) ProduceCallCount() int {
	fake.produceMutex.RLock()
	defer fake.produceMutex.RUnlock()
	return len(fake.produceArgsForCall)
}

func (fake *FakeProducer) ProduceCalls(stub func(collector.Store) error) {
	fake.produceMutex.Lock()
	defer fake.produceMutex.Unlock()
	fake.ProduceStub = stub
}

func (fake *FakeProducer) ProduceArgsForCall(i int) collector.Store {
	fake.produceMutex.RLock()
	defer fake.produceMutex.RUnlock()
	argsForCall := fake.produceArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeProducer) ProduceReturns(result1 error) {
	fake.produceMutex.Lock()
	defer fake.produceMutex.Unlock()
	fake.ProduceStub = nil
	fake.produceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeProducer) ProduceReturnsOnCall(i int, result1 error) {
	fake.produceMutex.Lock()
	defer fake.produceMutex.Unlock()
	fake.ProduceStub = nil
	if fake.produceReturnsOnCall == nil {
		fake.produceReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.produceReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeProducer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.produceMutex.RLock()
	defer fake.produceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeProducer) recordInvocation(key string, args []interface{}) {
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

var _ collector.Producer = new(FakeProducer)
