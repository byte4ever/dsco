// Code generated by mockery v2.13.1. DO NOT EDIT.

package model

import (
	fvalue "github.com/byte4ever/dsco/internal/fvalue"
	mock "github.com/stretchr/testify/mock"

	plocation "github.com/byte4ever/dsco/internal/plocation"

	reflect "reflect"
)

// MockNode is an autogenerated mock type for the Node type
type MockNode struct {
	mock.Mock
}

// BuildGetList provides a mock function with given fields: s
func (_m *MockNode) BuildGetList(s *GetList) {
	_m.Called(s)
}

// FeedFieldValues provides a mock function with given fields: srcID, fieldValues, value
func (_m *MockNode) FeedFieldValues(srcID string, fieldValues fvalue.Values, value reflect.Value) {
	_m.Called(srcID, fieldValues, value)
}

// Fill provides a mock function with given fields: value, layers
func (_m *MockNode) Fill(value reflect.Value, layers []fvalue.Values) (plocation.Locations, error) {
	ret := _m.Called(value, layers)

	var r0 plocation.Locations
	if rf, ok := ret.Get(0).(func(reflect.Value, []fvalue.Values) plocation.Locations); ok {
		r0 = rf(value, layers)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(plocation.Locations)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(reflect.Value, []fvalue.Values) error); ok {
		r1 = rf(value, layers)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockNode interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockNode creates a new instance of MockNode. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockNode(t mockConstructorTestingTNewMockNode) *MockNode {
	mock := &MockNode{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
