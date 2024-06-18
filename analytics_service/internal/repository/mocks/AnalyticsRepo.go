// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	domain "analytics_service/internal/domain"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// AnalyticsRepo is an autogenerated mock type for the AnalyticsRepo type
type AnalyticsRepo struct {
	mock.Mock
}

// GetTopUrls provides a mock function with given fields: ctx, paginationParams
func (_m *AnalyticsRepo) GetTopUrls(ctx context.Context, paginationParams domain.PaginationParams) ([]domain.TopURLData, error) {
	ret := _m.Called(ctx, paginationParams)

	if len(ret) == 0 {
		panic("no return value specified for GetTopUrls")
	}

	var r0 []domain.TopURLData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.PaginationParams) ([]domain.TopURLData, error)); ok {
		return rf(ctx, paginationParams)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.PaginationParams) []domain.TopURLData); ok {
		r0 = rf(ctx, paginationParams)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.TopURLData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.PaginationParams) error); ok {
		r1 = rf(ctx, paginationParams)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAnalyticsRepo creates a new instance of AnalyticsRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAnalyticsRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *AnalyticsRepo {
	mock := &AnalyticsRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}