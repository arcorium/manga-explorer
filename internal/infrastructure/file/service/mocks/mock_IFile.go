// Code generated by mockery v2.40.1. DO NOT EDIT.

package service

import (
	file "manga-explorer/internal/infrastructure/file"

	mock "github.com/stretchr/testify/mock"

	multipart "mime/multipart"

	status "manga-explorer/internal/common/status"
)

// FileMock is an autogenerated mock type for the IFile type
type FileMock struct {
	mock.Mock
}

type FileMock_Expecter struct {
	mock *mock.Mock
}

func (_m *FileMock) EXPECT() *FileMock_Expecter {
	return &FileMock_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: types, filename
func (_m *FileMock) Delete(types file.AssetType, filename file.Name) status.Object {
	ret := _m.Called(types, filename)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(file.AssetType, file.Name) status.Object); ok {
		r0 = rf(types, filename)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// FileMock_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type FileMock_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - types file.AssetType
//   - filename file.Name
func (_e *FileMock_Expecter) Delete(types interface{}, filename interface{}) *FileMock_Delete_Call {
	return &FileMock_Delete_Call{Call: _e.mock.On("Delete", types, filename)}
}

func (_c *FileMock_Delete_Call) Run(run func(types file.AssetType, filename file.Name)) *FileMock_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(file.AssetType), args[1].(file.Name))
	})
	return _c
}

func (_c *FileMock_Delete_Call) Return(_a0 status.Object) *FileMock_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileMock_Delete_Call) RunAndReturn(run func(file.AssetType, file.Name) status.Object) *FileMock_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Endpoint provides a mock function with given fields: assetType
func (_m *FileMock) Endpoint(assetType file.AssetType) string {
	ret := _m.Called(assetType)

	if len(ret) == 0 {
		panic("no return value specified for Endpoint")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(file.AssetType) string); ok {
		r0 = rf(assetType)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// FileMock_Endpoint_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Endpoint'
type FileMock_Endpoint_Call struct {
	*mock.Call
}

// Endpoint is a helper method to define mock.On call
//   - assetType file.AssetType
func (_e *FileMock_Expecter) Endpoint(assetType interface{}) *FileMock_Endpoint_Call {
	return &FileMock_Endpoint_Call{Call: _e.mock.On("Endpoint", assetType)}
}

func (_c *FileMock_Endpoint_Call) Run(run func(assetType file.AssetType)) *FileMock_Endpoint_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(file.AssetType))
	})
	return _c
}

func (_c *FileMock_Endpoint_Call) Return(_a0 string) *FileMock_Endpoint_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileMock_Endpoint_Call) RunAndReturn(run func(file.AssetType) string) *FileMock_Endpoint_Call {
	_c.Call.Return(run)
	return _c
}

// GetFullpath provides a mock function with given fields: assetType, filename
func (_m *FileMock) GetFullpath(assetType file.AssetType, filename file.Name) string {
	ret := _m.Called(assetType, filename)

	if len(ret) == 0 {
		panic("no return value specified for GetFullpath")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(file.AssetType, file.Name) string); ok {
		r0 = rf(assetType, filename)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// FileMock_GetFullpath_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFullpath'
type FileMock_GetFullpath_Call struct {
	*mock.Call
}

// GetFullpath is a helper method to define mock.On call
//   - assetType file.AssetType
//   - filename file.Name
func (_e *FileMock_Expecter) GetFullpath(assetType interface{}, filename interface{}) *FileMock_GetFullpath_Call {
	return &FileMock_GetFullpath_Call{Call: _e.mock.On("GetFullpath", assetType, filename)}
}

func (_c *FileMock_GetFullpath_Call) Run(run func(assetType file.AssetType, filename file.Name)) *FileMock_GetFullpath_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(file.AssetType), args[1].(file.Name))
	})
	return _c
}

func (_c *FileMock_GetFullpath_Call) Return(_a0 string) *FileMock_GetFullpath_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileMock_GetFullpath_Call) RunAndReturn(run func(file.AssetType, file.Name) string) *FileMock_GetFullpath_Call {
	_c.Call.Return(run)
	return _c
}

// Upload provides a mock function with given fields: types, header
func (_m *FileMock) Upload(types file.AssetType, header *multipart.FileHeader) (file.Name, status.Object) {
	ret := _m.Called(types, header)

	if len(ret) == 0 {
		panic("no return value specified for Upload")
	}

	var r0 file.Name
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(file.AssetType, *multipart.FileHeader) (file.Name, status.Object)); ok {
		return rf(types, header)
	}
	if rf, ok := ret.Get(0).(func(file.AssetType, *multipart.FileHeader) file.Name); ok {
		r0 = rf(types, header)
	} else {
		r0 = ret.Get(0).(file.Name)
	}

	if rf, ok := ret.Get(1).(func(file.AssetType, *multipart.FileHeader) status.Object); ok {
		r1 = rf(types, header)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// FileMock_Upload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upload'
type FileMock_Upload_Call struct {
	*mock.Call
}

// Upload is a helper method to define mock.On call
//   - types file.AssetType
//   - header *multipart.FileHeader
func (_e *FileMock_Expecter) Upload(types interface{}, header interface{}) *FileMock_Upload_Call {
	return &FileMock_Upload_Call{Call: _e.mock.On("Upload", types, header)}
}

func (_c *FileMock_Upload_Call) Run(run func(types file.AssetType, header *multipart.FileHeader)) *FileMock_Upload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(file.AssetType), args[1].(*multipart.FileHeader))
	})
	return _c
}

func (_c *FileMock_Upload_Call) Return(_a0 file.Name, _a1 status.Object) *FileMock_Upload_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FileMock_Upload_Call) RunAndReturn(run func(file.AssetType, *multipart.FileHeader) (file.Name, status.Object)) *FileMock_Upload_Call {
	_c.Call.Return(run)
	return _c
}

// Uploads provides a mock function with given fields: types, header
func (_m *FileMock) Uploads(types file.AssetType, header []multipart.FileHeader) ([]file.Name, status.Object) {
	ret := _m.Called(types, header)

	if len(ret) == 0 {
		panic("no return value specified for Uploads")
	}

	var r0 []file.Name
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(file.AssetType, []multipart.FileHeader) ([]file.Name, status.Object)); ok {
		return rf(types, header)
	}
	if rf, ok := ret.Get(0).(func(file.AssetType, []multipart.FileHeader) []file.Name); ok {
		r0 = rf(types, header)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]file.Name)
		}
	}

	if rf, ok := ret.Get(1).(func(file.AssetType, []multipart.FileHeader) status.Object); ok {
		r1 = rf(types, header)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// FileMock_Uploads_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Uploads'
type FileMock_Uploads_Call struct {
	*mock.Call
}

// Uploads is a helper method to define mock.On call
//   - types file.AssetType
//   - header []multipart.FileHeader
func (_e *FileMock_Expecter) Uploads(types interface{}, header interface{}) *FileMock_Uploads_Call {
	return &FileMock_Uploads_Call{Call: _e.mock.On("Uploads", types, header)}
}

func (_c *FileMock_Uploads_Call) Run(run func(types file.AssetType, header []multipart.FileHeader)) *FileMock_Uploads_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(file.AssetType), args[1].([]multipart.FileHeader))
	})
	return _c
}

func (_c *FileMock_Uploads_Call) Return(_a0 []file.Name, _a1 status.Object) *FileMock_Uploads_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FileMock_Uploads_Call) RunAndReturn(run func(file.AssetType, []multipart.FileHeader) ([]file.Name, status.Object)) *FileMock_Uploads_Call {
	_c.Call.Return(run)
	return _c
}

// NewFileMock creates a new instance of FileMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFileMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *FileMock {
	mock := &FileMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}