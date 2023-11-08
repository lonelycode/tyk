// Code generated by MockGen. DO NOT EDIT.
// Source: otel_graphql_engine_detailed.go

// Package graphql is a generated GoMock package.
package graphql

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	plan "github.com/TykTechnologies/graphql-go-tools/pkg/engine/plan"
	resolve "github.com/TykTechnologies/graphql-go-tools/pkg/engine/resolve"
	graphql "github.com/TykTechnologies/graphql-go-tools/pkg/graphql"
	operationreport "github.com/TykTechnologies/graphql-go-tools/pkg/operationreport"
	postprocess "github.com/TykTechnologies/graphql-go-tools/pkg/postprocess"
)

// MockExecutionEngineI is a mock of ExecutionEngineI interface.
type MockExecutionEngineI struct {
	ctrl     *gomock.Controller
	recorder *MockExecutionEngineIMockRecorder
}

// MockExecutionEngineIMockRecorder is the mock recorder for MockExecutionEngineI.
type MockExecutionEngineIMockRecorder struct {
	mock *MockExecutionEngineI
}

// NewMockExecutionEngineI creates a new mock instance.
func NewMockExecutionEngineI(ctrl *gomock.Controller) *MockExecutionEngineI {
	mock := &MockExecutionEngineI{ctrl: ctrl}
	mock.recorder = &MockExecutionEngineIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExecutionEngineI) EXPECT() *MockExecutionEngineIMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockExecutionEngineI) Execute(ctx context.Context, operation *graphql.Request, writer resolve.FlushWriter, options ...graphql.ExecutionOptionsV2) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, operation, writer}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Execute", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute.
func (mr *MockExecutionEngineIMockRecorder) Execute(ctx, operation, writer interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, operation, writer}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockExecutionEngineI)(nil).Execute), varargs...)
}

// InputValidation mocks base method.
func (m *MockExecutionEngineI) InputValidation(operation *graphql.Request) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InputValidation", operation)
	ret0, _ := ret[0].(error)
	return ret0
}

// InputValidation indicates an expected call of InputValidation.
func (mr *MockExecutionEngineIMockRecorder) InputValidation(operation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InputValidation", reflect.TypeOf((*MockExecutionEngineI)(nil).InputValidation), operation)
}

// Normalize mocks base method.
func (m *MockExecutionEngineI) Normalize(operation *graphql.Request) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Normalize", operation)
	ret0, _ := ret[0].(error)
	return ret0
}

// Normalize indicates an expected call of Normalize.
func (mr *MockExecutionEngineIMockRecorder) Normalize(operation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Normalize", reflect.TypeOf((*MockExecutionEngineI)(nil).Normalize), operation)
}

// Plan mocks base method.
func (m *MockExecutionEngineI) Plan(postProcessor *postprocess.Processor, operation *graphql.Request, report *operationreport.Report) (plan.Plan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Plan", postProcessor, operation, report)
	ret0, _ := ret[0].(plan.Plan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Plan indicates an expected call of Plan.
func (mr *MockExecutionEngineIMockRecorder) Plan(postProcessor, operation, report interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Plan", reflect.TypeOf((*MockExecutionEngineI)(nil).Plan), postProcessor, operation, report)
}

// Resolve mocks base method.
func (m *MockExecutionEngineI) Resolve(resolveContext *resolve.Context, planResult plan.Plan, writer resolve.FlushWriter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resolve", resolveContext, planResult, writer)
	ret0, _ := ret[0].(error)
	return ret0
}

// Resolve indicates an expected call of Resolve.
func (mr *MockExecutionEngineIMockRecorder) Resolve(resolveContext, planResult, writer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resolve", reflect.TypeOf((*MockExecutionEngineI)(nil).Resolve), resolveContext, planResult, writer)
}

// Setup mocks base method.
func (m *MockExecutionEngineI) Setup(ctx context.Context, postProcessor *postprocess.Processor, resolveContext *resolve.Context, operation *graphql.Request, options ...graphql.ExecutionOptionsV2) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, postProcessor, resolveContext, operation}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Setup", varargs...)
}

// Setup indicates an expected call of Setup.
func (mr *MockExecutionEngineIMockRecorder) Setup(ctx, postProcessor, resolveContext, operation interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, postProcessor, resolveContext, operation}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockExecutionEngineI)(nil).Setup), varargs...)
}

// Teardown mocks base method.
func (m *MockExecutionEngineI) Teardown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Teardown")
}

// Teardown indicates an expected call of Teardown.
func (mr *MockExecutionEngineIMockRecorder) Teardown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Teardown", reflect.TypeOf((*MockExecutionEngineI)(nil).Teardown))
}

// ValidateForSchema mocks base method.
func (m *MockExecutionEngineI) ValidateForSchema(operation *graphql.Request) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateForSchema", operation)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateForSchema indicates an expected call of ValidateForSchema.
func (mr *MockExecutionEngineIMockRecorder) ValidateForSchema(operation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateForSchema", reflect.TypeOf((*MockExecutionEngineI)(nil).ValidateForSchema), operation)
}
