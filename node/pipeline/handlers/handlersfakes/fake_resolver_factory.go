// Code generated by counterfeiter. DO NOT EDIT.
package handlersfakes

import (
	"sync"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline/handlers"
)

type FakeResolverFactory struct {
	NewArchiveStub        func() node.Resolver
	newArchiveMutex       sync.RWMutex
	newArchiveArgsForCall []struct {
	}
	newArchiveReturns struct {
		result1 node.Resolver
	}
	newArchiveReturnsOnCall map[int]struct {
		result1 node.Resolver
	}
	NewCompileCStub        func([]string) node.Resolver
	newCompileCMutex       sync.RWMutex
	newCompileCArgsForCall []struct {
		arg1 []string
	}
	newCompileCReturns struct {
		result1 node.Resolver
	}
	newCompileCReturnsOnCall map[int]struct {
		result1 node.Resolver
	}
	NewCompileCCStub        func([]string) node.Resolver
	newCompileCCMutex       sync.RWMutex
	newCompileCCArgsForCall []struct {
		arg1 []string
	}
	newCompileCCReturns struct {
		result1 node.Resolver
	}
	newCompileCCReturnsOnCall map[int]struct {
		result1 node.Resolver
	}
	NewDownloadStub        func(string, string) node.Resolver
	newDownloadMutex       sync.RWMutex
	newDownloadArgsForCall []struct {
		arg1 string
		arg2 string
	}
	newDownloadReturns struct {
		result1 node.Resolver
	}
	newDownloadReturnsOnCall map[int]struct {
		result1 node.Resolver
	}
	NewLinkStub        func() node.Resolver
	newLinkMutex       sync.RWMutex
	newLinkArgsForCall []struct {
	}
	newLinkReturns struct {
		result1 node.Resolver
	}
	newLinkReturnsOnCall map[int]struct {
		result1 node.Resolver
	}
	NewSymlinkStub        func() node.Resolver
	newSymlinkMutex       sync.RWMutex
	newSymlinkArgsForCall []struct {
	}
	newSymlinkReturns struct {
		result1 node.Resolver
	}
	newSymlinkReturnsOnCall map[int]struct {
		result1 node.Resolver
	}
	NewUnzipStub        func(string) node.Resolver
	newUnzipMutex       sync.RWMutex
	newUnzipArgsForCall []struct {
		arg1 string
	}
	newUnzipReturns struct {
		result1 node.Resolver
	}
	newUnzipReturnsOnCall map[int]struct {
		result1 node.Resolver
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeResolverFactory) NewArchive() node.Resolver {
	fake.newArchiveMutex.Lock()
	ret, specificReturn := fake.newArchiveReturnsOnCall[len(fake.newArchiveArgsForCall)]
	fake.newArchiveArgsForCall = append(fake.newArchiveArgsForCall, struct {
	}{})
	fake.recordInvocation("NewArchive", []interface{}{})
	fake.newArchiveMutex.Unlock()
	if fake.NewArchiveStub != nil {
		return fake.NewArchiveStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newArchiveReturns
	return fakeReturns.result1
}

func (fake *FakeResolverFactory) NewArchiveCallCount() int {
	fake.newArchiveMutex.RLock()
	defer fake.newArchiveMutex.RUnlock()
	return len(fake.newArchiveArgsForCall)
}

func (fake *FakeResolverFactory) NewArchiveCalls(stub func() node.Resolver) {
	fake.newArchiveMutex.Lock()
	defer fake.newArchiveMutex.Unlock()
	fake.NewArchiveStub = stub
}

func (fake *FakeResolverFactory) NewArchiveReturns(result1 node.Resolver) {
	fake.newArchiveMutex.Lock()
	defer fake.newArchiveMutex.Unlock()
	fake.NewArchiveStub = nil
	fake.newArchiveReturns = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewArchiveReturnsOnCall(i int, result1 node.Resolver) {
	fake.newArchiveMutex.Lock()
	defer fake.newArchiveMutex.Unlock()
	fake.NewArchiveStub = nil
	if fake.newArchiveReturnsOnCall == nil {
		fake.newArchiveReturnsOnCall = make(map[int]struct {
			result1 node.Resolver
		})
	}
	fake.newArchiveReturnsOnCall[i] = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewCompileC(arg1 []string) node.Resolver {
	var arg1Copy []string
	if arg1 != nil {
		arg1Copy = make([]string, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.newCompileCMutex.Lock()
	ret, specificReturn := fake.newCompileCReturnsOnCall[len(fake.newCompileCArgsForCall)]
	fake.newCompileCArgsForCall = append(fake.newCompileCArgsForCall, struct {
		arg1 []string
	}{arg1Copy})
	fake.recordInvocation("NewCompileC", []interface{}{arg1Copy})
	fake.newCompileCMutex.Unlock()
	if fake.NewCompileCStub != nil {
		return fake.NewCompileCStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newCompileCReturns
	return fakeReturns.result1
}

func (fake *FakeResolverFactory) NewCompileCCallCount() int {
	fake.newCompileCMutex.RLock()
	defer fake.newCompileCMutex.RUnlock()
	return len(fake.newCompileCArgsForCall)
}

func (fake *FakeResolverFactory) NewCompileCCalls(stub func([]string) node.Resolver) {
	fake.newCompileCMutex.Lock()
	defer fake.newCompileCMutex.Unlock()
	fake.NewCompileCStub = stub
}

func (fake *FakeResolverFactory) NewCompileCArgsForCall(i int) []string {
	fake.newCompileCMutex.RLock()
	defer fake.newCompileCMutex.RUnlock()
	argsForCall := fake.newCompileCArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeResolverFactory) NewCompileCReturns(result1 node.Resolver) {
	fake.newCompileCMutex.Lock()
	defer fake.newCompileCMutex.Unlock()
	fake.NewCompileCStub = nil
	fake.newCompileCReturns = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewCompileCReturnsOnCall(i int, result1 node.Resolver) {
	fake.newCompileCMutex.Lock()
	defer fake.newCompileCMutex.Unlock()
	fake.NewCompileCStub = nil
	if fake.newCompileCReturnsOnCall == nil {
		fake.newCompileCReturnsOnCall = make(map[int]struct {
			result1 node.Resolver
		})
	}
	fake.newCompileCReturnsOnCall[i] = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewCompileCC(arg1 []string) node.Resolver {
	var arg1Copy []string
	if arg1 != nil {
		arg1Copy = make([]string, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.newCompileCCMutex.Lock()
	ret, specificReturn := fake.newCompileCCReturnsOnCall[len(fake.newCompileCCArgsForCall)]
	fake.newCompileCCArgsForCall = append(fake.newCompileCCArgsForCall, struct {
		arg1 []string
	}{arg1Copy})
	fake.recordInvocation("NewCompileCC", []interface{}{arg1Copy})
	fake.newCompileCCMutex.Unlock()
	if fake.NewCompileCCStub != nil {
		return fake.NewCompileCCStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newCompileCCReturns
	return fakeReturns.result1
}

func (fake *FakeResolverFactory) NewCompileCCCallCount() int {
	fake.newCompileCCMutex.RLock()
	defer fake.newCompileCCMutex.RUnlock()
	return len(fake.newCompileCCArgsForCall)
}

func (fake *FakeResolverFactory) NewCompileCCCalls(stub func([]string) node.Resolver) {
	fake.newCompileCCMutex.Lock()
	defer fake.newCompileCCMutex.Unlock()
	fake.NewCompileCCStub = stub
}

func (fake *FakeResolverFactory) NewCompileCCArgsForCall(i int) []string {
	fake.newCompileCCMutex.RLock()
	defer fake.newCompileCCMutex.RUnlock()
	argsForCall := fake.newCompileCCArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeResolverFactory) NewCompileCCReturns(result1 node.Resolver) {
	fake.newCompileCCMutex.Lock()
	defer fake.newCompileCCMutex.Unlock()
	fake.NewCompileCCStub = nil
	fake.newCompileCCReturns = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewCompileCCReturnsOnCall(i int, result1 node.Resolver) {
	fake.newCompileCCMutex.Lock()
	defer fake.newCompileCCMutex.Unlock()
	fake.NewCompileCCStub = nil
	if fake.newCompileCCReturnsOnCall == nil {
		fake.newCompileCCReturnsOnCall = make(map[int]struct {
			result1 node.Resolver
		})
	}
	fake.newCompileCCReturnsOnCall[i] = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewDownload(arg1 string, arg2 string) node.Resolver {
	fake.newDownloadMutex.Lock()
	ret, specificReturn := fake.newDownloadReturnsOnCall[len(fake.newDownloadArgsForCall)]
	fake.newDownloadArgsForCall = append(fake.newDownloadArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("NewDownload", []interface{}{arg1, arg2})
	fake.newDownloadMutex.Unlock()
	if fake.NewDownloadStub != nil {
		return fake.NewDownloadStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newDownloadReturns
	return fakeReturns.result1
}

func (fake *FakeResolverFactory) NewDownloadCallCount() int {
	fake.newDownloadMutex.RLock()
	defer fake.newDownloadMutex.RUnlock()
	return len(fake.newDownloadArgsForCall)
}

func (fake *FakeResolverFactory) NewDownloadCalls(stub func(string, string) node.Resolver) {
	fake.newDownloadMutex.Lock()
	defer fake.newDownloadMutex.Unlock()
	fake.NewDownloadStub = stub
}

func (fake *FakeResolverFactory) NewDownloadArgsForCall(i int) (string, string) {
	fake.newDownloadMutex.RLock()
	defer fake.newDownloadMutex.RUnlock()
	argsForCall := fake.newDownloadArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeResolverFactory) NewDownloadReturns(result1 node.Resolver) {
	fake.newDownloadMutex.Lock()
	defer fake.newDownloadMutex.Unlock()
	fake.NewDownloadStub = nil
	fake.newDownloadReturns = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewDownloadReturnsOnCall(i int, result1 node.Resolver) {
	fake.newDownloadMutex.Lock()
	defer fake.newDownloadMutex.Unlock()
	fake.NewDownloadStub = nil
	if fake.newDownloadReturnsOnCall == nil {
		fake.newDownloadReturnsOnCall = make(map[int]struct {
			result1 node.Resolver
		})
	}
	fake.newDownloadReturnsOnCall[i] = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewLink() node.Resolver {
	fake.newLinkMutex.Lock()
	ret, specificReturn := fake.newLinkReturnsOnCall[len(fake.newLinkArgsForCall)]
	fake.newLinkArgsForCall = append(fake.newLinkArgsForCall, struct {
	}{})
	fake.recordInvocation("NewLink", []interface{}{})
	fake.newLinkMutex.Unlock()
	if fake.NewLinkStub != nil {
		return fake.NewLinkStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newLinkReturns
	return fakeReturns.result1
}

func (fake *FakeResolverFactory) NewLinkCallCount() int {
	fake.newLinkMutex.RLock()
	defer fake.newLinkMutex.RUnlock()
	return len(fake.newLinkArgsForCall)
}

func (fake *FakeResolverFactory) NewLinkCalls(stub func() node.Resolver) {
	fake.newLinkMutex.Lock()
	defer fake.newLinkMutex.Unlock()
	fake.NewLinkStub = stub
}

func (fake *FakeResolverFactory) NewLinkReturns(result1 node.Resolver) {
	fake.newLinkMutex.Lock()
	defer fake.newLinkMutex.Unlock()
	fake.NewLinkStub = nil
	fake.newLinkReturns = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewLinkReturnsOnCall(i int, result1 node.Resolver) {
	fake.newLinkMutex.Lock()
	defer fake.newLinkMutex.Unlock()
	fake.NewLinkStub = nil
	if fake.newLinkReturnsOnCall == nil {
		fake.newLinkReturnsOnCall = make(map[int]struct {
			result1 node.Resolver
		})
	}
	fake.newLinkReturnsOnCall[i] = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewSymlink() node.Resolver {
	fake.newSymlinkMutex.Lock()
	ret, specificReturn := fake.newSymlinkReturnsOnCall[len(fake.newSymlinkArgsForCall)]
	fake.newSymlinkArgsForCall = append(fake.newSymlinkArgsForCall, struct {
	}{})
	fake.recordInvocation("NewSymlink", []interface{}{})
	fake.newSymlinkMutex.Unlock()
	if fake.NewSymlinkStub != nil {
		return fake.NewSymlinkStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newSymlinkReturns
	return fakeReturns.result1
}

func (fake *FakeResolverFactory) NewSymlinkCallCount() int {
	fake.newSymlinkMutex.RLock()
	defer fake.newSymlinkMutex.RUnlock()
	return len(fake.newSymlinkArgsForCall)
}

func (fake *FakeResolverFactory) NewSymlinkCalls(stub func() node.Resolver) {
	fake.newSymlinkMutex.Lock()
	defer fake.newSymlinkMutex.Unlock()
	fake.NewSymlinkStub = stub
}

func (fake *FakeResolverFactory) NewSymlinkReturns(result1 node.Resolver) {
	fake.newSymlinkMutex.Lock()
	defer fake.newSymlinkMutex.Unlock()
	fake.NewSymlinkStub = nil
	fake.newSymlinkReturns = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewSymlinkReturnsOnCall(i int, result1 node.Resolver) {
	fake.newSymlinkMutex.Lock()
	defer fake.newSymlinkMutex.Unlock()
	fake.NewSymlinkStub = nil
	if fake.newSymlinkReturnsOnCall == nil {
		fake.newSymlinkReturnsOnCall = make(map[int]struct {
			result1 node.Resolver
		})
	}
	fake.newSymlinkReturnsOnCall[i] = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewUnzip(arg1 string) node.Resolver {
	fake.newUnzipMutex.Lock()
	ret, specificReturn := fake.newUnzipReturnsOnCall[len(fake.newUnzipArgsForCall)]
	fake.newUnzipArgsForCall = append(fake.newUnzipArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("NewUnzip", []interface{}{arg1})
	fake.newUnzipMutex.Unlock()
	if fake.NewUnzipStub != nil {
		return fake.NewUnzipStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newUnzipReturns
	return fakeReturns.result1
}

func (fake *FakeResolverFactory) NewUnzipCallCount() int {
	fake.newUnzipMutex.RLock()
	defer fake.newUnzipMutex.RUnlock()
	return len(fake.newUnzipArgsForCall)
}

func (fake *FakeResolverFactory) NewUnzipCalls(stub func(string) node.Resolver) {
	fake.newUnzipMutex.Lock()
	defer fake.newUnzipMutex.Unlock()
	fake.NewUnzipStub = stub
}

func (fake *FakeResolverFactory) NewUnzipArgsForCall(i int) string {
	fake.newUnzipMutex.RLock()
	defer fake.newUnzipMutex.RUnlock()
	argsForCall := fake.newUnzipArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeResolverFactory) NewUnzipReturns(result1 node.Resolver) {
	fake.newUnzipMutex.Lock()
	defer fake.newUnzipMutex.Unlock()
	fake.NewUnzipStub = nil
	fake.newUnzipReturns = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) NewUnzipReturnsOnCall(i int, result1 node.Resolver) {
	fake.newUnzipMutex.Lock()
	defer fake.newUnzipMutex.Unlock()
	fake.NewUnzipStub = nil
	if fake.newUnzipReturnsOnCall == nil {
		fake.newUnzipReturnsOnCall = make(map[int]struct {
			result1 node.Resolver
		})
	}
	fake.newUnzipReturnsOnCall[i] = struct {
		result1 node.Resolver
	}{result1}
}

func (fake *FakeResolverFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newArchiveMutex.RLock()
	defer fake.newArchiveMutex.RUnlock()
	fake.newCompileCMutex.RLock()
	defer fake.newCompileCMutex.RUnlock()
	fake.newCompileCCMutex.RLock()
	defer fake.newCompileCCMutex.RUnlock()
	fake.newDownloadMutex.RLock()
	defer fake.newDownloadMutex.RUnlock()
	fake.newLinkMutex.RLock()
	defer fake.newLinkMutex.RUnlock()
	fake.newSymlinkMutex.RLock()
	defer fake.newSymlinkMutex.RUnlock()
	fake.newUnzipMutex.RLock()
	defer fake.newUnzipMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeResolverFactory) recordInvocation(key string, args []interface{}) {
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

var _ handlers.ResolverFactory = new(FakeResolverFactory)
