// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"sync"

	"fabricbypeer/bccsp/idemix/handlers"
)

type Issuer struct {
	NewKeyStub        func([]string) (handlers.IssuerSecretKey, error)
	newKeyMutex       sync.RWMutex
	newKeyArgsForCall []struct {
		arg1 []string
	}
	newKeyReturns struct {
		result1 handlers.IssuerSecretKey
		result2 error
	}
	newKeyReturnsOnCall map[int]struct {
		result1 handlers.IssuerSecretKey
		result2 error
	}
	NewPublicKeyFromBytesStub        func([]byte, []string) (handlers.IssuerPublicKey, error)
	newPublicKeyFromBytesMutex       sync.RWMutex
	newPublicKeyFromBytesArgsForCall []struct {
		arg1 []byte
		arg2 []string
	}
	newPublicKeyFromBytesReturns struct {
		result1 handlers.IssuerPublicKey
		result2 error
	}
	newPublicKeyFromBytesReturnsOnCall map[int]struct {
		result1 handlers.IssuerPublicKey
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Issuer) NewKey(arg1 []string) (handlers.IssuerSecretKey, error) {
	var arg1Copy []string
	if arg1 != nil {
		arg1Copy = make([]string, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.newKeyMutex.Lock()
	ret, specificReturn := fake.newKeyReturnsOnCall[len(fake.newKeyArgsForCall)]
	fake.newKeyArgsForCall = append(fake.newKeyArgsForCall, struct {
		arg1 []string
	}{arg1Copy})
	fake.recordInvocation("NewKey", []interface{}{arg1Copy})
	fake.newKeyMutex.Unlock()
	if fake.NewKeyStub != nil {
		return fake.NewKeyStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.newKeyReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Issuer) NewKeyCallCount() int {
	fake.newKeyMutex.RLock()
	defer fake.newKeyMutex.RUnlock()
	return len(fake.newKeyArgsForCall)
}

func (fake *Issuer) NewKeyCalls(stub func([]string) (handlers.IssuerSecretKey, error)) {
	fake.newKeyMutex.Lock()
	defer fake.newKeyMutex.Unlock()
	fake.NewKeyStub = stub
}

func (fake *Issuer) NewKeyArgsForCall(i int) []string {
	fake.newKeyMutex.RLock()
	defer fake.newKeyMutex.RUnlock()
	argsForCall := fake.newKeyArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Issuer) NewKeyReturns(result1 handlers.IssuerSecretKey, result2 error) {
	fake.newKeyMutex.Lock()
	defer fake.newKeyMutex.Unlock()
	fake.NewKeyStub = nil
	fake.newKeyReturns = struct {
		result1 handlers.IssuerSecretKey
		result2 error
	}{result1, result2}
}

func (fake *Issuer) NewKeyReturnsOnCall(i int, result1 handlers.IssuerSecretKey, result2 error) {
	fake.newKeyMutex.Lock()
	defer fake.newKeyMutex.Unlock()
	fake.NewKeyStub = nil
	if fake.newKeyReturnsOnCall == nil {
		fake.newKeyReturnsOnCall = make(map[int]struct {
			result1 handlers.IssuerSecretKey
			result2 error
		})
	}
	fake.newKeyReturnsOnCall[i] = struct {
		result1 handlers.IssuerSecretKey
		result2 error
	}{result1, result2}
}

func (fake *Issuer) NewPublicKeyFromBytes(arg1 []byte, arg2 []string) (handlers.IssuerPublicKey, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	var arg2Copy []string
	if arg2 != nil {
		arg2Copy = make([]string, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.newPublicKeyFromBytesMutex.Lock()
	ret, specificReturn := fake.newPublicKeyFromBytesReturnsOnCall[len(fake.newPublicKeyFromBytesArgsForCall)]
	fake.newPublicKeyFromBytesArgsForCall = append(fake.newPublicKeyFromBytesArgsForCall, struct {
		arg1 []byte
		arg2 []string
	}{arg1Copy, arg2Copy})
	fake.recordInvocation("NewPublicKeyFromBytes", []interface{}{arg1Copy, arg2Copy})
	fake.newPublicKeyFromBytesMutex.Unlock()
	if fake.NewPublicKeyFromBytesStub != nil {
		return fake.NewPublicKeyFromBytesStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.newPublicKeyFromBytesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Issuer) NewPublicKeyFromBytesCallCount() int {
	fake.newPublicKeyFromBytesMutex.RLock()
	defer fake.newPublicKeyFromBytesMutex.RUnlock()
	return len(fake.newPublicKeyFromBytesArgsForCall)
}

func (fake *Issuer) NewPublicKeyFromBytesCalls(stub func([]byte, []string) (handlers.IssuerPublicKey, error)) {
	fake.newPublicKeyFromBytesMutex.Lock()
	defer fake.newPublicKeyFromBytesMutex.Unlock()
	fake.NewPublicKeyFromBytesStub = stub
}

func (fake *Issuer) NewPublicKeyFromBytesArgsForCall(i int) ([]byte, []string) {
	fake.newPublicKeyFromBytesMutex.RLock()
	defer fake.newPublicKeyFromBytesMutex.RUnlock()
	argsForCall := fake.newPublicKeyFromBytesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Issuer) NewPublicKeyFromBytesReturns(result1 handlers.IssuerPublicKey, result2 error) {
	fake.newPublicKeyFromBytesMutex.Lock()
	defer fake.newPublicKeyFromBytesMutex.Unlock()
	fake.NewPublicKeyFromBytesStub = nil
	fake.newPublicKeyFromBytesReturns = struct {
		result1 handlers.IssuerPublicKey
		result2 error
	}{result1, result2}
}

func (fake *Issuer) NewPublicKeyFromBytesReturnsOnCall(i int, result1 handlers.IssuerPublicKey, result2 error) {
	fake.newPublicKeyFromBytesMutex.Lock()
	defer fake.newPublicKeyFromBytesMutex.Unlock()
	fake.NewPublicKeyFromBytesStub = nil
	if fake.newPublicKeyFromBytesReturnsOnCall == nil {
		fake.newPublicKeyFromBytesReturnsOnCall = make(map[int]struct {
			result1 handlers.IssuerPublicKey
			result2 error
		})
	}
	fake.newPublicKeyFromBytesReturnsOnCall[i] = struct {
		result1 handlers.IssuerPublicKey
		result2 error
	}{result1, result2}
}

func (fake *Issuer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newKeyMutex.RLock()
	defer fake.newKeyMutex.RUnlock()
	fake.newPublicKeyFromBytesMutex.RLock()
	defer fake.newPublicKeyFromBytesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Issuer) recordInvocation(key string, args []interface{}) {
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

var _ handlers.Issuer = new(Issuer)
