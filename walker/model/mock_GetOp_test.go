// Code generated by mockery v2.13.1. DO NOT EDIT.

package model

import (
	fvalues "github.com/byte4ever/dsco/walker/fvalues"
	ifaces "github.com/byte4ever/dsco/walker/ifaces"

	mock "github.com/stretchr/testify/mock"
)

// MockGetOp is an autogenerated mock type for the GetOp type
type MockGetOp struct {
	mock.Mock
}

// Execute provides a mock function with given fields: g
func (_m *MockGetOp) Execute(g ifaces.Getter) (uint, *fvalues.FieldValue, error) {
	ret := _m.Called(g)

	var r0 uint
	if rf, ok := ret.Get(0).(func(ifaces.Getter) uint); ok {
		r0 = rf(g)
	} else {
		r0 = ret.Get(0).(uint)
	}

	var r1 *fvalues.FieldValue
	if rf, ok := ret.Get(1).(func(ifaces.Getter) *fvalues.FieldValue); ok {
		r1 = rf(g)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*fvalues.FieldValue)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(ifaces.Getter) error); ok {
		r2 = rf(g)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewMockGetOp interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockGetOp creates a new instance of MockGetOp. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockGetOp(t mockConstructorTestingTNewMockGetOp) *MockGetOp {
	mock := &MockGetOp{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}