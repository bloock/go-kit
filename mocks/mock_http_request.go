// Code generated by MockGen. DO NOT EDIT.
// Source: request/http_request.go

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
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

// Get mocks base method.
func (m *MockHttpRequest) Get(url string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", url)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockHttpRequestMockRecorder) Get(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockHttpRequest)(nil).Get), url)
}
