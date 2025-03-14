// Code generated by MockGen. DO NOT EDIT.
// Source: internal/common/storage/storage.service.go
//
// Generated by this command:
//
//	mockgen -source=internal/common/storage/storage.service.go -destination=test/mocks/mock_storage.service.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	io "io"
	reflect "reflect"

	storage "github.com/baothaihcmut/Bibox/storage-app/internal/common/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockStorageService is a mock of StorageService interface.
type MockStorageService struct {
	ctrl     *gomock.Controller
	recorder *MockStorageServiceMockRecorder
	isgomock struct{}
}

// MockStorageServiceMockRecorder is the mock recorder for MockStorageService.
type MockStorageServiceMockRecorder struct {
	mock *MockStorageService
}

// NewMockStorageService creates a new mock instance.
func NewMockStorageService(ctrl *gomock.Controller) *MockStorageService {
	mock := &MockStorageService{ctrl: ctrl}
	mock.recorder = &MockStorageServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageService) EXPECT() *MockStorageServiceMockRecorder {
	return m.recorder
}

// GetFile mocks base method.
func (m *MockStorageService) GetFile(arg0 context.Context, arg1 string) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFile", arg0, arg1)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFile indicates an expected call of GetFile.
func (mr *MockStorageServiceMockRecorder) GetFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFile", reflect.TypeOf((*MockStorageService)(nil).GetFile), arg0, arg1)
}

// GetPresignUrl mocks base method.
func (m *MockStorageService) GetPresignUrl(arg0 context.Context, arg1 storage.GetPresignUrlArg) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPresignUrl", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPresignUrl indicates an expected call of GetPresignUrl.
func (mr *MockStorageServiceMockRecorder) GetPresignUrl(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPresignUrl", reflect.TypeOf((*MockStorageService)(nil).GetPresignUrl), arg0, arg1)
}

// GetStorageBucket mocks base method.
func (m *MockStorageService) GetStorageBucket() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStorageBucket")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetStorageBucket indicates an expected call of GetStorageBucket.
func (mr *MockStorageServiceMockRecorder) GetStorageBucket() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStorageBucket", reflect.TypeOf((*MockStorageService)(nil).GetStorageBucket))
}

// GetStorageProviderName mocks base method.
func (m *MockStorageService) GetStorageProviderName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStorageProviderName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetStorageProviderName indicates an expected call of GetStorageProviderName.
func (mr *MockStorageServiceMockRecorder) GetStorageProviderName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStorageProviderName", reflect.TypeOf((*MockStorageService)(nil).GetStorageProviderName))
}
