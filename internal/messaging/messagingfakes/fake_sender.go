// Code generated by counterfeiter. DO NOT EDIT.
package messagingfakes

import (
	"sync"

	"github.com/taxibeat/pigeon/internal/messaging"
)

type FakeSender struct {
	SendStub        func(messaging.Message) error
	sendMutex       sync.RWMutex
	sendArgsForCall []struct {
		arg1 messaging.Message
	}
	sendReturns struct {
		result1 error
	}
	sendReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSender) Send(arg1 messaging.Message) error {
	fake.sendMutex.Lock()
	ret, specificReturn := fake.sendReturnsOnCall[len(fake.sendArgsForCall)]
	fake.sendArgsForCall = append(fake.sendArgsForCall, struct {
		arg1 messaging.Message
	}{arg1})
	fake.recordInvocation("Send", []interface{}{arg1})
	fake.sendMutex.Unlock()
	if fake.SendStub != nil {
		return fake.SendStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sendReturns
	return fakeReturns.result1
}

func (fake *FakeSender) SendCallCount() int {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return len(fake.sendArgsForCall)
}

func (fake *FakeSender) SendCalls(stub func(messaging.Message) error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = stub
}

func (fake *FakeSender) SendArgsForCall(i int) messaging.Message {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	argsForCall := fake.sendArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeSender) SendReturns(result1 error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = nil
	fake.sendReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSender) SendReturnsOnCall(i int, result1 error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = nil
	if fake.sendReturnsOnCall == nil {
		fake.sendReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.sendReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeSender) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeSender) recordInvocation(key string, args []interface{}) {
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

var _ messaging.Sender = new(FakeSender)