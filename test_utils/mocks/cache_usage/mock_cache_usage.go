// Code generated by MockGen. DO NOT EDIT.
// Source: cache/cache_usage_repository.go

// Package mock_cache is a generated GoMock package.
package mock_cache

import (
	context "context"
	reflect "reflect"

	domain "github.com/bloock/go-kit/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockCacheUsageRepository is a mock of CacheUsageRepository interface.
type MockCacheUsageRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCacheUsageRepositoryMockRecorder
}

// MockCacheUsageRepositoryMockRecorder is the mock recorder for MockCacheUsageRepository.
type MockCacheUsageRepositoryMockRecorder struct {
	mock *MockCacheUsageRepository
}

// NewMockCacheUsageRepository creates a new mock instance.
func NewMockCacheUsageRepository(ctrl *gomock.Controller) *MockCacheUsageRepository {
	mock := &MockCacheUsageRepository{ctrl: ctrl}
	mock.recorder = &MockCacheUsageRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacheUsageRepository) EXPECT() *MockCacheUsageRepositoryMockRecorder {
	return m.recorder
}

// FindValueByKey mocks base method.
func (m *MockCacheUsageRepository) FindValueByKey(ctx context.Context, key string) (domain.CacheUsage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindValueByKey", ctx, key)
	ret0, _ := ret[0].(domain.CacheUsage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindValueByKey indicates an expected call of FindValueByKey.
func (mr *MockCacheUsageRepositoryMockRecorder) FindValueByKey(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindValueByKey", reflect.TypeOf((*MockCacheUsageRepository)(nil).FindValueByKey), ctx, key)
}

// GetValueByKey mocks base method.
func (m *MockCacheUsageRepository) GetValueByKey(ctx context.Context, key string) (domain.CacheUsage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValueByKey", ctx, key)
	ret0, _ := ret[0].(domain.CacheUsage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValueByKey indicates an expected call of GetValueByKey.
func (mr *MockCacheUsageRepositoryMockRecorder) GetValueByKey(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValueByKey", reflect.TypeOf((*MockCacheUsageRepository)(nil).GetValueByKey), ctx, key)
}

// Save mocks base method.
func (m *MockCacheUsageRepository) Save(ctx context.Context, usage domain.CacheUsage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, usage)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockCacheUsageRepositoryMockRecorder) Save(ctx, usage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockCacheUsageRepository)(nil).Save), ctx, usage)
}

// Update mocks base method.
func (m *MockCacheUsageRepository) Update(ctx context.Context, usage domain.CacheUsage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, usage)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCacheUsageRepositoryMockRecorder) Update(ctx, usage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCacheUsageRepository)(nil).Update), ctx, usage)
}
