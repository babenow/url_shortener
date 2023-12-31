// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/babenow/url_shortener/intrernal/model"
	mock "github.com/stretchr/testify/mock"
)

// URLGetter is an autogenerated mock type for the URLGetter type
type URLGetter struct {
	mock.Mock
}

// AddRedirect provides a mock function with given fields: ctx, alias
func (_m *URLGetter) AddRedirect(ctx context.Context, alias string) error {
	ret := _m.Called(ctx, alias)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, alias)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetURLByAlias provides a mock function with given fields: ctx, alias
func (_m *URLGetter) GetURLByAlias(ctx context.Context, alias string) (*model.Url, error) {
	ret := _m.Called(ctx, alias)

	var r0 *model.Url
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Url, error)); ok {
		return rf(ctx, alias)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Url); ok {
		r0 = rf(ctx, alias)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Url)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewURLGetter creates a new instance of URLGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewURLGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *URLGetter {
	mock := &URLGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
