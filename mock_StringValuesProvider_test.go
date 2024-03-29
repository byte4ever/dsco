// Code generated by mockery. DO NOT EDIT.

package dsco

import (
	svalue "github.com/byte4ever/dsco/svalue"
	mock "github.com/stretchr/testify/mock"
)

// MockStringValuesProvider is an autogenerated mock type for the StringValuesProvider type
type MockStringValuesProvider struct {
	mock.Mock
}

type MockStringValuesProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockStringValuesProvider) EXPECT() *MockStringValuesProvider_Expecter {
	return &MockStringValuesProvider_Expecter{mock: &_m.Mock}
}

// GetStringValues provides a mock function with given fields:
func (_m *MockStringValuesProvider) GetStringValues() svalue.Values {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetStringValues")
	}

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

// MockStringValuesProvider_GetStringValues_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetStringValues'
type MockStringValuesProvider_GetStringValues_Call struct {
	*mock.Call
}

// GetStringValues is a helper method to define mock.On call
func (_e *MockStringValuesProvider_Expecter) GetStringValues() *MockStringValuesProvider_GetStringValues_Call {
	return &MockStringValuesProvider_GetStringValues_Call{Call: _e.mock.On("GetStringValues")}
}

func (_c *MockStringValuesProvider_GetStringValues_Call) Run(run func()) *MockStringValuesProvider_GetStringValues_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockStringValuesProvider_GetStringValues_Call) Return(_a0 svalue.Values) *MockStringValuesProvider_GetStringValues_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockStringValuesProvider_GetStringValues_Call) RunAndReturn(run func() svalue.Values) *MockStringValuesProvider_GetStringValues_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockStringValuesProvider creates a new instance of MockStringValuesProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockStringValuesProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockStringValuesProvider {
	mock := &MockStringValuesProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
