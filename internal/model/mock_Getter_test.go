// Code generated by mockery v2.13.1. DO NOT EDIT.

package model

import (
	"github.com/stretchr/testify/mock"

	"reflect"

	"github.com/byte4ever/dsco/internal/fvalue"
)

// MockGetter is an autogenerated mock type for the Getter type
type MockGetter struct {
	mock.Mock
}

// Get provides a mock function with given fields: path, _type
func (_m *MockGetter) Get(path string, _type reflect.Type) (*fvalue.Value, error) {
	ret := _m.Called(path, _type)

	var r0 *fvalue.Value
	if rf, ok := ret.Get(0).(func(string, reflect.Type) *fvalue.Value); ok {
		r0 = rf(path, _type)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fvalue.Value)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, reflect.Type) error); ok {
		r1 = rf(path, _type)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockGetter interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockGetter creates a new instance of MockGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockGetter(t mockConstructorTestingTNewMockGetter) *MockGetter {
	mock := &MockGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
