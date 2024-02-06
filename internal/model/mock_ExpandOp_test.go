// Code generated by mockery. DO NOT EDIT.

package model

import (
	internal "github.com/byte4ever/dsco/internal"
	mock "github.com/stretchr/testify/mock"
)

// MockExpandOp is an autogenerated mock type for the ExpandOp type
type MockExpandOp struct {
	mock.Mock
}

type MockExpandOp_Expecter struct {
	mock *mock.Mock
}

func (_m *MockExpandOp) EXPECT() *MockExpandOp_Expecter {
	return &MockExpandOp_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: g
func (_m *MockExpandOp) Execute(g internal.Expander) error {
	ret := _m.Called(g)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(internal.Expander) error); ok {
		r0 = rf(g)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockExpandOp_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockExpandOp_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - g internal.Expander
func (_e *MockExpandOp_Expecter) Execute(g interface{}) *MockExpandOp_Execute_Call {
	return &MockExpandOp_Execute_Call{Call: _e.mock.On("Execute", g)}
}

func (_c *MockExpandOp_Execute_Call) Run(run func(g internal.Expander)) *MockExpandOp_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(internal.Expander))
	})
	return _c
}

func (_c *MockExpandOp_Execute_Call) Return(err error) *MockExpandOp_Execute_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockExpandOp_Execute_Call) RunAndReturn(run func(internal.Expander) error) *MockExpandOp_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockExpandOp creates a new instance of MockExpandOp. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExpandOp(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExpandOp {
	mock := &MockExpandOp{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
