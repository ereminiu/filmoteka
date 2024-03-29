// Code generated by MockGen. DO NOT EDIT.
// Source: repositories.go

// Package mock_db is a generated GoMock package.
package mock_db

import (
	reflect "reflect"

	models "github.com/ereminiu/filmoteka/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockMovie is a mock of Movie interface.
type MockMovie struct {
	ctrl     *gomock.Controller
	recorder *MockMovieMockRecorder
}

// MockMovieMockRecorder is the mock recorder for MockMovie.
type MockMovieMockRecorder struct {
	mock *MockMovie
}

// NewMockMovie creates a new mock instance.
func NewMockMovie(ctrl *gomock.Controller) *MockMovie {
	mock := &MockMovie{ctrl: ctrl}
	mock.recorder = &MockMovieMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMovie) EXPECT() *MockMovieMockRecorder {
	return m.recorder
}

// AddActorToMovie mocks base method.
func (m *MockMovie) AddActorToMovie(actorId, movieId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddActorToMovie", actorId, movieId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddActorToMovie indicates an expected call of AddActorToMovie.
func (mr *MockMovieMockRecorder) AddActorToMovie(actorId, movieId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddActorToMovie", reflect.TypeOf((*MockMovie)(nil).AddActorToMovie), actorId, movieId)
}

// ChangeField mocks base method.
func (m *MockMovie) ChangeField(movieId int, field, newValue string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeField", movieId, field, newValue)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeField indicates an expected call of ChangeField.
func (mr *MockMovieMockRecorder) ChangeField(movieId, field, newValue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeField", reflect.TypeOf((*MockMovie)(nil).ChangeField), movieId, field, newValue)
}

// CreateMovie mocks base method.
func (m *MockMovie) CreateMovie(name, description, date string, rate int, actorIds []int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMovie", name, description, date, rate, actorIds)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMovie indicates an expected call of CreateMovie.
func (mr *MockMovieMockRecorder) CreateMovie(name, description, date, rate, actorIds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMovie", reflect.TypeOf((*MockMovie)(nil).CreateMovie), name, description, date, rate, actorIds)
}

// DeleteField mocks base method.
func (m *MockMovie) DeleteField(movieId int, field any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteField", movieId, field)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteField indicates an expected call of DeleteField.
func (mr *MockMovieMockRecorder) DeleteField(movieId, field interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteField", reflect.TypeOf((*MockMovie)(nil).DeleteField), movieId, field)
}

// DeleteMovie mocks base method.
func (m *MockMovie) DeleteMovie(movieId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMovie", movieId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMovie indicates an expected call of DeleteMovie.
func (mr *MockMovieMockRecorder) DeleteMovie(movieId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMovie", reflect.TypeOf((*MockMovie)(nil).DeleteMovie), movieId)
}

// GetAllMovies mocks base method.
func (m *MockMovie) GetAllMovies(sortBy string) ([]models.MovieWithActors, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMovies", sortBy)
	ret0, _ := ret[0].([]models.MovieWithActors)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMovies indicates an expected call of GetAllMovies.
func (mr *MockMovieMockRecorder) GetAllMovies(sortBy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMovies", reflect.TypeOf((*MockMovie)(nil).GetAllMovies), sortBy)
}

// SearchMovie mocks base method.
func (m *MockMovie) SearchMovie(moviePattern, actorPattern string) ([]models.MovieWithActors, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchMovie", moviePattern, actorPattern)
	ret0, _ := ret[0].([]models.MovieWithActors)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchMovie indicates an expected call of SearchMovie.
func (mr *MockMovieMockRecorder) SearchMovie(moviePattern, actorPattern interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchMovie", reflect.TypeOf((*MockMovie)(nil).SearchMovie), moviePattern, actorPattern)
}

// MockActor is a mock of Actor interface.
type MockActor struct {
	ctrl     *gomock.Controller
	recorder *MockActorMockRecorder
}

// MockActorMockRecorder is the mock recorder for MockActor.
type MockActorMockRecorder struct {
	mock *MockActor
}

// NewMockActor creates a new mock instance.
func NewMockActor(ctrl *gomock.Controller) *MockActor {
	mock := &MockActor{ctrl: ctrl}
	mock.recorder = &MockActorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockActor) EXPECT() *MockActorMockRecorder {
	return m.recorder
}

// ChangeField mocks base method.
func (m *MockActor) ChangeField(field, newValue string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeField", field, newValue)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeField indicates an expected call of ChangeField.
func (mr *MockActorMockRecorder) ChangeField(field, newValue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeField", reflect.TypeOf((*MockActor)(nil).ChangeField), field, newValue)
}

// CreateActor mocks base method.
func (m *MockActor) CreateActor(name, gender, birthday string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateActor", name, gender, birthday)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateActor indicates an expected call of CreateActor.
func (mr *MockActorMockRecorder) CreateActor(name, gender, birthday interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateActor", reflect.TypeOf((*MockActor)(nil).CreateActor), name, gender, birthday)
}

// DeleteActor mocks base method.
func (m *MockActor) DeleteActor(actorId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteActor", actorId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteActor indicates an expected call of DeleteActor.
func (mr *MockActorMockRecorder) DeleteActor(actorId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteActor", reflect.TypeOf((*MockActor)(nil).DeleteActor), actorId)
}

// DeleteField mocks base method.
func (m *MockActor) DeleteField(field string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteField", field)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteField indicates an expected call of DeleteField.
func (mr *MockActorMockRecorder) DeleteField(field interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteField", reflect.TypeOf((*MockActor)(nil).DeleteField), field)
}

// GetAllActors mocks base method.
func (m *MockActor) GetAllActors() ([]models.ActorWithMovies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllActors")
	ret0, _ := ret[0].([]models.ActorWithMovies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllActors indicates an expected call of GetAllActors.
func (mr *MockActorMockRecorder) GetAllActors() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllActors", reflect.TypeOf((*MockActor)(nil).GetAllActors))
}

// MockAuthorization is a mock of Authorization interface.
type MockAuthorization struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationMockRecorder
}

// MockAuthorizationMockRecorder is the mock recorder for MockAuthorization.
type MockAuthorizationMockRecorder struct {
	mock *MockAuthorization
}

// NewMockAuthorization creates a new mock instance.
func NewMockAuthorization(ctrl *gomock.Controller) *MockAuthorization {
	mock := &MockAuthorization{ctrl: ctrl}
	mock.recorder = &MockAuthorizationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorization) EXPECT() *MockAuthorizationMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockAuthorization) CreateUser(name, username, passwordHash string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", name, username, passwordHash)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthorizationMockRecorder) CreateUser(name, username, passwordHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthorization)(nil).CreateUser), name, username, passwordHash)
}

// GetUser mocks base method.
func (m *MockAuthorization) GetUser(username, passwordHash string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", username, passwordHash)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockAuthorizationMockRecorder) GetUser(username, passwordHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockAuthorization)(nil).GetUser), username, passwordHash)
}
