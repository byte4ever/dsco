// Code generated by mockery. DO NOT EDIT.

package dsco

import (
	fvalue "github.com/byte4ever/dsco/internal/fvalue"
	mock "github.com/stretchr/testify/mock"
)

// MockFieldValuesGetter is an autogenerated mock type for the FieldValuesGetter type
type MockFieldValuesGetter struct {
	mock.Mock
}

type MockFieldValuesGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFieldValuesGetter) EXPECT() *MockFieldValuesGetter_Expecter {
	return &MockFieldValuesGetter_Expecter{mock: &_m.Mock}
}

// GetFieldValuesFrom provides a mock function with given fields: model
func (_m *MockFieldValuesGetter) GetFieldValuesFrom(model ModelInterface) (fvalue.Values, error) {
	ret := _m.Called(model)

	if len(ret) == 0 {
		panic("no return value specified for GetFieldValuesFrom")
	}

	var r0 fvalue.Values
	var r1 error
	if rf, ok := ret.Get(0).(func(ModelInterface) (fvalue.Values, error)); ok {
		return rf(model)
	}
	if rf, ok := ret.Get(0).(func(ModelInterface) fvalue.Values); ok {
		r0 = rf(model)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fvalue.Values)
		}
	}

	if rf, ok := ret.Get(1).(func(ModelInterface) error); ok {
		r1 = rf(model)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFieldValuesGetter_GetFieldValuesFrom_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFieldValuesFrom'
type MockFieldValuesGetter_GetFieldValuesFrom_Call struct {
	*mock.Call
}

// GetFieldValuesFrom is a helper method to define mock.On call
//   - model ModelInterface
func (_e *MockFieldValuesGetter_Expecter) GetFieldValuesFrom(model interface{}) *MockFieldValuesGetter_GetFieldValuesFrom_Call {
	return &MockFieldValuesGetter_GetFieldValuesFrom_Call{Call: _e.mock.On("GetFieldValuesFrom", model)}
}

func (_c *MockFieldValuesGetter_GetFieldValuesFrom_Call) Run(run func(model ModelInterface)) *MockFieldValuesGetter_GetFieldValuesFrom_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ModelInterface))
	})
	return _c
}

func (_c *MockFieldValuesGetter_GetFieldValuesFrom_Call) Return(_a0 fvalue.Values, _a1 error) *MockFieldValuesGetter_GetFieldValuesFrom_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFieldValuesGetter_GetFieldValuesFrom_Call) RunAndReturn(run func(ModelInterface) (fvalue.Values, error)) *MockFieldValuesGetter_GetFieldValuesFrom_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockFieldValuesGetter creates a new instance of MockFieldValuesGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFieldValuesGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFieldValuesGetter {
	mock := &MockFieldValuesGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
