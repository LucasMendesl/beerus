// Code generated by MockGen. DO NOT EDIT.
// Source: client.go
//
// Generated by this command:
//
//	mockgen -source client.go -package docker -destination ./mocks/client.go Client
//

// Package docker is a generated GoMock package.
package docker

import (
	context "context"
	reflect "reflect"

	image "github.com/docker/docker/api/types/image"
	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
	isgomock struct{}
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// ImageList mocks base method.
func (m *MockClient) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageList", ctx, options)
	ret0, _ := ret[0].([]image.Summary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageList indicates an expected call of ImageList.
func (mr *MockClientMockRecorder) ImageList(ctx, options any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageList", reflect.TypeOf((*MockClient)(nil).ImageList), ctx, options)
}

// ImageRemove mocks base method.
func (m *MockClient) ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageRemove", ctx, imageID, options)
	ret0, _ := ret[0].([]image.DeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageRemove indicates an expected call of ImageRemove.
func (mr *MockClientMockRecorder) ImageRemove(ctx, imageID, options any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageRemove", reflect.TypeOf((*MockClient)(nil).ImageRemove), ctx, imageID, options)
}
