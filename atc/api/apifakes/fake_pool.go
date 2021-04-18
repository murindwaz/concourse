// Code generated by counterfeiter. DO NOT EDIT.
package apifakes

import (
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/concourse/atc/api"
	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/runtime"
	"github.com/concourse/concourse/atc/worker"
)

type FakePool struct {
	CreateVolumeForArtifactStub        func(lager.Logger, worker.Spec) (runtime.Volume, db.WorkerArtifact, error)
	createVolumeForArtifactMutex       sync.RWMutex
	createVolumeForArtifactArgsForCall []struct {
		arg1 lager.Logger
		arg2 worker.Spec
	}
	createVolumeForArtifactReturns struct {
		result1 runtime.Volume
		result2 db.WorkerArtifact
		result3 error
	}
	createVolumeForArtifactReturnsOnCall map[int]struct {
		result1 runtime.Volume
		result2 db.WorkerArtifact
		result3 error
	}
	LocateContainerStub        func(lager.Logger, int, string) (runtime.Container, runtime.Worker, bool, error)
	locateContainerMutex       sync.RWMutex
	locateContainerArgsForCall []struct {
		arg1 lager.Logger
		arg2 int
		arg3 string
	}
	locateContainerReturns struct {
		result1 runtime.Container
		result2 runtime.Worker
		result3 bool
		result4 error
	}
	locateContainerReturnsOnCall map[int]struct {
		result1 runtime.Container
		result2 runtime.Worker
		result3 bool
		result4 error
	}
	LocateVolumeStub        func(lager.Logger, int, string) (runtime.Volume, runtime.Worker, bool, error)
	locateVolumeMutex       sync.RWMutex
	locateVolumeArgsForCall []struct {
		arg1 lager.Logger
		arg2 int
		arg3 string
	}
	locateVolumeReturns struct {
		result1 runtime.Volume
		result2 runtime.Worker
		result3 bool
		result4 error
	}
	locateVolumeReturnsOnCall map[int]struct {
		result1 runtime.Volume
		result2 runtime.Worker
		result3 bool
		result4 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePool) CreateVolumeForArtifact(arg1 lager.Logger, arg2 worker.Spec) (runtime.Volume, db.WorkerArtifact, error) {
	fake.createVolumeForArtifactMutex.Lock()
	ret, specificReturn := fake.createVolumeForArtifactReturnsOnCall[len(fake.createVolumeForArtifactArgsForCall)]
	fake.createVolumeForArtifactArgsForCall = append(fake.createVolumeForArtifactArgsForCall, struct {
		arg1 lager.Logger
		arg2 worker.Spec
	}{arg1, arg2})
	stub := fake.CreateVolumeForArtifactStub
	fakeReturns := fake.createVolumeForArtifactReturns
	fake.recordInvocation("CreateVolumeForArtifact", []interface{}{arg1, arg2})
	fake.createVolumeForArtifactMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakePool) CreateVolumeForArtifactCallCount() int {
	fake.createVolumeForArtifactMutex.RLock()
	defer fake.createVolumeForArtifactMutex.RUnlock()
	return len(fake.createVolumeForArtifactArgsForCall)
}

func (fake *FakePool) CreateVolumeForArtifactCalls(stub func(lager.Logger, worker.Spec) (runtime.Volume, db.WorkerArtifact, error)) {
	fake.createVolumeForArtifactMutex.Lock()
	defer fake.createVolumeForArtifactMutex.Unlock()
	fake.CreateVolumeForArtifactStub = stub
}

func (fake *FakePool) CreateVolumeForArtifactArgsForCall(i int) (lager.Logger, worker.Spec) {
	fake.createVolumeForArtifactMutex.RLock()
	defer fake.createVolumeForArtifactMutex.RUnlock()
	argsForCall := fake.createVolumeForArtifactArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePool) CreateVolumeForArtifactReturns(result1 runtime.Volume, result2 db.WorkerArtifact, result3 error) {
	fake.createVolumeForArtifactMutex.Lock()
	defer fake.createVolumeForArtifactMutex.Unlock()
	fake.CreateVolumeForArtifactStub = nil
	fake.createVolumeForArtifactReturns = struct {
		result1 runtime.Volume
		result2 db.WorkerArtifact
		result3 error
	}{result1, result2, result3}
}

func (fake *FakePool) CreateVolumeForArtifactReturnsOnCall(i int, result1 runtime.Volume, result2 db.WorkerArtifact, result3 error) {
	fake.createVolumeForArtifactMutex.Lock()
	defer fake.createVolumeForArtifactMutex.Unlock()
	fake.CreateVolumeForArtifactStub = nil
	if fake.createVolumeForArtifactReturnsOnCall == nil {
		fake.createVolumeForArtifactReturnsOnCall = make(map[int]struct {
			result1 runtime.Volume
			result2 db.WorkerArtifact
			result3 error
		})
	}
	fake.createVolumeForArtifactReturnsOnCall[i] = struct {
		result1 runtime.Volume
		result2 db.WorkerArtifact
		result3 error
	}{result1, result2, result3}
}

func (fake *FakePool) LocateContainer(arg1 lager.Logger, arg2 int, arg3 string) (runtime.Container, runtime.Worker, bool, error) {
	fake.locateContainerMutex.Lock()
	ret, specificReturn := fake.locateContainerReturnsOnCall[len(fake.locateContainerArgsForCall)]
	fake.locateContainerArgsForCall = append(fake.locateContainerArgsForCall, struct {
		arg1 lager.Logger
		arg2 int
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.LocateContainerStub
	fakeReturns := fake.locateContainerReturns
	fake.recordInvocation("LocateContainer", []interface{}{arg1, arg2, arg3})
	fake.locateContainerMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3, ret.result4
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3, fakeReturns.result4
}

func (fake *FakePool) LocateContainerCallCount() int {
	fake.locateContainerMutex.RLock()
	defer fake.locateContainerMutex.RUnlock()
	return len(fake.locateContainerArgsForCall)
}

func (fake *FakePool) LocateContainerCalls(stub func(lager.Logger, int, string) (runtime.Container, runtime.Worker, bool, error)) {
	fake.locateContainerMutex.Lock()
	defer fake.locateContainerMutex.Unlock()
	fake.LocateContainerStub = stub
}

func (fake *FakePool) LocateContainerArgsForCall(i int) (lager.Logger, int, string) {
	fake.locateContainerMutex.RLock()
	defer fake.locateContainerMutex.RUnlock()
	argsForCall := fake.locateContainerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePool) LocateContainerReturns(result1 runtime.Container, result2 runtime.Worker, result3 bool, result4 error) {
	fake.locateContainerMutex.Lock()
	defer fake.locateContainerMutex.Unlock()
	fake.LocateContainerStub = nil
	fake.locateContainerReturns = struct {
		result1 runtime.Container
		result2 runtime.Worker
		result3 bool
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *FakePool) LocateContainerReturnsOnCall(i int, result1 runtime.Container, result2 runtime.Worker, result3 bool, result4 error) {
	fake.locateContainerMutex.Lock()
	defer fake.locateContainerMutex.Unlock()
	fake.LocateContainerStub = nil
	if fake.locateContainerReturnsOnCall == nil {
		fake.locateContainerReturnsOnCall = make(map[int]struct {
			result1 runtime.Container
			result2 runtime.Worker
			result3 bool
			result4 error
		})
	}
	fake.locateContainerReturnsOnCall[i] = struct {
		result1 runtime.Container
		result2 runtime.Worker
		result3 bool
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *FakePool) LocateVolume(arg1 lager.Logger, arg2 int, arg3 string) (runtime.Volume, runtime.Worker, bool, error) {
	fake.locateVolumeMutex.Lock()
	ret, specificReturn := fake.locateVolumeReturnsOnCall[len(fake.locateVolumeArgsForCall)]
	fake.locateVolumeArgsForCall = append(fake.locateVolumeArgsForCall, struct {
		arg1 lager.Logger
		arg2 int
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.LocateVolumeStub
	fakeReturns := fake.locateVolumeReturns
	fake.recordInvocation("LocateVolume", []interface{}{arg1, arg2, arg3})
	fake.locateVolumeMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3, ret.result4
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3, fakeReturns.result4
}

func (fake *FakePool) LocateVolumeCallCount() int {
	fake.locateVolumeMutex.RLock()
	defer fake.locateVolumeMutex.RUnlock()
	return len(fake.locateVolumeArgsForCall)
}

func (fake *FakePool) LocateVolumeCalls(stub func(lager.Logger, int, string) (runtime.Volume, runtime.Worker, bool, error)) {
	fake.locateVolumeMutex.Lock()
	defer fake.locateVolumeMutex.Unlock()
	fake.LocateVolumeStub = stub
}

func (fake *FakePool) LocateVolumeArgsForCall(i int) (lager.Logger, int, string) {
	fake.locateVolumeMutex.RLock()
	defer fake.locateVolumeMutex.RUnlock()
	argsForCall := fake.locateVolumeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePool) LocateVolumeReturns(result1 runtime.Volume, result2 runtime.Worker, result3 bool, result4 error) {
	fake.locateVolumeMutex.Lock()
	defer fake.locateVolumeMutex.Unlock()
	fake.LocateVolumeStub = nil
	fake.locateVolumeReturns = struct {
		result1 runtime.Volume
		result2 runtime.Worker
		result3 bool
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *FakePool) LocateVolumeReturnsOnCall(i int, result1 runtime.Volume, result2 runtime.Worker, result3 bool, result4 error) {
	fake.locateVolumeMutex.Lock()
	defer fake.locateVolumeMutex.Unlock()
	fake.LocateVolumeStub = nil
	if fake.locateVolumeReturnsOnCall == nil {
		fake.locateVolumeReturnsOnCall = make(map[int]struct {
			result1 runtime.Volume
			result2 runtime.Worker
			result3 bool
			result4 error
		})
	}
	fake.locateVolumeReturnsOnCall[i] = struct {
		result1 runtime.Volume
		result2 runtime.Worker
		result3 bool
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *FakePool) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createVolumeForArtifactMutex.RLock()
	defer fake.createVolumeForArtifactMutex.RUnlock()
	fake.locateContainerMutex.RLock()
	defer fake.locateContainerMutex.RUnlock()
	fake.locateVolumeMutex.RLock()
	defer fake.locateVolumeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePool) recordInvocation(key string, args []interface{}) {
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

var _ api.Pool = new(FakePool)
