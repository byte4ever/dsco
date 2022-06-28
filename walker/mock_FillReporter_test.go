// Code generated by mockery v2.13.1. DO NOT EDIT.

package walker

import (
	"github.com/stretchr/testify/mock"

	"github.com/byte4ever/dsco/walker/plocation"
)

// MockFillReporter is an autogenerated mock type for the FillReporter type
type MockFillReporter struct {
	mock.Mock
}

// Failed provides a mock function with given fields:
func (_m *MockFillReporter) Failed() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ReportError provides a mock function with given fields: err
func (_m *MockFillReporter) ReportError(err error) {
	_m.Called(err)
}

// ReportOverride provides a mock function with given fields: uid, Location
func (_m *MockFillReporter) ReportOverride(uid uint, location string) {
	_m.Called(uid, location)
}

// ReportUnused provides a mock function with given fields: path
func (_m *MockFillReporter) ReportUnused(path string) {
	_m.Called(path)
}

// ReportUse provides a mock function with given fields: uid, path, Location
func (_m *MockFillReporter) ReportUse(uid uint, path string, location string) {
	_m.Called(uid, path, location)
}

// Result provides a mock function with given fields:
func (_m *MockFillReporter) Result() (plocation.PathLocations, error) {
	ret := _m.Called()

	var r0 plocation.PathLocations
	if rf, ok := ret.Get(0).(func() plocation.PathLocations); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(plocation.PathLocations)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockFillReporter interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockFillReporter creates a new instance of MockFillReporter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockFillReporter(t mockConstructorTestingTNewMockFillReporter) *MockFillReporter {
	mock := &MockFillReporter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
