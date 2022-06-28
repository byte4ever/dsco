// Code generated by mockery v2.13.1. DO NOT EDIT.

package walker

import (
	"reflect"

	"github.com/stretchr/testify/mock"

	"github.com/byte4ever/dsco/walker/plocation"
)

// MockModelInterface is an autogenerated mock type for the ModelInterface type
type MockModelInterface struct {
	mock.Mock
}

// ApplyOn provides a mock function with given fields: g
func (_m *MockModelInterface) ApplyOn(g Getter) (FieldValues, error) {
	ret := _m.Called(g)

	var r0 FieldValues
	if rf, ok := ret.Get(0).(func(Getter) FieldValues); ok {
		r0 = rf(g)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(FieldValues)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(Getter) error); ok {
		r1 = rf(g)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FeedFieldValues provides a mock function with given fields: id, v
func (_m *MockModelInterface) FeedFieldValues(id string, v reflect.Value) FieldValues {
	ret := _m.Called(id, v)

	var r0 FieldValues
	if rf, ok := ret.Get(0).(func(string, reflect.Value) FieldValues); ok {
		r0 = rf(id, v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(FieldValues)
		}
	}

	return r0
}

// Fill provides a mock function with given fields: inputModelValue, layers
func (_m *MockModelInterface) Fill(inputModelValue reflect.Value, layers []FieldValues) (plocation.PathLocations, error) {
	ret := _m.Called(inputModelValue, layers)

	var r0 plocation.PathLocations
	if rf, ok := ret.Get(0).(func(reflect.Value, []FieldValues) plocation.PathLocations); ok {
		r0 = rf(inputModelValue, layers)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(plocation.PathLocations)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(reflect.Value, []FieldValues) error); ok {
		r1 = rf(inputModelValue, layers)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TypeName provides a mock function with given fields:
func (_m *MockModelInterface) TypeName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewMockModelInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockModelInterface creates a new instance of MockModelInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockModelInterface(t mockConstructorTestingTNewMockModelInterface) *MockModelInterface {
	mock := &MockModelInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}