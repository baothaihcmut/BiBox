// Code generated by MockGen. DO NOT EDIT.
// Source: internal/modules/files/interactors/file.interactor.go
//
// Generated by this command:
//
//	mockgen -source=internal/modules/files/interactors/file.interactor.go -destination=test/files/mocks/mock_file.interactor.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	presenters "github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	gomock "go.uber.org/mock/gomock"
)

// MockFileInteractor is a mock of FileInteractor interface.
type MockFileInteractor struct {
	ctrl     *gomock.Controller
	recorder *MockFileInteractorMockRecorder
	isgomock struct{}
}

// MockFileInteractorMockRecorder is the mock recorder for MockFileInteractor.
type MockFileInteractorMockRecorder struct {
	mock *MockFileInteractor
}

// NewMockFileInteractor creates a new mock instance.
func NewMockFileInteractor(ctrl *gomock.Controller) *MockFileInteractor {
	mock := &MockFileInteractor{ctrl: ctrl}
	mock.recorder = &MockFileInteractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileInteractor) EXPECT() *MockFileInteractorMockRecorder {
	return m.recorder
}

// CreatFile mocks base method.
func (m *MockFileInteractor) CreatFile(arg0 context.Context, arg1 *presenters.CreateFileInput) (*presenters.CreateFileOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatFile", arg0, arg1)
	ret0, _ := ret[0].(*presenters.CreateFileOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatFile indicates an expected call of CreatFile.
func (mr *MockFileInteractorMockRecorder) CreatFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatFile", reflect.TypeOf((*MockFileInteractor)(nil).CreatFile), arg0, arg1)
}

// FindAllFileOfUser mocks base method.
func (m *MockFileInteractor) FindAllFileOfUser(ctx context.Context, input *presenters.FindFileOfUserInput) (*presenters.FindFileOfUserOuput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllFileOfUser", ctx, input)
	ret0, _ := ret[0].(*presenters.FindFileOfUserOuput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllFileOfUser indicates an expected call of FindAllFileOfUser.
func (mr *MockFileInteractorMockRecorder) FindAllFileOfUser(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllFileOfUser", reflect.TypeOf((*MockFileInteractor)(nil).FindAllFileOfUser), ctx, input)
}

// UploadedFile mocks base method.
func (m *MockFileInteractor) UploadedFile(arg0 context.Context, arg1 *presenters.UploadedFileInput) (*presenters.UploadedFileOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadedFile", arg0, arg1)
	ret0, _ := ret[0].(*presenters.UploadedFileOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadedFile indicates an expected call of UploadedFile.
func (mr *MockFileInteractorMockRecorder) UploadedFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadedFile", reflect.TypeOf((*MockFileInteractor)(nil).UploadedFile), arg0, arg1)
}
