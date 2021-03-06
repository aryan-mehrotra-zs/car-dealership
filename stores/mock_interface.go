// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package stores is a generated GoMock package.
package stores

import (
	reflect "reflect"

	filters "github.com/amehrotra/car-dealership/filters"
	models "github.com/amehrotra/car-dealership/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockCar is a mock of Car interface.
type MockCar struct {
	ctrl     *gomock.Controller
	recorder *MockCarMockRecorder
}

// MockCarMockRecorder is the mock recorder for MockCar.
type MockCarMockRecorder struct {
	mock *MockCar
}

// NewMockCar creates a new mock instance.
func NewMockCar(ctrl *gomock.Controller) *MockCar {
	mock := &MockCar{ctrl: ctrl}
	mock.recorder = &MockCarMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCar) EXPECT() *MockCarMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCar) Create(car *models.Car) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", car)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCarMockRecorder) Create(car interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCar)(nil).Create), car)
}

// Delete mocks base method.
func (m *MockCar) Delete(id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCarMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCar)(nil).Delete), id)
}

// GetAll mocks base method.
func (m *MockCar) GetAll(filter filters.Car) ([]models.Car, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", filter)
	ret0, _ := ret[0].([]models.Car)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockCarMockRecorder) GetAll(filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockCar)(nil).GetAll), filter)
}

// GetByID mocks base method.
func (m *MockCar) GetByID(id uuid.UUID) (models.Car, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(models.Car)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockCarMockRecorder) GetByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockCar)(nil).GetByID), id)
}

// Update mocks base method.
func (m *MockCar) Update(car *models.Car) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", car)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCarMockRecorder) Update(car interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCar)(nil).Update), car)
}

// MockEngine is a mock of Engine interface.
type MockEngine struct {
	ctrl     *gomock.Controller
	recorder *MockEngineMockRecorder
}

// MockEngineMockRecorder is the mock recorder for MockEngine.
type MockEngineMockRecorder struct {
	mock *MockEngine
}

// NewMockEngine creates a new mock instance.
func NewMockEngine(ctrl *gomock.Controller) *MockEngine {
	mock := &MockEngine{ctrl: ctrl}
	mock.recorder = &MockEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEngine) EXPECT() *MockEngineMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockEngine) Create(engine *models.Engine) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", engine)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockEngineMockRecorder) Create(engine interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEngine)(nil).Create), engine)
}

// Delete mocks base method.
func (m *MockEngine) Delete(id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockEngineMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockEngine)(nil).Delete), id)
}

// GetByID mocks base method.
func (m *MockEngine) GetByID(id uuid.UUID) (models.Engine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(models.Engine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockEngineMockRecorder) GetByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockEngine)(nil).GetByID), id)
}

// Update mocks base method.
func (m *MockEngine) Update(engine *models.Engine) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", engine)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockEngineMockRecorder) Update(engine interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockEngine)(nil).Update), engine)
}
