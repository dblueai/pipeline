// Code generated by mockery v1.0.0. DO NOT EDIT.

package auth

import context "context"
import mock "github.com/stretchr/testify/mock"

// MockOrganizationStore is an autogenerated mock type for the OrganizationStore type
type MockOrganizationStore struct {
	mock.Mock
}

// ApplyUserMembership provides a mock function with given fields: ctx, organizationID, userID, role
func (_m *MockOrganizationStore) ApplyUserMembership(ctx context.Context, organizationID uint, userID uint, role string) error {
	ret := _m.Called(ctx, organizationID, userID, role)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint, string) error); ok {
		r0 = rf(ctx, organizationID, userID, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EnsureOrganizationExists provides a mock function with given fields: ctx, name, provider
func (_m *MockOrganizationStore) EnsureOrganizationExists(ctx context.Context, name string, provider string) (bool, uint, error) {
	ret := _m.Called(ctx, name, provider)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, name, provider)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 uint
	if rf, ok := ret.Get(1).(func(context.Context, string, string) uint); ok {
		r1 = rf(ctx, name, provider)
	} else {
		r1 = ret.Get(1).(uint)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, name, provider)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetOrganizationMembershipsOf provides a mock function with given fields: ctx, userID
func (_m *MockOrganizationStore) GetOrganizationMembershipsOf(ctx context.Context, userID uint) ([]UserOrganization, error) {
	ret := _m.Called(ctx, userID)

	var r0 []UserOrganization
	if rf, ok := ret.Get(0).(func(context.Context, uint) []UserOrganization); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]UserOrganization)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveFromOrganization provides a mock function with given fields: ctx, organizationID, userID
func (_m *MockOrganizationStore) RemoveFromOrganization(ctx context.Context, organizationID uint, userID uint) error {
	ret := _m.Called(ctx, organizationID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) error); ok {
		r0 = rf(ctx, organizationID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
