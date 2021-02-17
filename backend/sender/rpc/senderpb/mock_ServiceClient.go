// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package senderpb

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// MockServiceClient is an autogenerated mock type for the ServiceClient type
type MockServiceClient struct {
	mock.Mock
}

// CreateSendersFromCSVDataURL provides a mock function with given fields: ctx, in, opts
func (_m *MockServiceClient) CreateSendersFromCSVDataURL(ctx context.Context, in *CreateSendersFromCSVDataURLParams, opts ...grpc.CallOption) (*CreateSendersFromCSVDataURLReply, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *CreateSendersFromCSVDataURLReply
	if rf, ok := ret.Get(0).(func(context.Context, *CreateSendersFromCSVDataURLParams, ...grpc.CallOption) *CreateSendersFromCSVDataURLReply); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*CreateSendersFromCSVDataURLReply)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *CreateSendersFromCSVDataURLParams, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindSenderByAddressAndAccountID provides a mock function with given fields: ctx, in, opts
func (_m *MockServiceClient) FindSenderByAddressAndAccountID(ctx context.Context, in *FindSenderByAddressAndAccountIDParams, opts ...grpc.CallOption) (*FindSenderByAddressAndAccountIDReply, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *FindSenderByAddressAndAccountIDReply
	if rf, ok := ret.Get(0).(func(context.Context, *FindSenderByAddressAndAccountIDParams, ...grpc.CallOption) *FindSenderByAddressAndAccountIDReply); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*FindSenderByAddressAndAccountIDReply)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *FindSenderByAddressAndAccountIDParams, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindSendersByAccountId provides a mock function with given fields: ctx, in, opts
func (_m *MockServiceClient) FindSendersByAccountId(ctx context.Context, in *FindSendersByAccountIdParams, opts ...grpc.CallOption) (*FindSendersByAccountIdReply, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *FindSendersByAccountIdReply
	if rf, ok := ret.Get(0).(func(context.Context, *FindSendersByAccountIdParams, ...grpc.CallOption) *FindSendersByAccountIdReply); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*FindSendersByAccountIdReply)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *FindSendersByAccountIdParams, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindSendersByAddress provides a mock function with given fields: ctx, in, opts
func (_m *MockServiceClient) FindSendersByAddress(ctx context.Context, in *FindSendersByAddressParams, opts ...grpc.CallOption) (*FindSendersByAddressReply, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *FindSendersByAddressReply
	if rf, ok := ret.Get(0).(func(context.Context, *FindSendersByAddressParams, ...grpc.CallOption) *FindSendersByAddressReply); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*FindSendersByAddressReply)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *FindSendersByAddressParams, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
