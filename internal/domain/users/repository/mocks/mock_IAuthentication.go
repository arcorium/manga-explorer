// Code generated by mockery v2.38.0. DO NOT EDIT.

package repository

import (
	auth "manga-explorer/internal/domain/users"

	mock "github.com/stretchr/testify/mock"
)

// AuthenticationMock is an autogenerated mock type for the IAuthentication type
type AuthenticationMock struct {
	mock.Mock
}

type AuthenticationMock_Expecter struct {
	mock *mock.Mock
}

func (_m *AuthenticationMock) EXPECT() *AuthenticationMock_Expecter {
	return &AuthenticationMock_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: credential
func (_m *AuthenticationMock) Create(credential *auth.Credential) error {
	ret := _m.Called(credential)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*auth.Credential) error); ok {
		r0 = rf(credential)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthenticationMock_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type AuthenticationMock_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - credential *auth.Credential
func (_e *AuthenticationMock_Expecter) Create(credential interface{}) *AuthenticationMock_Create_Call {
	return &AuthenticationMock_Create_Call{Call: _e.mock.On("Create", credential)}
}

func (_c *AuthenticationMock_Create_Call) Run(run func(credential *auth.Credential)) *AuthenticationMock_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*auth.Credential))
	})
	return _c
}

func (_c *AuthenticationMock_Create_Call) Return(_a0 error) *AuthenticationMock_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthenticationMock_Create_Call) RunAndReturn(run func(*auth.Credential) error) *AuthenticationMock_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: userId, credId
func (_m *AuthenticationMock) Find(userId string, credId string) (*auth.Credential, error) {
	ret := _m.Called(userId, credId)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 *auth.Credential
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*auth.Credential, error)); ok {
		return rf(userId, credId)
	}
	if rf, ok := ret.Get(0).(func(string, string) *auth.Credential); ok {
		r0 = rf(userId, credId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Credential)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(userId, credId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthenticationMock_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type AuthenticationMock_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - userId string
//   - credId string
func (_e *AuthenticationMock_Expecter) Find(userId interface{}, credId interface{}) *AuthenticationMock_Find_Call {
	return &AuthenticationMock_Find_Call{Call: _e.mock.On("Find", userId, credId)}
}

func (_c *AuthenticationMock_Find_Call) Run(run func(userId string, credId string)) *AuthenticationMock_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *AuthenticationMock_Find_Call) Return(_a0 *auth.Credential, _a1 error) *AuthenticationMock_Find_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthenticationMock_Find_Call) RunAndReturn(run func(string, string) (*auth.Credential, error)) *AuthenticationMock_Find_Call {
	_c.Call.Return(run)
	return _c
}

// FindByAccessTokenId provides a mock function with given fields: accessTokenId
func (_m *AuthenticationMock) FindByAccessTokenId(accessTokenId string) (*auth.Credential, error) {
	ret := _m.Called(accessTokenId)

	if len(ret) == 0 {
		panic("no return value specified for FindByAccessTokenId")
	}

	var r0 *auth.Credential
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*auth.Credential, error)); ok {
		return rf(accessTokenId)
	}
	if rf, ok := ret.Get(0).(func(string) *auth.Credential); ok {
		r0 = rf(accessTokenId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Credential)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(accessTokenId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthenticationMock_FindByAccessTokenId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByAccessTokenId'
type AuthenticationMock_FindByAccessTokenId_Call struct {
	*mock.Call
}

// FindByAccessTokenId is a helper method to define mock.On call
//   - accessTokenId string
func (_e *AuthenticationMock_Expecter) FindByAccessTokenId(accessTokenId interface{}) *AuthenticationMock_FindByAccessTokenId_Call {
	return &AuthenticationMock_FindByAccessTokenId_Call{Call: _e.mock.On("FindByAccessTokenId", accessTokenId)}
}

func (_c *AuthenticationMock_FindByAccessTokenId_Call) Run(run func(accessTokenId string)) *AuthenticationMock_FindByAccessTokenId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *AuthenticationMock_FindByAccessTokenId_Call) Return(_a0 *auth.Credential, _a1 error) *AuthenticationMock_FindByAccessTokenId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthenticationMock_FindByAccessTokenId_Call) RunAndReturn(run func(string) (*auth.Credential, error)) *AuthenticationMock_FindByAccessTokenId_Call {
	_c.Call.Return(run)
	return _c
}

// FindUserCredentials provides a mock function with given fields: userId
func (_m *AuthenticationMock) FindUserCredentials(userId string) ([]auth.Credential, error) {
	ret := _m.Called(userId)

	if len(ret) == 0 {
		panic("no return value specified for FindUserCredentials")
	}

	var r0 []auth.Credential
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]auth.Credential, error)); ok {
		return rf(userId)
	}
	if rf, ok := ret.Get(0).(func(string) []auth.Credential); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]auth.Credential)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthenticationMock_FindUserCredentials_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindUserCredentials'
type AuthenticationMock_FindUserCredentials_Call struct {
	*mock.Call
}

// FindUserCredentials is a helper method to define mock.On call
//   - userId string
func (_e *AuthenticationMock_Expecter) FindUserCredentials(userId interface{}) *AuthenticationMock_FindUserCredentials_Call {
	return &AuthenticationMock_FindUserCredentials_Call{Call: _e.mock.On("FindUserCredentials", userId)}
}

func (_c *AuthenticationMock_FindUserCredentials_Call) Run(run func(userId string)) *AuthenticationMock_FindUserCredentials_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *AuthenticationMock_FindUserCredentials_Call) Return(_a0 []auth.Credential, _a1 error) *AuthenticationMock_FindUserCredentials_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthenticationMock_FindUserCredentials_Call) RunAndReturn(run func(string) ([]auth.Credential, error)) *AuthenticationMock_FindUserCredentials_Call {
	_c.Call.Return(run)
	return _c
}

// Remove provides a mock function with given fields: userId, credId
func (_m *AuthenticationMock) Remove(userId string, credId string) error {
	ret := _m.Called(userId, credId)

	if len(ret) == 0 {
		panic("no return value specified for Remove")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(userId, credId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthenticationMock_Remove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remove'
type AuthenticationMock_Remove_Call struct {
	*mock.Call
}

// Remove is a helper method to define mock.On call
//   - userId string
//   - credId string
func (_e *AuthenticationMock_Expecter) Remove(userId interface{}, credId interface{}) *AuthenticationMock_Remove_Call {
	return &AuthenticationMock_Remove_Call{Call: _e.mock.On("Remove", userId, credId)}
}

func (_c *AuthenticationMock_Remove_Call) Run(run func(userId string, credId string)) *AuthenticationMock_Remove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *AuthenticationMock_Remove_Call) Return(_a0 error) *AuthenticationMock_Remove_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthenticationMock_Remove_Call) RunAndReturn(run func(string, string) error) *AuthenticationMock_Remove_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveByAccessTokenId provides a mock function with given fields: userId, accessTokenId
func (_m *AuthenticationMock) RemoveByAccessTokenId(userId string, accessTokenId string) error {
	ret := _m.Called(userId, accessTokenId)

	if len(ret) == 0 {
		panic("no return value specified for RemoveByAccessTokenId")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(userId, accessTokenId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthenticationMock_RemoveByAccessTokenId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveByAccessTokenId'
type AuthenticationMock_RemoveByAccessTokenId_Call struct {
	*mock.Call
}

// RemoveByAccessTokenId is a helper method to define mock.On call
//   - userId string
//   - accessTokenId string
func (_e *AuthenticationMock_Expecter) RemoveByAccessTokenId(userId interface{}, accessTokenId interface{}) *AuthenticationMock_RemoveByAccessTokenId_Call {
	return &AuthenticationMock_RemoveByAccessTokenId_Call{Call: _e.mock.On("RemoveByAccessTokenId", userId, accessTokenId)}
}

func (_c *AuthenticationMock_RemoveByAccessTokenId_Call) Run(run func(userId string, accessTokenId string)) *AuthenticationMock_RemoveByAccessTokenId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *AuthenticationMock_RemoveByAccessTokenId_Call) Return(_a0 error) *AuthenticationMock_RemoveByAccessTokenId_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthenticationMock_RemoveByAccessTokenId_Call) RunAndReturn(run func(string, string) error) *AuthenticationMock_RemoveByAccessTokenId_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveUserCredentials provides a mock function with given fields: userId
func (_m *AuthenticationMock) RemoveUserCredentials(userId string) error {
	ret := _m.Called(userId)

	if len(ret) == 0 {
		panic("no return value specified for RemoveUserCredentials")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthenticationMock_RemoveUserCredentials_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveUserCredentials'
type AuthenticationMock_RemoveUserCredentials_Call struct {
	*mock.Call
}

// RemoveUserCredentials is a helper method to define mock.On call
//   - userId string
func (_e *AuthenticationMock_Expecter) RemoveUserCredentials(userId interface{}) *AuthenticationMock_RemoveUserCredentials_Call {
	return &AuthenticationMock_RemoveUserCredentials_Call{Call: _e.mock.On("RemoveUserCredentials", userId)}
}

func (_c *AuthenticationMock_RemoveUserCredentials_Call) Run(run func(userId string)) *AuthenticationMock_RemoveUserCredentials_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *AuthenticationMock_RemoveUserCredentials_Call) Return(_a0 error) *AuthenticationMock_RemoveUserCredentials_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthenticationMock_RemoveUserCredentials_Call) RunAndReturn(run func(string) error) *AuthenticationMock_RemoveUserCredentials_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateAccessTokenId provides a mock function with given fields: credentialId, accessTokenId
func (_m *AuthenticationMock) UpdateAccessTokenId(credentialId string, accessTokenId string) error {
	ret := _m.Called(credentialId, accessTokenId)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAccessTokenId")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(credentialId, accessTokenId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthenticationMock_UpdateAccessTokenId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateAccessTokenId'
type AuthenticationMock_UpdateAccessTokenId_Call struct {
	*mock.Call
}

// UpdateAccessTokenId is a helper method to define mock.On call
//   - credentialId string
//   - accessTokenId string
func (_e *AuthenticationMock_Expecter) UpdateAccessTokenId(credentialId interface{}, accessTokenId interface{}) *AuthenticationMock_UpdateAccessTokenId_Call {
	return &AuthenticationMock_UpdateAccessTokenId_Call{Call: _e.mock.On("UpdateAccessTokenId", credentialId, accessTokenId)}
}

func (_c *AuthenticationMock_UpdateAccessTokenId_Call) Run(run func(credentialId string, accessTokenId string)) *AuthenticationMock_UpdateAccessTokenId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *AuthenticationMock_UpdateAccessTokenId_Call) Return(_a0 error) *AuthenticationMock_UpdateAccessTokenId_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthenticationMock_UpdateAccessTokenId_Call) RunAndReturn(run func(string, string) error) *AuthenticationMock_UpdateAccessTokenId_Call {
	_c.Call.Return(run)
	return _c
}

// NewAuthenticationMock creates a new instance of AuthenticationMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthenticationMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthenticationMock {
	mock := &AuthenticationMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}