// Code generated by counterfeiter. DO NOT EDIT.
package resourcefakes

import (
	"context"
	"io"
	"sync"

	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/resource"
	"github.com/concourse/concourse/atc/runtime"
)

type FakeGetter struct {
	GetStub        func(context.Context, runtime.Worker, func(ctx context.Context) (runtime.Container, []runtime.VolumeMount, error), resource.Resource, db.UsedResourceCache, io.Writer) (resource.VersionResult, runtime.ProcessResult, runtime.Volume, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		arg1 context.Context
		arg2 runtime.Worker
		arg3 func(ctx context.Context) (runtime.Container, []runtime.VolumeMount, error)
		arg4 resource.Resource
		arg5 db.UsedResourceCache
		arg6 io.Writer
	}
	getReturns struct {
		result1 resource.VersionResult
		result2 runtime.ProcessResult
		result3 runtime.Volume
		result4 error
	}
	getReturnsOnCall map[int]struct {
		result1 resource.VersionResult
		result2 runtime.ProcessResult
		result3 runtime.Volume
		result4 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeGetter) Get(arg1 context.Context, arg2 runtime.Worker, arg3 func(ctx context.Context) (runtime.Container, []runtime.VolumeMount, error), arg4 resource.Resource, arg5 db.UsedResourceCache, arg6 io.Writer) (resource.VersionResult, runtime.ProcessResult, runtime.Volume, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		arg1 context.Context
		arg2 runtime.Worker
		arg3 func(ctx context.Context) (runtime.Container, []runtime.VolumeMount, error)
		arg4 resource.Resource
		arg5 db.UsedResourceCache
		arg6 io.Writer
	}{arg1, arg2, arg3, arg4, arg5, arg6})
	stub := fake.GetStub
	fakeReturns := fake.getReturns
	fake.recordInvocation("Get", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.getMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3, ret.result4
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3, fakeReturns.result4
}

func (fake *FakeGetter) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeGetter) GetCalls(stub func(context.Context, runtime.Worker, func(ctx context.Context) (runtime.Container, []runtime.VolumeMount, error), resource.Resource, db.UsedResourceCache, io.Writer) (resource.VersionResult, runtime.ProcessResult, runtime.Volume, error)) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = stub
}

func (fake *FakeGetter) GetArgsForCall(i int) (context.Context, runtime.Worker, func(ctx context.Context) (runtime.Container, []runtime.VolumeMount, error), resource.Resource, db.UsedResourceCache, io.Writer) {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	argsForCall := fake.getArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6
}

func (fake *FakeGetter) GetReturns(result1 resource.VersionResult, result2 runtime.ProcessResult, result3 runtime.Volume, result4 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 resource.VersionResult
		result2 runtime.ProcessResult
		result3 runtime.Volume
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *FakeGetter) GetReturnsOnCall(i int, result1 resource.VersionResult, result2 runtime.ProcessResult, result3 runtime.Volume, result4 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 resource.VersionResult
			result2 runtime.ProcessResult
			result3 runtime.Volume
			result4 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 resource.VersionResult
		result2 runtime.ProcessResult
		result3 runtime.Volume
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *FakeGetter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeGetter) recordInvocation(key string, args []interface{}) {
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

var _ resource.Getter = new(FakeGetter)
