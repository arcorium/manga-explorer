// Code generated by mockery v2.40.1. DO NOT EDIT.

package repository

import (
	common "manga-explorer/internal/common"
	mangas "manga-explorer/internal/domain/mangas"

	mock "github.com/stretchr/testify/mock"
)

// TranslationMock is an autogenerated mock type for the ITranslation type
type TranslationMock struct {
	mock.Mock
}

type TranslationMock_Expecter struct {
	mock *mock.Mock
}

func (_m *TranslationMock) EXPECT() *TranslationMock_Expecter {
	return &TranslationMock_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: translation
func (_m *TranslationMock) Create(translation []mangas.Translation) error {
	ret := _m.Called(translation)

	if len(ret) == 0 {
		panic("no return value specified for Upsert")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]mangas.Translation) error); ok {
		r0 = rf(translation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TranslationMock_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upsert'
type TranslationMock_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - translation []mangas.Translation
func (_e *TranslationMock_Expecter) Create(translation interface{}) *TranslationMock_Create_Call {
	return &TranslationMock_Create_Call{Call: _e.mock.On("Upsert", translation)}
}

func (_c *TranslationMock_Create_Call) Run(run func(translation []mangas.Translation)) *TranslationMock_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]mangas.Translation))
	})
	return _c
}

func (_c *TranslationMock_Create_Call) Return(_a0 error) *TranslationMock_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TranslationMock_Create_Call) RunAndReturn(run func([]mangas.Translation) error) *TranslationMock_Create_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteByIds provides a mock function with given fields: translationIds
func (_m *TranslationMock) DeleteByIds(translationIds []string) error {
	ret := _m.Called(translationIds)

	if len(ret) == 0 {
		panic("no return value specified for DeleteByIds")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(translationIds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TranslationMock_DeleteByIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteByIds'
type TranslationMock_DeleteByIds_Call struct {
	*mock.Call
}

// DeleteByIds is a helper method to define mock.On call
//   - translationIds []string
func (_e *TranslationMock_Expecter) DeleteByIds(translationIds interface{}) *TranslationMock_DeleteByIds_Call {
	return &TranslationMock_DeleteByIds_Call{Call: _e.mock.On("DeleteByIds", translationIds)}
}

func (_c *TranslationMock_DeleteByIds_Call) Run(run func(translationIds []string)) *TranslationMock_DeleteByIds_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string))
	})
	return _c
}

func (_c *TranslationMock_DeleteByIds_Call) Return(_a0 error) *TranslationMock_DeleteByIds_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TranslationMock_DeleteByIds_Call) RunAndReturn(run func([]string) error) *TranslationMock_DeleteByIds_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteByMangaId provides a mock function with given fields: mangaId
func (_m *TranslationMock) DeleteByMangaId(mangaId string) error {
	ret := _m.Called(mangaId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteByMangaId")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(mangaId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TranslationMock_DeleteByMangaId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteByMangaId'
type TranslationMock_DeleteByMangaId_Call struct {
	*mock.Call
}

// DeleteByMangaId is a helper method to define mock.On call
//   - mangaId string
func (_e *TranslationMock_Expecter) DeleteByMangaId(mangaId interface{}) *TranslationMock_DeleteByMangaId_Call {
	return &TranslationMock_DeleteByMangaId_Call{Call: _e.mock.On("DeleteByMangaId", mangaId)}
}

func (_c *TranslationMock_DeleteByMangaId_Call) Run(run func(mangaId string)) *TranslationMock_DeleteByMangaId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *TranslationMock_DeleteByMangaId_Call) Return(_a0 error) *TranslationMock_DeleteByMangaId_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TranslationMock_DeleteByMangaId_Call) RunAndReturn(run func(string) error) *TranslationMock_DeleteByMangaId_Call {
	_c.Call.Return(run)
	return _c
}

// FindById provides a mock function with given fields: id
func (_m *TranslationMock) FindById(id string) (*mangas.Translation, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for FindById")
	}

	var r0 *mangas.Translation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*mangas.Translation, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *mangas.Translation); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mangas.Translation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TranslationMock_FindById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindById'
type TranslationMock_FindById_Call struct {
	*mock.Call
}

// FindById is a helper method to define mock.On call
//   - id string
func (_e *TranslationMock_Expecter) FindById(id interface{}) *TranslationMock_FindById_Call {
	return &TranslationMock_FindById_Call{Call: _e.mock.On("FindById", id)}
}

func (_c *TranslationMock_FindById_Call) Run(run func(id string)) *TranslationMock_FindById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *TranslationMock_FindById_Call) Return(_a0 *mangas.Translation, _a1 error) *TranslationMock_FindById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TranslationMock_FindById_Call) RunAndReturn(run func(string) (*mangas.Translation, error)) *TranslationMock_FindById_Call {
	_c.Call.Return(run)
	return _c
}

// FindByMangaId provides a mock function with given fields: mangaId
func (_m *TranslationMock) FindByMangaId(mangaId string) ([]mangas.Translation, error) {
	ret := _m.Called(mangaId)

	if len(ret) == 0 {
		panic("no return value specified for FindByMangaId")
	}

	var r0 []mangas.Translation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]mangas.Translation, error)); ok {
		return rf(mangaId)
	}
	if rf, ok := ret.Get(0).(func(string) []mangas.Translation); ok {
		r0 = rf(mangaId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]mangas.Translation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(mangaId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TranslationMock_FindByMangaId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByMangaId'
type TranslationMock_FindByMangaId_Call struct {
	*mock.Call
}

// FindByMangaId is a helper method to define mock.On call
//   - mangaId string
func (_e *TranslationMock_Expecter) FindByMangaId(mangaId interface{}) *TranslationMock_FindByMangaId_Call {
	return &TranslationMock_FindByMangaId_Call{Call: _e.mock.On("FindByMangaId", mangaId)}
}

func (_c *TranslationMock_FindByMangaId_Call) Run(run func(mangaId string)) *TranslationMock_FindByMangaId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *TranslationMock_FindByMangaId_Call) Return(_a0 []mangas.Translation, _a1 error) *TranslationMock_FindByMangaId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TranslationMock_FindByMangaId_Call) RunAndReturn(run func(string) ([]mangas.Translation, error)) *TranslationMock_FindByMangaId_Call {
	_c.Call.Return(run)
	return _c
}

// FindMangaSpecific provides a mock function with given fields: mangaId, language
func (_m *TranslationMock) FindMangaSpecific(mangaId string, language common.Language) (*mangas.Translation, error) {
	ret := _m.Called(mangaId, language)

	if len(ret) == 0 {
		panic("no return value specified for FindMangaSpecific")
	}

	var r0 *mangas.Translation
	var r1 error
	if rf, ok := ret.Get(0).(func(string, common.Language) (*mangas.Translation, error)); ok {
		return rf(mangaId, language)
	}
	if rf, ok := ret.Get(0).(func(string, common.Language) *mangas.Translation); ok {
		r0 = rf(mangaId, language)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mangas.Translation)
		}
	}

	if rf, ok := ret.Get(1).(func(string, common.Language) error); ok {
		r1 = rf(mangaId, language)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TranslationMock_FindMangaSpecific_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMangaSpecific'
type TranslationMock_FindMangaSpecific_Call struct {
	*mock.Call
}

// FindMangaSpecific is a helper method to define mock.On call
//   - mangaId string
//   - language common.Language
func (_e *TranslationMock_Expecter) FindMangaSpecific(mangaId interface{}, language interface{}) *TranslationMock_FindMangaSpecific_Call {
	return &TranslationMock_FindMangaSpecific_Call{Call: _e.mock.On("FindMangaSpecific", mangaId, language)}
}

func (_c *TranslationMock_FindMangaSpecific_Call) Run(run func(mangaId string, language common.Language)) *TranslationMock_FindMangaSpecific_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(common.Language))
	})
	return _c
}

func (_c *TranslationMock_FindMangaSpecific_Call) Return(_a0 *mangas.Translation, _a1 error) *TranslationMock_FindMangaSpecific_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TranslationMock_FindMangaSpecific_Call) RunAndReturn(run func(string, common.Language) (*mangas.Translation, error)) *TranslationMock_FindMangaSpecific_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: translation
func (_m *TranslationMock) Update(translation *mangas.Translation) error {
	ret := _m.Called(translation)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*mangas.Translation) error); ok {
		r0 = rf(translation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TranslationMock_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type TranslationMock_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - translation *mangas.Translation
func (_e *TranslationMock_Expecter) Update(translation interface{}) *TranslationMock_Update_Call {
	return &TranslationMock_Update_Call{Call: _e.mock.On("Update", translation)}
}

func (_c *TranslationMock_Update_Call) Run(run func(translation *mangas.Translation)) *TranslationMock_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*mangas.Translation))
	})
	return _c
}

func (_c *TranslationMock_Update_Call) Return(_a0 error) *TranslationMock_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TranslationMock_Update_Call) RunAndReturn(run func(*mangas.Translation) error) *TranslationMock_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewTranslationMock creates a new instance of TranslationMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTranslationMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *TranslationMock {
	mock := &TranslationMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
