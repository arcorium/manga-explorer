// Code generated by mockery v2.40.1. DO NOT EDIT.

package service

import (
	common "manga-explorer/internal/app/common"
	dto2 "manga-explorer/internal/common/dto"

	dto "manga-explorer/internal/domain/mangas/dto"

	mock "github.com/stretchr/testify/mock"

	status "manga-explorer/internal/app/common/status"
)

// MangaMock is an autogenerated mock type for the IManga type
type MangaMock struct {
	mock.Mock
}

type MangaMock_Expecter struct {
	mock *mock.Mock
}

func (_m *MangaMock) EXPECT() *MangaMock_Expecter {
	return &MangaMock_Expecter{mock: &_m.Mock}
}

// CreateComments provides a mock function with given fields: input
func (_m *MangaMock) CreateComments(input *dto.MangaCommentCreateInput) status.Object {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for CreateComments")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(*dto.MangaCommentCreateInput) status.Object); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// MangaMock_CreateComments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateComments'
type MangaMock_CreateComments_Call struct {
	*mock.Call
}

// CreateComments is a helper method to define mock.On call
//   - input *dto.MangaCommentCreateInput
func (_e *MangaMock_Expecter) CreateComments(input interface{}) *MangaMock_CreateComments_Call {
	return &MangaMock_CreateComments_Call{Call: _e.mock.On("CreateComments", input)}
}

func (_c *MangaMock_CreateComments_Call) Run(run func(input *dto.MangaCommentCreateInput)) *MangaMock_CreateComments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.MangaCommentCreateInput))
	})
	return _c
}

func (_c *MangaMock_CreateComments_Call) Return(_a0 status.Object) *MangaMock_CreateComments_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MangaMock_CreateComments_Call) RunAndReturn(run func(*dto.MangaCommentCreateInput) status.Object) *MangaMock_CreateComments_Call {
	_c.Call.Return(run)
	return _c
}

// CreateManga provides a mock function with given fields: input
func (_m *MangaMock) CreateManga(input *dto.MangaCreateInput) status.Object {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for CreateManga")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(*dto.MangaCreateInput) status.Object); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// MangaMock_CreateManga_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateManga'
type MangaMock_CreateManga_Call struct {
	*mock.Call
}

// CreateManga is a helper method to define mock.On call
//   - input *dto.MangaCreateInput
func (_e *MangaMock_Expecter) CreateManga(input interface{}) *MangaMock_CreateManga_Call {
	return &MangaMock_CreateManga_Call{Call: _e.mock.On("CreateManga", input)}
}

func (_c *MangaMock_CreateManga_Call) Run(run func(input *dto.MangaCreateInput)) *MangaMock_CreateManga_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.MangaCreateInput))
	})
	return _c
}

func (_c *MangaMock_CreateManga_Call) Return(_a0 status.Object) *MangaMock_CreateManga_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MangaMock_CreateManga_Call) RunAndReturn(run func(*dto.MangaCreateInput) status.Object) *MangaMock_CreateManga_Call {
	_c.Call.Return(run)
	return _c
}

// CreateVolume provides a mock function with given fields: input
func (_m *MangaMock) CreateVolume(input *dto.VolumeCreateInput) status.Object {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for CreateVolume")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(*dto.VolumeCreateInput) status.Object); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// MangaMock_CreateVolume_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateVolume'
type MangaMock_CreateVolume_Call struct {
	*mock.Call
}

// CreateVolume is a helper method to define mock.On call
//   - input *dto.VolumeCreateInput
func (_e *MangaMock_Expecter) CreateVolume(input interface{}) *MangaMock_CreateVolume_Call {
	return &MangaMock_CreateVolume_Call{Call: _e.mock.On("CreateVolume", input)}
}

func (_c *MangaMock_CreateVolume_Call) Run(run func(input *dto.VolumeCreateInput)) *MangaMock_CreateVolume_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.VolumeCreateInput))
	})
	return _c
}

func (_c *MangaMock_CreateVolume_Call) Return(_a0 status.Object) *MangaMock_CreateVolume_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MangaMock_CreateVolume_Call) RunAndReturn(run func(*dto.VolumeCreateInput) status.Object) *MangaMock_CreateVolume_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteMangaTranslations provides a mock function with given fields: mangaId
func (_m *MangaMock) DeleteMangaTranslations(mangaId string) status.Object {
	ret := _m.Called(mangaId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMangaTranslations")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(string) status.Object); ok {
		r0 = rf(mangaId)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// MangaMock_DeleteMangaTranslations_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteMangaTranslations'
type MangaMock_DeleteMangaTranslations_Call struct {
	*mock.Call
}

// DeleteMangaTranslations is a helper method to define mock.On call
//   - mangaId string
func (_e *MangaMock_Expecter) DeleteMangaTranslations(mangaId interface{}) *MangaMock_DeleteMangaTranslations_Call {
	return &MangaMock_DeleteMangaTranslations_Call{Call: _e.mock.On("DeleteMangaTranslations", mangaId)}
}

func (_c *MangaMock_DeleteMangaTranslations_Call) Run(run func(mangaId string)) *MangaMock_DeleteMangaTranslations_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MangaMock_DeleteMangaTranslations_Call) Return(_a0 status.Object) *MangaMock_DeleteMangaTranslations_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MangaMock_DeleteMangaTranslations_Call) RunAndReturn(run func(string) status.Object) *MangaMock_DeleteMangaTranslations_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteTranslations provides a mock function with given fields: input
func (_m *MangaMock) DeleteTranslations(input *dto.TranslationDeleteInput) status.Object {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for DeleteTranslations")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(*dto.TranslationDeleteInput) status.Object); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// MangaMock_DeleteTranslations_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteTranslations'
type MangaMock_DeleteTranslations_Call struct {
	*mock.Call
}

// DeleteTranslations is a helper method to define mock.On call
//   - input *dto.TranslationDeleteInput
func (_e *MangaMock_Expecter) DeleteTranslations(input interface{}) *MangaMock_DeleteTranslations_Call {
	return &MangaMock_DeleteTranslations_Call{Call: _e.mock.On("DeleteTranslations", input)}
}

func (_c *MangaMock_DeleteTranslations_Call) Run(run func(input *dto.TranslationDeleteInput)) *MangaMock_DeleteTranslations_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.TranslationDeleteInput))
	})
	return _c
}

func (_c *MangaMock_DeleteTranslations_Call) Return(_a0 status.Object) *MangaMock_DeleteTranslations_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MangaMock_DeleteTranslations_Call) RunAndReturn(run func(*dto.TranslationDeleteInput) status.Object) *MangaMock_DeleteTranslations_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteVolume provides a mock function with given fields: input
func (_m *MangaMock) DeleteVolume(input *dto.VolumeDeleteInput) status.Object {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for DeleteVolume")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(*dto.VolumeDeleteInput) status.Object); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// MangaMock_DeleteVolume_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteVolume'
type MangaMock_DeleteVolume_Call struct {
	*mock.Call
}

// DeleteVolume is a helper method to define mock.On call
//   - input *dto.VolumeDeleteInput
func (_e *MangaMock_Expecter) DeleteVolume(input interface{}) *MangaMock_DeleteVolume_Call {
	return &MangaMock_DeleteVolume_Call{Call: _e.mock.On("DeleteVolume", input)}
}

func (_c *MangaMock_DeleteVolume_Call) Run(run func(input *dto.VolumeDeleteInput)) *MangaMock_DeleteVolume_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.VolumeDeleteInput))
	})
	return _c
}

func (_c *MangaMock_DeleteVolume_Call) Return(_a0 status.Object) *MangaMock_DeleteVolume_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MangaMock_DeleteVolume_Call) RunAndReturn(run func(*dto.VolumeDeleteInput) status.Object) *MangaMock_DeleteVolume_Call {
	_c.Call.Return(run)
	return _c
}

// EditManga provides a mock function with given fields: input
func (_m *MangaMock) EditManga(input *dto.MangaEditInput) status.Object {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for EditManga")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(*dto.MangaEditInput) status.Object); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// MangaMock_EditManga_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EditManga'
type MangaMock_EditManga_Call struct {
	*mock.Call
}

// EditManga is a helper method to define mock.On call
//   - input *dto.MangaEditInput
func (_e *MangaMock_Expecter) EditManga(input interface{}) *MangaMock_EditManga_Call {
	return &MangaMock_EditManga_Call{Call: _e.mock.On("EditManga", input)}
}

func (_c *MangaMock_EditManga_Call) Run(run func(input *dto.MangaEditInput)) *MangaMock_EditManga_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.MangaEditInput))
	})
	return _c
}

func (_c *MangaMock_EditManga_Call) Return(_a0 status.Object) *MangaMock_EditManga_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MangaMock_EditManga_Call) RunAndReturn(run func(*dto.MangaEditInput) status.Object) *MangaMock_EditManga_Call {
	_c.Call.Return(run)
	return _c
}

// FindMangaByIds provides a mock function with given fields: mangaId
func (_m *MangaMock) FindMangaByIds(mangaId ...string) ([]dto.MangaResponse, status.Object) {
	_va := make([]interface{}, len(mangaId))
	for _i := range mangaId {
		_va[_i] = mangaId[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindMangaByIds")
	}

	var r0 []dto.MangaResponse
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(...string) ([]dto.MangaResponse, status.Object)); ok {
		return rf(mangaId...)
	}
	if rf, ok := ret.Get(0).(func(...string) []dto.MangaResponse); ok {
		r0 = rf(mangaId...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.MangaResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(...string) status.Object); ok {
		r1 = rf(mangaId...)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// MangaMock_FindMangaByIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMangaByIds'
type MangaMock_FindMangaByIds_Call struct {
	*mock.Call
}

// FindMangaByIds is a helper method to define mock.On call
//   - mangaId ...string
func (_e *MangaMock_Expecter) FindMangaByIds(mangaId ...interface{}) *MangaMock_FindMangaByIds_Call {
	return &MangaMock_FindMangaByIds_Call{Call: _e.mock.On("FindMangaByIds",
		append([]interface{}{}, mangaId...)...)}
}

func (_c *MangaMock_FindMangaByIds_Call) Run(run func(mangaId ...string)) *MangaMock_FindMangaByIds_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *MangaMock_FindMangaByIds_Call) Return(_a0 []dto.MangaResponse, _a1 status.Object) *MangaMock_FindMangaByIds_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MangaMock_FindMangaByIds_Call) RunAndReturn(run func(...string) ([]dto.MangaResponse, status.Object)) *MangaMock_FindMangaByIds_Call {
	_c.Call.Return(run)
	return _c
}

// FindMangaComments provides a mock function with given fields: mangaId
func (_m *MangaMock) FindMangaComments(mangaId string) ([]dto.CommentResponse, status.Object) {
	ret := _m.Called(mangaId)

	if len(ret) == 0 {
		panic("no return value specified for FindMangaComments")
	}

	var r0 []dto.CommentResponse
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(string) ([]dto.CommentResponse, status.Object)); ok {
		return rf(mangaId)
	}
	if rf, ok := ret.Get(0).(func(string) []dto.CommentResponse); ok {
		r0 = rf(mangaId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.CommentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) status.Object); ok {
		r1 = rf(mangaId)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// MangaMock_FindMangaComments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMangaComments'
type MangaMock_FindMangaComments_Call struct {
	*mock.Call
}

// FindMangaComments is a helper method to define mock.On call
//   - mangaId string
func (_e *MangaMock_Expecter) FindMangaComments(mangaId interface{}) *MangaMock_FindMangaComments_Call {
	return &MangaMock_FindMangaComments_Call{Call: _e.mock.On("FindMangaComments", mangaId)}
}

func (_c *MangaMock_FindMangaComments_Call) Run(run func(mangaId string)) *MangaMock_FindMangaComments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MangaMock_FindMangaComments_Call) Return(_a0 []dto.CommentResponse, _a1 status.Object) *MangaMock_FindMangaComments_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MangaMock_FindMangaComments_Call) RunAndReturn(run func(string) ([]dto.CommentResponse, status.Object)) *MangaMock_FindMangaComments_Call {
	_c.Call.Return(run)
	return _c
}

// FindMangaFavorites provides a mock function with given fields: userId, query
func (_m *MangaMock) FindMangaFavorites(userId string, query *dto2.PagedQueryInput) ([]dto.MangaFavoriteResponse, *dto2.ResponsePage, status.Object) {
	ret := _m.Called(userId, query)

	if len(ret) == 0 {
		panic("no return value specified for FindMangaFavorites")
	}

	var r0 []dto.MangaFavoriteResponse
	var r1 *dto2.ResponsePage
	var r2 status.Object
	if rf, ok := ret.Get(0).(func(string, *dto2.PagedQueryInput) ([]dto.MangaFavoriteResponse, *dto2.ResponsePage, status.Object)); ok {
		return rf(userId, query)
	}
	if rf, ok := ret.Get(0).(func(string, *dto2.PagedQueryInput) []dto.MangaFavoriteResponse); ok {
		r0 = rf(userId, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.MangaFavoriteResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *dto2.PagedQueryInput) *dto2.ResponsePage); ok {
		r1 = rf(userId, query)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dto2.ResponsePage)
		}
	}

	if rf, ok := ret.Get(2).(func(string, *dto2.PagedQueryInput) status.Object); ok {
		r2 = rf(userId, query)
	} else {
		r2 = ret.Get(2).(status.Object)
	}

	return r0, r1, r2
}

// MangaMock_FindMangaFavorites_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMangaFavorites'
type MangaMock_FindMangaFavorites_Call struct {
	*mock.Call
}

// FindMangaFavorites is a helper method to define mock.On call
//   - userId string
//   - query *appdto.PagedQueryInput
func (_e *MangaMock_Expecter) FindMangaFavorites(userId interface{}, query interface{}) *MangaMock_FindMangaFavorites_Call {
	return &MangaMock_FindMangaFavorites_Call{Call: _e.mock.On("FindMangaFavorites", userId, query)}
}

func (_c *MangaMock_FindMangaFavorites_Call) Run(run func(userId string, query *dto2.PagedQueryInput)) *MangaMock_FindMangaFavorites_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*dto2.PagedQueryInput))
	})
	return _c
}

func (_c *MangaMock_FindMangaFavorites_Call) Return(_a0 []dto.MangaFavoriteResponse, _a1 *dto2.ResponsePage, _a2 status.Object) *MangaMock_FindMangaFavorites_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MangaMock_FindMangaFavorites_Call) RunAndReturn(run func(string, *dto2.PagedQueryInput) ([]dto.MangaFavoriteResponse, *dto2.ResponsePage, status.Object)) *MangaMock_FindMangaFavorites_Call {
	_c.Call.Return(run)
	return _c
}

// FindMangaHistories provides a mock function with given fields: userId, query
func (_m *MangaMock) FindMangaHistories(userId string, query *dto2.PagedQueryInput) ([]dto.MangaHistoryResponse, *dto2.ResponsePage, status.Object) {
	ret := _m.Called(userId, query)

	if len(ret) == 0 {
		panic("no return value specified for FindMangaHistories")
	}

	var r0 []dto.MangaHistoryResponse
	var r1 *dto2.ResponsePage
	var r2 status.Object
	if rf, ok := ret.Get(0).(func(string, *dto2.PagedQueryInput) ([]dto.MangaHistoryResponse, *dto2.ResponsePage, status.Object)); ok {
		return rf(userId, query)
	}
	if rf, ok := ret.Get(0).(func(string, *dto2.PagedQueryInput) []dto.MangaHistoryResponse); ok {
		r0 = rf(userId, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.MangaHistoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *dto2.PagedQueryInput) *dto2.ResponsePage); ok {
		r1 = rf(userId, query)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dto2.ResponsePage)
		}
	}

	if rf, ok := ret.Get(2).(func(string, *dto2.PagedQueryInput) status.Object); ok {
		r2 = rf(userId, query)
	} else {
		r2 = ret.Get(2).(status.Object)
	}

	return r0, r1, r2
}

// MangaMock_FindMangaHistories_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMangaHistories'
type MangaMock_FindMangaHistories_Call struct {
	*mock.Call
}

// FindMangaHistories is a helper method to define mock.On call
//   - userId string
//   - query *appdto.PagedQueryInput
func (_e *MangaMock_Expecter) FindMangaHistories(userId interface{}, query interface{}) *MangaMock_FindMangaHistories_Call {
	return &MangaMock_FindMangaHistories_Call{Call: _e.mock.On("FindMangaHistories", userId, query)}
}

func (_c *MangaMock_FindMangaHistories_Call) Run(run func(userId string, query *dto2.PagedQueryInput)) *MangaMock_FindMangaHistories_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*dto2.PagedQueryInput))
	})
	return _c
}

func (_c *MangaMock_FindMangaHistories_Call) Return(_a0 []dto.MangaHistoryResponse, _a1 *dto2.ResponsePage, _a2 status.Object) *MangaMock_FindMangaHistories_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MangaMock_FindMangaHistories_Call) RunAndReturn(run func(string, *dto2.PagedQueryInput) ([]dto.MangaHistoryResponse, *dto2.ResponsePage, status.Object)) *MangaMock_FindMangaHistories_Call {
	_c.Call.Return(run)
	return _c
}

// FindMangaRatings provides a mock function with given fields: mangaId
func (_m *MangaMock) FindMangaRatings(mangaId string) ([]dto.RateResponse, status.Object) {
	ret := _m.Called(mangaId)

	if len(ret) == 0 {
		panic("no return value specified for FindMangaRatings")
	}

	var r0 []dto.RateResponse
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(string) ([]dto.RateResponse, status.Object)); ok {
		return rf(mangaId)
	}
	if rf, ok := ret.Get(0).(func(string) []dto.RateResponse); ok {
		r0 = rf(mangaId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.RateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) status.Object); ok {
		r1 = rf(mangaId)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// MangaMock_FindMangaRatings_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMangaRatings'
type MangaMock_FindMangaRatings_Call struct {
	*mock.Call
}

// FindMangaRatings is a helper method to define mock.On call
//   - mangaId string
func (_e *MangaMock_Expecter) FindMangaRatings(mangaId interface{}) *MangaMock_FindMangaRatings_Call {
	return &MangaMock_FindMangaRatings_Call{Call: _e.mock.On("FindMangaRatings", mangaId)}
}

func (_c *MangaMock_FindMangaRatings_Call) Run(run func(mangaId string)) *MangaMock_FindMangaRatings_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MangaMock_FindMangaRatings_Call) Return(_a0 []dto.RateResponse, _a1 status.Object) *MangaMock_FindMangaRatings_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MangaMock_FindMangaRatings_Call) RunAndReturn(run func(string) ([]dto.RateResponse, status.Object)) *MangaMock_FindMangaRatings_Call {
	_c.Call.Return(run)
	return _c
}

// FindMangaTranslations provides a mock function with given fields: mangaId
func (_m *MangaMock) FindMangaTranslations(mangaId string) ([]dto.TranslationResponse, status.Object) {
	ret := _m.Called(mangaId)

	if len(ret) == 0 {
		panic("no return value specified for FindMangaTranslations")
	}

	var r0 []dto.TranslationResponse
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(string) ([]dto.TranslationResponse, status.Object)); ok {
		return rf(mangaId)
	}
	if rf, ok := ret.Get(0).(func(string) []dto.TranslationResponse); ok {
		r0 = rf(mangaId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.TranslationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) status.Object); ok {
		r1 = rf(mangaId)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// MangaMock_FindMangaTranslations_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMangaTranslations'
type MangaMock_FindMangaTranslations_Call struct {
	*mock.Call
}

// FindMangaTranslations is a helper method to define mock.On call
//   - mangaId string
func (_e *MangaMock_Expecter) FindMangaTranslations(mangaId interface{}) *MangaMock_FindMangaTranslations_Call {
	return &MangaMock_FindMangaTranslations_Call{Call: _e.mock.On("FindMangaTranslations", mangaId)}
}

func (_c *MangaMock_FindMangaTranslations_Call) Run(run func(mangaId string)) *MangaMock_FindMangaTranslations_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MangaMock_FindMangaTranslations_Call) Return(_a0 []dto.TranslationResponse, _a1 status.Object) *MangaMock_FindMangaTranslations_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MangaMock_FindMangaTranslations_Call) RunAndReturn(run func(string) ([]dto.TranslationResponse, status.Object)) *MangaMock_FindMangaTranslations_Call {
	_c.Call.Return(run)
	return _c
}

// FindRandomMangas provides a mock function with given fields: limit
func (_m *MangaMock) FindRandomMangas(limit uint64) ([]dto.MangaResponse, status.Object) {
	ret := _m.Called(limit)

	if len(ret) == 0 {
		panic("no return value specified for FindRandomMangas")
	}

	var r0 []dto.MangaResponse
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(uint64) ([]dto.MangaResponse, status.Object)); ok {
		return rf(limit)
	}
	if rf, ok := ret.Get(0).(func(uint64) []dto.MangaResponse); ok {
		r0 = rf(limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.MangaResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) status.Object); ok {
		r1 = rf(limit)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// MangaMock_FindRandomMangas_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindRandomMangas'
type MangaMock_FindRandomMangas_Call struct {
	*mock.Call
}

// FindRandomMangas is a helper method to define mock.On call
//   - limit uint64
func (_e *MangaMock_Expecter) FindRandomMangas(limit interface{}) *MangaMock_FindRandomMangas_Call {
	return &MangaMock_FindRandomMangas_Call{Call: _e.mock.On("FindRandomMangas", limit)}
}

func (_c *MangaMock_FindRandomMangas_Call) Run(run func(limit uint64)) *MangaMock_FindRandomMangas_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint64))
	})
	return _c
}

func (_c *MangaMock_FindRandomMangas_Call) Return(_a0 []dto.MangaResponse, _a1 status.Object) *MangaMock_FindRandomMangas_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MangaMock_FindRandomMangas_Call) RunAndReturn(run func(uint64) ([]dto.MangaResponse, status.Object)) *MangaMock_FindRandomMangas_Call {
	_c.Call.Return(run)
	return _c
}

// FindSpecificMangaTranslation provides a mock function with given fields: mangaId, language
func (_m *MangaMock) FindSpecificMangaTranslation(mangaId string, language common.Language) (dto.TranslationResponse, status.Object) {
	ret := _m.Called(mangaId, language)

	if len(ret) == 0 {
		panic("no return value specified for FindSpecificMangaTranslation")
	}

	var r0 dto.TranslationResponse
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(string, common.Language) (dto.TranslationResponse, status.Object)); ok {
		return rf(mangaId, language)
	}
	if rf, ok := ret.Get(0).(func(string, common.Language) dto.TranslationResponse); ok {
		r0 = rf(mangaId, language)
	} else {
		r0 = ret.Get(0).(dto.TranslationResponse)
	}

	if rf, ok := ret.Get(1).(func(string, common.Language) status.Object); ok {
		r1 = rf(mangaId, language)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// MangaMock_FindSpecificMangaTranslation_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindSpecificMangaTranslation'
type MangaMock_FindSpecificMangaTranslation_Call struct {
	*mock.Call
}

// FindSpecificMangaTranslation is a helper method to define mock.On call
//   - mangaId string
//   - language common.Language
func (_e *MangaMock_Expecter) FindSpecificMangaTranslation(mangaId interface{}, language interface{}) *MangaMock_FindSpecificMangaTranslation_Call {
	return &MangaMock_FindSpecificMangaTranslation_Call{Call: _e.mock.On("FindSpecificMangaTranslation", mangaId, language)}
}

func (_c *MangaMock_FindSpecificMangaTranslation_Call) Run(run func(mangaId string, language common.Language)) *MangaMock_FindSpecificMangaTranslation_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(common.Language))
	})
	return _c
}

func (_c *MangaMock_FindSpecificMangaTranslation_Call) Return(_a0 dto.TranslationResponse, _a1 status.Object) *MangaMock_FindSpecificMangaTranslation_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MangaMock_FindSpecificMangaTranslation_Call) RunAndReturn(run func(string, common.Language) (dto.TranslationResponse, status.Object)) *MangaMock_FindSpecificMangaTranslation_Call {
	_c.Call.Return(run)
	return _c
}

// InsertMangaTranslations provides a mock function with given fields: input
func (_m *MangaMock) InsertMangaTranslations(input *dto.MangaInsertTranslationInput) status.Object {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for InsertMangaTranslations")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(*dto.MangaInsertTranslationInput) status.Object); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// MangaMock_InsertMangaTranslations_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertMangaTranslations'
type MangaMock_InsertMangaTranslations_Call struct {
	*mock.Call
}

// InsertMangaTranslations is a helper method to define mock.On call
//   - input *dto.MangaInsertTranslationInput
func (_e *MangaMock_Expecter) InsertMangaTranslations(input interface{}) *MangaMock_InsertMangaTranslations_Call {
	return &MangaMock_InsertMangaTranslations_Call{Call: _e.mock.On("InsertMangaTranslations", input)}
}

func (_c *MangaMock_InsertMangaTranslations_Call) Run(run func(input *dto.MangaInsertTranslationInput)) *MangaMock_InsertMangaTranslations_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.MangaInsertTranslationInput))
	})
	return _c
}

func (_c *MangaMock_InsertMangaTranslations_Call) Return(_a0 status.Object) *MangaMock_InsertMangaTranslations_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MangaMock_InsertMangaTranslations_Call) RunAndReturn(run func(*dto.MangaInsertTranslationInput) status.Object) *MangaMock_InsertMangaTranslations_Call {
	_c.Call.Return(run)
	return _c
}

// ListMangas provides a mock function with given fields: query
func (_m *MangaMock) ListMangas(query *dto2.PagedQueryInput) ([]dto.MangaResponse, *dto2.ResponsePage, status.Object) {
	ret := _m.Called(query)

	if len(ret) == 0 {
		panic("no return value specified for ListMangas")
	}

	var r0 []dto.MangaResponse
	var r1 *dto2.ResponsePage
	var r2 status.Object
	if rf, ok := ret.Get(0).(func(*dto2.PagedQueryInput) ([]dto.MangaResponse, *dto2.ResponsePage, status.Object)); ok {
		return rf(query)
	}
	if rf, ok := ret.Get(0).(func(*dto2.PagedQueryInput) []dto.MangaResponse); ok {
		r0 = rf(query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.MangaResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*dto2.PagedQueryInput) *dto2.ResponsePage); ok {
		r1 = rf(query)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dto2.ResponsePage)
		}
	}

	if rf, ok := ret.Get(2).(func(*dto2.PagedQueryInput) status.Object); ok {
		r2 = rf(query)
	} else {
		r2 = ret.Get(2).(status.Object)
	}

	return r0, r1, r2
}

// MangaMock_ListMangas_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListMangas'
type MangaMock_ListMangas_Call struct {
	*mock.Call
}

// ListMangas is a helper method to define mock.On call
//   - query *appdto.PagedQueryInput
func (_e *MangaMock_Expecter) ListMangas(query interface{}) *MangaMock_ListMangas_Call {
	return &MangaMock_ListMangas_Call{Call: _e.mock.On("ListMangas", query)}
}

func (_c *MangaMock_ListMangas_Call) Run(run func(query *dto2.PagedQueryInput)) *MangaMock_ListMangas_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto2.PagedQueryInput))
	})
	return _c
}

func (_c *MangaMock_ListMangas_Call) Return(_a0 []dto.MangaResponse, _a1 *dto2.ResponsePage, _a2 status.Object) *MangaMock_ListMangas_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MangaMock_ListMangas_Call) RunAndReturn(run func(*dto2.PagedQueryInput) ([]dto.MangaResponse, *dto2.ResponsePage, status.Object)) *MangaMock_ListMangas_Call {
	_c.Call.Return(run)
	return _c
}

// SearchMangas provides a mock function with given fields: query
func (_m *MangaMock) SearchMangas(query *dto.MangaSearchQuery) ([]dto.MangaResponse, *dto2.ResponsePage, status.Object) {
	ret := _m.Called(query)

	if len(ret) == 0 {
		panic("no return value specified for SearchMangas")
	}

	var r0 []dto.MangaResponse
	var r1 *dto2.ResponsePage
	var r2 status.Object
	if rf, ok := ret.Get(0).(func(*dto.MangaSearchQuery) ([]dto.MangaResponse, *dto2.ResponsePage, status.Object)); ok {
		return rf(query)
	}
	if rf, ok := ret.Get(0).(func(*dto.MangaSearchQuery) []dto.MangaResponse); ok {
		r0 = rf(query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.MangaResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*dto.MangaSearchQuery) *dto2.ResponsePage); ok {
		r1 = rf(query)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dto2.ResponsePage)
		}
	}

	if rf, ok := ret.Get(2).(func(*dto.MangaSearchQuery) status.Object); ok {
		r2 = rf(query)
	} else {
		r2 = ret.Get(2).(status.Object)
	}

	return r0, r1, r2
}

// MangaMock_SearchMangas_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SearchMangas'
type MangaMock_SearchMangas_Call struct {
	*mock.Call
}

// SearchMangas is a helper method to define mock.On call
//   - query *dto.MangaSearchQuery
func (_e *MangaMock_Expecter) SearchMangas(query interface{}) *MangaMock_SearchMangas_Call {
	return &MangaMock_SearchMangas_Call{Call: _e.mock.On("SearchMangas", query)}
}

func (_c *MangaMock_SearchMangas_Call) Run(run func(query *dto.MangaSearchQuery)) *MangaMock_SearchMangas_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.MangaSearchQuery))
	})
	return _c
}

func (_c *MangaMock_SearchMangas_Call) Return(_a0 []dto.MangaResponse, _a1 *dto2.ResponsePage, _a2 status.Object) *MangaMock_SearchMangas_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MangaMock_SearchMangas_Call) RunAndReturn(run func(*dto.MangaSearchQuery) ([]dto.MangaResponse, *dto2.ResponsePage, status.Object)) *MangaMock_SearchMangas_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateMangaCover provides a mock function with given fields: input
func (_m *MangaMock) UpdateMangaCover(input *dto.MangaCoverUpdateInput) status.Object {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMangaCover")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(*dto.MangaCoverUpdateInput) status.Object); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// MangaMock_UpdateMangaCover_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMangaCover'
type MangaMock_UpdateMangaCover_Call struct {
	*mock.Call
}

// UpdateMangaCover is a helper method to define mock.On call
//   - input *dto.MangaCoverUpdateInput
func (_e *MangaMock_Expecter) UpdateMangaCover(input interface{}) *MangaMock_UpdateMangaCover_Call {
	return &MangaMock_UpdateMangaCover_Call{Call: _e.mock.On("UpdateMangaCover", input)}
}

func (_c *MangaMock_UpdateMangaCover_Call) Run(run func(input *dto.MangaCoverUpdateInput)) *MangaMock_UpdateMangaCover_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.MangaCoverUpdateInput))
	})
	return _c
}

func (_c *MangaMock_UpdateMangaCover_Call) Return(_a0 status.Object) *MangaMock_UpdateMangaCover_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MangaMock_UpdateMangaCover_Call) RunAndReturn(run func(*dto.MangaCoverUpdateInput) status.Object) *MangaMock_UpdateMangaCover_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateTranslation provides a mock function with given fields: input
func (_m *MangaMock) UpdateTranslation(input *dto.TranslationUpdateInput) status.Object {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for UpdateTranslation")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(*dto.TranslationUpdateInput) status.Object); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// MangaMock_UpdateTranslation_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateTranslation'
type MangaMock_UpdateTranslation_Call struct {
	*mock.Call
}

// UpdateTranslation is a helper method to define mock.On call
//   - input *dto.TranslationUpdateInput
func (_e *MangaMock_Expecter) UpdateTranslation(input interface{}) *MangaMock_UpdateTranslation_Call {
	return &MangaMock_UpdateTranslation_Call{Call: _e.mock.On("UpdateTranslation", input)}
}

func (_c *MangaMock_UpdateTranslation_Call) Run(run func(input *dto.TranslationUpdateInput)) *MangaMock_UpdateTranslation_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.TranslationUpdateInput))
	})
	return _c
}

func (_c *MangaMock_UpdateTranslation_Call) Return(_a0 status.Object) *MangaMock_UpdateTranslation_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MangaMock_UpdateTranslation_Call) RunAndReturn(run func(*dto.TranslationUpdateInput) status.Object) *MangaMock_UpdateTranslation_Call {
	_c.Call.Return(run)
	return _c
}

// UpsertMangaRating provides a mock function with given fields: input
func (_m *MangaMock) UpsertMangaRating(input *dto.RateUpsertInput) status.Object {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for UpsertMangaRating")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(*dto.RateUpsertInput) status.Object); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// MangaMock_UpsertMangaRating_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertMangaRating'
type MangaMock_UpsertMangaRating_Call struct {
	*mock.Call
}

// UpsertMangaRating is a helper method to define mock.On call
//   - input *dto.RateUpsertInput
func (_e *MangaMock_Expecter) UpsertMangaRating(input interface{}) *MangaMock_UpsertMangaRating_Call {
	return &MangaMock_UpsertMangaRating_Call{Call: _e.mock.On("UpsertMangaRating", input)}
}

func (_c *MangaMock_UpsertMangaRating_Call) Run(run func(input *dto.RateUpsertInput)) *MangaMock_UpsertMangaRating_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.RateUpsertInput))
	})
	return _c
}

func (_c *MangaMock_UpsertMangaRating_Call) Return(_a0 status.Object) *MangaMock_UpsertMangaRating_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MangaMock_UpsertMangaRating_Call) RunAndReturn(run func(*dto.RateUpsertInput) status.Object) *MangaMock_UpsertMangaRating_Call {
	_c.Call.Return(run)
	return _c
}

// NewMangaMock creates a new instance of MangaMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMangaMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *MangaMock {
	mock := &MangaMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}