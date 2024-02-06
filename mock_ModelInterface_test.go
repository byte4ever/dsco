// Code generated by mockery. DO NOT EDIT.

package dsco

import (
	internal "github.com/byte4ever/dsco/internal"
	fvalue "github.com/byte4ever/dsco/internal/fvalue"

	mock "github.com/stretchr/testify/mock"

	plocation "github.com/byte4ever/dsco/internal/plocation"

	reflect "reflect"
)

// MockModelInterface is an autogenerated mock type for the ModelInterface type
type MockModelInterface struct {
	mock.Mock
}

type MockModelInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockModelInterface) EXPECT() *MockModelInterface_Expecter {
	return &MockModelInterface_Expecter{mock: &_m.Mock}
}

// ApplyOn provides a mock function with given fields: g
func (_m *MockModelInterface) ApplyOn(g internal.Getter) (fvalue.Values, error) {
	ret := _m.Called(g)

	if len(ret) == 0 {
		panic("no return value specified for ApplyOn")
	}

	var r0 fvalue.Values
	var r1 error
	if rf, ok := ret.Get(0).(func(internal.Getter) (fvalue.Values, error)); ok {
		return rf(g)
	}
	if rf, ok := ret.Get(0).(func(internal.Getter) fvalue.Values); ok {
		r0 = rf(g)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fvalue.Values)
		}
	}

	if rf, ok := ret.Get(1).(func(internal.Getter) error); ok {
		r1 = rf(g)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockModelInterface_ApplyOn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ApplyOn'
type MockModelInterface_ApplyOn_Call struct {
	*mock.Call
}

// ApplyOn is a helper method to define mock.On call
//   - g internal.Getter
func (_e *MockModelInterface_Expecter) ApplyOn(g interface{}) *MockModelInterface_ApplyOn_Call {
	return &MockModelInterface_ApplyOn_Call{Call: _e.mock.On("ApplyOn", g)}
}

func (_c *MockModelInterface_ApplyOn_Call) Run(run func(g internal.Getter)) *MockModelInterface_ApplyOn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(internal.Getter))
	})
	return _c
}

func (_c *MockModelInterface_ApplyOn_Call) Return(_a0 fvalue.Values, _a1 error) *MockModelInterface_ApplyOn_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockModelInterface_ApplyOn_Call) RunAndReturn(run func(internal.Getter) (fvalue.Values, error)) *MockModelInterface_ApplyOn_Call {
	_c.Call.Return(run)
	return _c
}

// Expand provides a mock function with given fields: g
func (_m *MockModelInterface) Expand(g internal.Expander) error {
	ret := _m.Called(g)

	if len(ret) == 0 {
		panic("no return value specified for Expand")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(internal.Expander) error); ok {
		r0 = rf(g)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockModelInterface_Expand_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Expand'
type MockModelInterface_Expand_Call struct {
	*mock.Call
}

// Expand is a helper method to define mock.On call
//   - g internal.Expander
func (_e *MockModelInterface_Expecter) Expand(g interface{}) *MockModelInterface_Expand_Call {
	return &MockModelInterface_Expand_Call{Call: _e.mock.On("Expand", g)}
}

func (_c *MockModelInterface_Expand_Call) Run(run func(g internal.Expander)) *MockModelInterface_Expand_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(internal.Expander))
	})
	return _c
}

func (_c *MockModelInterface_Expand_Call) Return(_a0 error) *MockModelInterface_Expand_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockModelInterface_Expand_Call) RunAndReturn(run func(internal.Expander) error) *MockModelInterface_Expand_Call {
	_c.Call.Return(run)
	return _c
}

// Fill provides a mock function with given fields: inputModelValue, layers
func (_m *MockModelInterface) Fill(inputModelValue reflect.Value, layers []fvalue.Values) (plocation.Locations, error) {
	ret := _m.Called(inputModelValue, layers)

	if len(ret) == 0 {
		panic("no return value specified for Fill")
	}

	var r0 plocation.Locations
	var r1 error
	if rf, ok := ret.Get(0).(func(reflect.Value, []fvalue.Values) (plocation.Locations, error)); ok {
		return rf(inputModelValue, layers)
	}
	if rf, ok := ret.Get(0).(func(reflect.Value, []fvalue.Values) plocation.Locations); ok {
		r0 = rf(inputModelValue, layers)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(plocation.Locations)
		}
	}

	if rf, ok := ret.Get(1).(func(reflect.Value, []fvalue.Values) error); ok {
		r1 = rf(inputModelValue, layers)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockModelInterface_Fill_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Fill'
type MockModelInterface_Fill_Call struct {
	*mock.Call
}

// Fill is a helper method to define mock.On call
//   - inputModelValue reflect.Value
//   - layers []fvalue.Values
func (_e *MockModelInterface_Expecter) Fill(inputModelValue interface{}, layers interface{}) *MockModelInterface_Fill_Call {
	return &MockModelInterface_Fill_Call{Call: _e.mock.On("Fill", inputModelValue, layers)}
}

func (_c *MockModelInterface_Fill_Call) Run(run func(inputModelValue reflect.Value, layers []fvalue.Values)) *MockModelInterface_Fill_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(reflect.Value), args[1].([]fvalue.Values))
	})
	return _c
}

func (_c *MockModelInterface_Fill_Call) Return(_a0 plocation.Locations, _a1 error) *MockModelInterface_Fill_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockModelInterface_Fill_Call) RunAndReturn(run func(reflect.Value, []fvalue.Values) (plocation.Locations, error)) *MockModelInterface_Fill_Call {
	_c.Call.Return(run)
	return _c
}

// GetFieldValuesFor provides a mock function with given fields: id, v
func (_m *MockModelInterface) GetFieldValuesFor(id string, v reflect.Value) fvalue.Values {
	ret := _m.Called(id, v)

	if len(ret) == 0 {
		panic("no return value specified for GetFieldValuesFor")
	}

	var r0 fvalue.Values
	if rf, ok := ret.Get(0).(func(string, reflect.Value) fvalue.Values); ok {
		r0 = rf(id, v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fvalue.Values)
		}
	}

	return r0
}

// MockModelInterface_GetFieldValuesFor_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFieldValuesFor'
type MockModelInterface_GetFieldValuesFor_Call struct {
	*mock.Call
}

// GetFieldValuesFor is a helper method to define mock.On call
//   - id string
//   - v reflect.Value
func (_e *MockModelInterface_Expecter) GetFieldValuesFor(id interface{}, v interface{}) *MockModelInterface_GetFieldValuesFor_Call {
	return &MockModelInterface_GetFieldValuesFor_Call{Call: _e.mock.On("GetFieldValuesFor", id, v)}
}

func (_c *MockModelInterface_GetFieldValuesFor_Call) Run(run func(id string, v reflect.Value)) *MockModelInterface_GetFieldValuesFor_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(reflect.Value))
	})
	return _c
}

func (_c *MockModelInterface_GetFieldValuesFor_Call) Return(_a0 fvalue.Values) *MockModelInterface_GetFieldValuesFor_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockModelInterface_GetFieldValuesFor_Call) RunAndReturn(run func(string, reflect.Value) fvalue.Values) *MockModelInterface_GetFieldValuesFor_Call {
	_c.Call.Return(run)
	return _c
}

// TypeName provides a mock function with given fields:
func (_m *MockModelInterface) TypeName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TypeName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockModelInterface_TypeName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TypeName'
type MockModelInterface_TypeName_Call struct {
	*mock.Call
}

// TypeName is a helper method to define mock.On call
func (_e *MockModelInterface_Expecter) TypeName() *MockModelInterface_TypeName_Call {
	return &MockModelInterface_TypeName_Call{Call: _e.mock.On("TypeName")}
}

func (_c *MockModelInterface_TypeName_Call) Run(run func()) *MockModelInterface_TypeName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockModelInterface_TypeName_Call) Return(_a0 string) *MockModelInterface_TypeName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockModelInterface_TypeName_Call) RunAndReturn(run func() string) *MockModelInterface_TypeName_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockModelInterface creates a new instance of MockModelInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockModelInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockModelInterface {
	mock := &MockModelInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
