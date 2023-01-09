// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/moira-alert/moira/logging (interfaces: EventBuilder)

// Package mock_moira_alert is a generated GoMock package.
package mock_moira_alert

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	logging "github.com/moira-alert/moira/logging"
)

// MockEventBuilder is a mock of EventBuilder interface.
type MockEventBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockEventBuilderMockRecorder
}

// MockEventBuilderMockRecorder is the mock recorder for MockEventBuilder.
type MockEventBuilderMockRecorder struct {
	mock *MockEventBuilder
}

// NewMockEventBuilder creates a new mock instance.
func NewMockEventBuilder(ctrl *gomock.Controller) *MockEventBuilder {
	mock := &MockEventBuilder{ctrl: ctrl}
	mock.recorder = &MockEventBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventBuilder) EXPECT() *MockEventBuilderMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *MockEventBuilder) Error(arg0 error) logging.EventBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Error", arg0)
	ret0, _ := ret[0].(logging.EventBuilder)
	return ret0
}

// Error indicates an expected call of Error.
func (mr *MockEventBuilderMockRecorder) Error(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockEventBuilder)(nil).Error), arg0)
}

// Fields mocks base method.
func (m *MockEventBuilder) Fields(arg0 map[string]interface{}) logging.EventBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fields", arg0)
	ret0, _ := ret[0].(logging.EventBuilder)
	return ret0
}

// Fields indicates an expected call of Fields.
func (mr *MockEventBuilderMockRecorder) Fields(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fields", reflect.TypeOf((*MockEventBuilder)(nil).Fields), arg0)
}

// Int mocks base method.
func (m *MockEventBuilder) Int(arg0 string, arg1 int) logging.EventBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Int", arg0, arg1)
	ret0, _ := ret[0].(logging.EventBuilder)
	return ret0
}

// Int indicates an expected call of Int.
func (mr *MockEventBuilderMockRecorder) Int(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Int", reflect.TypeOf((*MockEventBuilder)(nil).Int), arg0, arg1)
}

// Int64 mocks base method.
func (m *MockEventBuilder) Int64(arg0 string, arg1 int64) logging.EventBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Int64", arg0, arg1)
	ret0, _ := ret[0].(logging.EventBuilder)
	return ret0
}

// Int64 indicates an expected call of Int64.
func (mr *MockEventBuilderMockRecorder) Int64(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Int64", reflect.TypeOf((*MockEventBuilder)(nil).Int64), arg0, arg1)
}

// Interface mocks base method.
func (m *MockEventBuilder) Interface(arg0 string, arg1 interface{}) logging.EventBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Interface", arg0, arg1)
	ret0, _ := ret[0].(logging.EventBuilder)
	return ret0
}

// Interface indicates an expected call of Interface.
func (mr *MockEventBuilderMockRecorder) Interface(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Interface", reflect.TypeOf((*MockEventBuilder)(nil).Interface), arg0, arg1)
}

// Msg mocks base method.
func (m *MockEventBuilder) Msg(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Msg", arg0)
}

// Msg indicates an expected call of Msg.
func (mr *MockEventBuilderMockRecorder) Msg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Msg", reflect.TypeOf((*MockEventBuilder)(nil).Msg), arg0)
}

// String mocks base method.
func (m *MockEventBuilder) String(arg0, arg1 string) logging.EventBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String", arg0, arg1)
	ret0, _ := ret[0].(logging.EventBuilder)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockEventBuilderMockRecorder) String(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockEventBuilder)(nil).String), arg0, arg1)
}
