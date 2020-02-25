// Code generated by MockGen. DO NOT EDIT.
// Source: types/interfaces.go

// Package mock_types is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/pisign/pisign-backend/types"
)

// MockDataObject is a mock of DataObject interface
type MockDataObject struct {
	ctrl     *gomock.Controller
	recorder *MockDataObjectMockRecorder
}

// MockDataObjectMockRecorder is the mock recorder for MockDataObject
type MockDataObjectMockRecorder struct {
	mock *MockDataObject
}

// NewMockDataObject creates a new mock instance
func NewMockDataObject(ctrl *gomock.Controller) *MockDataObject {
	mock := &MockDataObject{ctrl: ctrl}
	mock.recorder = &MockDataObjectMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataObject) EXPECT() *MockDataObjectMockRecorder {
	return m.recorder
}

// Update mocks base method
func (m *MockDataObject) Update(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockDataObjectMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDataObject)(nil).Update), arg0)
}

// Transform mocks base method
func (m *MockDataObject) Transform() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transform")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Transform indicates an expected call of Transform
func (mr *MockDataObjectMockRecorder) Transform() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transform", reflect.TypeOf((*MockDataObject)(nil).Transform))
}

// MockAPI is a mock of API interface
type MockAPI struct {
	ctrl     *gomock.Controller
	recorder *MockAPIMockRecorder
}

// MockAPIMockRecorder is the mock recorder for MockAPI
type MockAPIMockRecorder struct {
	mock *MockAPI
}

// NewMockAPI creates a new mock instance
func NewMockAPI(ctrl *gomock.Controller) *MockAPI {
	mock := &MockAPI{ctrl: ctrl}
	mock.recorder = &MockAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAPI) EXPECT() *MockAPIMockRecorder {
	return m.recorder
}

// Configure mocks base method
func (m *MockAPI) Configure(message types.ClientMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Configure", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Configure indicates an expected call of Configure
func (mr *MockAPIMockRecorder) Configure(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configure", reflect.TypeOf((*MockAPI)(nil).Configure), message)
}

// Run mocks base method
func (m *MockAPI) Run(w types.Socket) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", w)
}

// Run indicates an expected call of Run
func (mr *MockAPIMockRecorder) Run(w interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockAPI)(nil).Run), w)
}

// Data mocks base method
func (m *MockAPI) Data() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Data")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Data indicates an expected call of Data
func (mr *MockAPIMockRecorder) Data() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Data", reflect.TypeOf((*MockAPI)(nil).Data))
}

// GetName mocks base method
func (m *MockAPI) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName
func (mr *MockAPIMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockAPI)(nil).GetName))
}

// MockSocket is a mock of Socket interface
type MockSocket struct {
	ctrl     *gomock.Controller
	recorder *MockSocketMockRecorder
}

// MockSocketMockRecorder is the mock recorder for MockSocket
type MockSocketMockRecorder struct {
	mock *MockSocket
}

// NewMockSocket creates a new mock instance
func NewMockSocket(ctrl *gomock.Controller) *MockSocket {
	mock := &MockSocket{ctrl: ctrl}
	mock.recorder = &MockSocketMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSocket) EXPECT() *MockSocketMockRecorder {
	return m.recorder
}

// Read mocks base method
func (m *MockSocket) Read() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Read")
}

// Read indicates an expected call of Read
func (mr *MockSocketMockRecorder) Read() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockSocket)(nil).Read))
}

// Send mocks base method
func (m *MockSocket) Send(arg0 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Send", arg0)
}

// Send indicates an expected call of Send
func (mr *MockSocketMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSocket)(nil).Send), arg0)
}

// Close mocks base method
func (m *MockSocket) Close() chan bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(chan bool)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockSocketMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSocket)(nil).Close))
}

// SendErrorMessage mocks base method
func (m *MockSocket) SendErrorMessage(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendErrorMessage", arg0)
}

// SendErrorMessage indicates an expected call of SendErrorMessage
func (mr *MockSocketMockRecorder) SendErrorMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendErrorMessage", reflect.TypeOf((*MockSocket)(nil).SendErrorMessage), arg0)
}

// MockPool is a mock of Pool interface
type MockPool struct {
	ctrl     *gomock.Controller
	recorder *MockPoolMockRecorder
}

// MockPoolMockRecorder is the mock recorder for MockPool
type MockPoolMockRecorder struct {
	mock *MockPool
}

// NewMockPool creates a new mock instance
func NewMockPool(ctrl *gomock.Controller) *MockPool {
	mock := &MockPool{ctrl: ctrl}
	mock.recorder = &MockPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPool) EXPECT() *MockPoolMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockPool) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start
func (mr *MockPoolMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockPool)(nil).Start))
}

// Register mocks base method
func (m *MockPool) Register(arg0 types.API) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Register", arg0)
}

// Register indicates an expected call of Register
func (mr *MockPoolMockRecorder) Register(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockPool)(nil).Register), arg0)
}

// Unregister mocks base method
func (m *MockPool) Unregister(arg0 types.API) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Unregister", arg0)
}

// Unregister indicates an expected call of Unregister
func (mr *MockPoolMockRecorder) Unregister(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unregister", reflect.TypeOf((*MockPool)(nil).Unregister), arg0)
}

// Save mocks base method
func (m *MockPool) Save() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Save")
}

// Save indicates an expected call of Save
func (mr *MockPoolMockRecorder) Save() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockPool)(nil).Save))
}