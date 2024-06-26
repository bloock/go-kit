// Code generated by MockGen. DO NOT EDIT.
// Source: http/http_request.go

// Package mock_http is a generated GoMock package.
package mock_http

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHttpRequest is a mock of HttpRequest interface.
type MockHttpRequest struct {
	ctrl     *gomock.Controller
	recorder *MockHttpRequestMockRecorder
}

// MockHttpRequestMockRecorder is the mock recorder for MockHttpRequest.
type MockHttpRequestMockRecorder struct {
	mock *MockHttpRequest
}

// NewMockHttpRequest creates a new mock instance.
func NewMockHttpRequest(ctrl *gomock.Controller) *MockHttpRequest {
	mock := &MockHttpRequest{ctrl: ctrl}
	mock.recorder = &MockHttpRequestMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpRequest) EXPECT() *MockHttpRequestMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockHttpRequest) Delete(ctx context.Context, url string, body, response interface{}, headers map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, url, body, response, headers)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockHttpRequestMockRecorder) Delete(ctx, url, body, response, headers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockHttpRequest)(nil).Delete), ctx, url, body, response, headers)
}

// Get mocks base method.
func (m *MockHttpRequest) Get(ctx context.Context, url string, response interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, url, response)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockHttpRequestMockRecorder) Get(ctx, url, response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockHttpRequest)(nil).Get), ctx, url, response)
}

// GetWithHeaders mocks base method.
func (m *MockHttpRequest) GetWithHeaders(ctx context.Context, url string, response interface{}, headers map[string][]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithHeaders", ctx, url, response, headers)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetWithHeaders indicates an expected call of GetWithHeaders.
func (mr *MockHttpRequestMockRecorder) GetWithHeaders(ctx, url, response, headers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithHeaders", reflect.TypeOf((*MockHttpRequest)(nil).GetWithHeaders), ctx, url, response, headers)
}

// Post mocks base method.
func (m *MockHttpRequest) Post(ctx context.Context, url string, body, response interface{}, contentType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Post", ctx, url, body, response, contentType)
	ret0, _ := ret[0].(error)
	return ret0
}

// Post indicates an expected call of Post.
func (mr *MockHttpRequestMockRecorder) Post(ctx, url, body, response, contentType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockHttpRequest)(nil).Post), ctx, url, body, response, contentType)
}

// PostWithHeaders mocks base method.
func (m *MockHttpRequest) PostWithHeaders(ctx context.Context, url string, body, response interface{}, headers map[string]string, contentType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostWithHeaders", ctx, url, body, response, headers, contentType)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostWithHeaders indicates an expected call of PostWithHeaders.
func (mr *MockHttpRequestMockRecorder) PostWithHeaders(ctx, url, body, response, headers, contentType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostWithHeaders", reflect.TypeOf((*MockHttpRequest)(nil).PostWithHeaders), ctx, url, body, response, headers, contentType)
}
