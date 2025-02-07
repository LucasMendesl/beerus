// Code generated by MockGen. DO NOT EDIT.
// Source: types.go
//
// Generated by this command:
//
//	mockgen -source types.go -package docker -destination ./mocks/beerus_container_api.go BeerusContainerAPI
//

// Package docker is a generated GoMock package.
package docker

import (
	context "context"
	reflect "reflect"

	types "github.com/docker/docker/api/types"
	events "github.com/docker/docker/api/types/events"
	docker "github.com/lucasmendesl/beerus/docker"
	gomock "go.uber.org/mock/gomock"
)

// MockBeerusContainerAPI is a mock of BeerusContainerAPI interface.
type MockBeerusContainerAPI struct {
	ctrl     *gomock.Controller
	recorder *MockBeerusContainerAPIMockRecorder
	isgomock struct{}
}

// MockBeerusContainerAPIMockRecorder is the mock recorder for MockBeerusContainerAPI.
type MockBeerusContainerAPIMockRecorder struct {
	mock *MockBeerusContainerAPI
}

// NewMockBeerusContainerAPI creates a new mock instance.
func NewMockBeerusContainerAPI(ctrl *gomock.Controller) *MockBeerusContainerAPI {
	mock := &MockBeerusContainerAPI{ctrl: ctrl}
	mock.recorder = &MockBeerusContainerAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBeerusContainerAPI) EXPECT() *MockBeerusContainerAPIMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockBeerusContainerAPI) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockBeerusContainerAPIMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBeerusContainerAPI)(nil).Close))
}

// FromEvents mocks base method.
func (m *MockBeerusContainerAPI) FromEvents(ctx context.Context, actions ...events.Action) <-chan docker.EventResult {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range actions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FromEvents", varargs...)
	ret0, _ := ret[0].(<-chan docker.EventResult)
	return ret0
}

// FromEvents indicates an expected call of FromEvents.
func (mr *MockBeerusContainerAPIMockRecorder) FromEvents(ctx any, actions ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, actions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FromEvents", reflect.TypeOf((*MockBeerusContainerAPI)(nil).FromEvents), varargs...)
}

// Inspect mocks base method.
func (m *MockBeerusContainerAPI) Inspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Inspect", ctx, containerID)
	ret0, _ := ret[0].(types.ContainerJSON)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Inspect indicates an expected call of Inspect.
func (mr *MockBeerusContainerAPIMockRecorder) Inspect(ctx, containerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inspect", reflect.TypeOf((*MockBeerusContainerAPI)(nil).Inspect), ctx, containerID)
}

// ListContainers mocks base method.
func (m *MockBeerusContainerAPI) ListContainers(ctx context.Context, concurrency uint8, options ...docker.ListContainersOptions) ([]docker.Container, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, concurrency}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListContainers", varargs...)
	ret0, _ := ret[0].([]docker.Container)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListContainers indicates an expected call of ListContainers.
func (mr *MockBeerusContainerAPIMockRecorder) ListContainers(ctx, concurrency any, options ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, concurrency}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListContainers", reflect.TypeOf((*MockBeerusContainerAPI)(nil).ListContainers), varargs...)
}

// ListExpiredImages mocks base method.
func (m *MockBeerusContainerAPI) ListExpiredImages(ctx context.Context, options docker.ExpiredImageListOptions) ([]docker.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListExpiredImages", ctx, options)
	ret0, _ := ret[0].([]docker.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListExpiredImages indicates an expected call of ListExpiredImages.
func (mr *MockBeerusContainerAPIMockRecorder) ListExpiredImages(ctx, options any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListExpiredImages", reflect.TypeOf((*MockBeerusContainerAPI)(nil).ListExpiredImages), ctx, options)
}

// RemoveContainer mocks base method.
func (m *MockBeerusContainerAPI) RemoveContainer(ctx context.Context, options docker.RemoveContainerOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveContainer", ctx, options)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveContainer indicates an expected call of RemoveContainer.
func (mr *MockBeerusContainerAPIMockRecorder) RemoveContainer(ctx, options any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveContainer", reflect.TypeOf((*MockBeerusContainerAPI)(nil).RemoveContainer), ctx, options)
}

// RemoveImage mocks base method.
func (m *MockBeerusContainerAPI) RemoveImage(ctx context.Context, dockerImage string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveImage", ctx, dockerImage)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveImage indicates an expected call of RemoveImage.
func (mr *MockBeerusContainerAPIMockRecorder) RemoveImage(ctx, dockerImage any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveImage", reflect.TypeOf((*MockBeerusContainerAPI)(nil).RemoveImage), ctx, dockerImage)
}
