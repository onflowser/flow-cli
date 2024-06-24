// Code generated by mockery v2.43.2. DO NOT EDIT.

package migrate

import (
	common "github.com/onflow/cadence/runtime/common"
	mock "github.com/stretchr/testify/mock"
)

// mockStagingValidator is an autogenerated mock type for the stagingValidator type
type mockStagingValidator struct {
	mock.Mock
}

// PrettyPrintError provides a mock function with given fields: err, location
func (_m *mockStagingValidator) PrettyPrintError(err error, location common.Location) string {
	ret := _m.Called(err, location)

	if len(ret) == 0 {
		panic("no return value specified for PrettyPrintError")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(error, common.Location) string); ok {
		r0 = rf(err, location)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Validate provides a mock function with given fields: stagedContracts
func (_m *mockStagingValidator) Validate(stagedContracts []stagedContractUpdate) error {
	ret := _m.Called(stagedContracts)

	if len(ret) == 0 {
		panic("no return value specified for Validate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]stagedContractUpdate) error); ok {
		r0 = rf(stagedContracts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// newMockStagingValidator creates a new instance of mockStagingValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockStagingValidator(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockStagingValidator {
	mock := &mockStagingValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}