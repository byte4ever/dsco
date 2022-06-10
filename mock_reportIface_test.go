// Code generated by mockery v2.12.3. DO NOT EDIT.

package dsco

import mock "github.com/stretchr/testify/mock"

// mockReportIface is an autogenerated mock type for the reportInterface type
type mockReportIface struct {
	mock.Mock
}

// addEntry provides a mock function with given fields: e
func (_m *mockReportIface) addEntry(e ReportEntry) {
	_m.Called(e)
}

// perEntryReport provides a mock function with given fields:
func (_m *mockReportIface) perEntryReport() []error {
	ret := _m.Called()

	var r0 []error
	if rf, ok := ret.Get(0).(func() []error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]error)
		}
	}

	return r0
}

type newMockReportIfaceT interface {
	mock.TestingT
	Cleanup(func())
}

// newMockReportIface creates a new instance of mockReportIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockReportIface(t newMockReportIfaceT) *mockReportIface {
	mock := &mockReportIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
