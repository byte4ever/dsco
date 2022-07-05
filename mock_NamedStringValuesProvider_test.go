// Code generated by mockery v2.13.1. DO NOT EDIT.

package dsco

import (
	svalue "github.com/byte4ever/dsco/svalue"
	mock "github.com/stretchr/testify/mock"
)

// MockNamedStringValuesProvider is an autogenerated mock type for the NamedStringValuesProvider type
type MockNamedStringValuesProvider struct {
	mock.Mock
}

// GetName provides a mock function with given fields:
func (_m *MockNamedStringValuesProvider) GetName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetStringValues provides a mock function with given fields:
func (_m *MockNamedStringValuesProvider) GetStringValues() svalue.Values {
	ret := _m.Called()

	var r0 svalue.Values
	if rf, ok := ret.Get(0).(func() svalue.Values); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(svalue.Values)
		}
	}

	return r0
}

type mockConstructorTestingTNewMockNamedStringValuesProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockNamedStringValuesProvider creates a new instance of MockNamedStringValuesProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockNamedStringValuesProvider(t mockConstructorTestingTNewMockNamedStringValuesProvider) *MockNamedStringValuesProvider {
	mock := &MockNamedStringValuesProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}