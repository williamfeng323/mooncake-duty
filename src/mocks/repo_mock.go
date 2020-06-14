// Code generated by MockGen. DO NOT EDIT.
// Source: ./src/infrastructure/db/repo.go

// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
	reflect "reflect"
)

// MockValidator is a mock of Validator interface
type MockValidator struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorMockRecorder
}

// MockValidatorMockRecorder is the mock recorder for MockValidator
type MockValidatorMockRecorder struct {
	mock *MockValidator
}

// NewMockValidator creates a new mock instance
func NewMockValidator(ctrl *gomock.Controller) *MockValidator {
	mock := &MockValidator{ctrl: ctrl}
	mock.recorder = &MockValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockValidator) EXPECT() *MockValidatorMockRecorder {
	return m.recorder
}

// Verify mocks base method
func (m *MockValidator) Verify(arg0 interface{}) []error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", arg0)
	ret0, _ := ret[0].([]error)
	return ret0
}

// Verify indicates an expected call of Verify
func (mr *MockValidatorMockRecorder) Verify(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockValidator)(nil).Verify), arg0)
}

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetName mocks base method
func (m *MockRepository) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName
func (mr *MockRepositoryMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockRepository)(nil).GetName))
}

// SetCollection mocks base method
func (m *MockRepository) SetCollection(arg0 *mongo.Collection) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetCollection", arg0)
}

// SetCollection indicates an expected call of SetCollection
func (mr *MockRepositoryMockRecorder) SetCollection(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCollection", reflect.TypeOf((*MockRepository)(nil).SetCollection), arg0)
}

// InsertOne mocks base method
func (m *MockRepository) InsertOne(arg0 context.Context, arg1 interface{}, arg2 ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "InsertOne", varargs...)
	ret0, _ := ret[0].(*mongo.InsertOneResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertOne indicates an expected call of InsertOne
func (mr *MockRepositoryMockRecorder) InsertOne(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOne", reflect.TypeOf((*MockRepository)(nil).InsertOne), varargs...)
}

// Find mocks base method
func (m *MockRepository) Find(arg0 context.Context, arg1 interface{}, arg2 ...*options.FindOptions) (*mongo.Cursor, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Find", varargs...)
	ret0, _ := ret[0].(*mongo.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockRepositoryMockRecorder) Find(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockRepository)(nil).Find), varargs...)
}

// FindOne mocks base method
func (m *MockRepository) FindOne(arg0 context.Context, arg1 interface{}, arg2 ...*options.FindOneOptions) *mongo.SingleResult {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindOne", varargs...)
	ret0, _ := ret[0].(*mongo.SingleResult)
	return ret0
}

// FindOne indicates an expected call of FindOne
func (mr *MockRepositoryMockRecorder) FindOne(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockRepository)(nil).FindOne), varargs...)
}

// UpdateOne mocks base method
func (m *MockRepository) UpdateOne(arg0 context.Context, arg1, arg2 interface{}, arg3 ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateOne", varargs...)
	ret0, _ := ret[0].(*mongo.UpdateResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateOne indicates an expected call of UpdateOne
func (mr *MockRepositoryMockRecorder) UpdateOne(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOne", reflect.TypeOf((*MockRepository)(nil).UpdateOne), varargs...)
}

// DeleteOne mocks base method
func (m *MockRepository) DeleteOne(arg0 context.Context, arg1 interface{}, arg2 ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteOne", varargs...)
	ret0, _ := ret[0].(*mongo.DeleteResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteOne indicates an expected call of DeleteOne
func (mr *MockRepositoryMockRecorder) DeleteOne(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOne", reflect.TypeOf((*MockRepository)(nil).DeleteOne), varargs...)
}
