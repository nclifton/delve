//
// account.proto
//
// definition of the account rpc service

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: account/rpc/accountpb/account.proto

package accountpb

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type NoParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NoParams) Reset() {
	*x = NoParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_rpc_accountpb_account_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoParams) ProtoMessage() {}

func (x *NoParams) ProtoReflect() protoreflect.Message {
	mi := &file_account_rpc_accountpb_account_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoParams.ProtoReflect.Descriptor instead.
func (*NoParams) Descriptor() ([]byte, []int) {
	return file_account_rpc_accountpb_account_proto_rawDescGZIP(), []int{0}
}

type NoReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NoReply) Reset() {
	*x = NoReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_rpc_accountpb_account_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoReply) ProtoMessage() {}

func (x *NoReply) ProtoReflect() protoreflect.Message {
	mi := &file_account_rpc_accountpb_account_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoReply.ProtoReflect.Descriptor instead.
func (*NoReply) Descriptor() ([]byte, []int) {
	return file_account_rpc_accountpb_account_proto_rawDescGZIP(), []int{1}
}

type FindAccountByAPIKeyParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *FindAccountByAPIKeyParams) Reset() {
	*x = FindAccountByAPIKeyParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_rpc_accountpb_account_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindAccountByAPIKeyParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindAccountByAPIKeyParams) ProtoMessage() {}

func (x *FindAccountByAPIKeyParams) ProtoReflect() protoreflect.Message {
	mi := &file_account_rpc_accountpb_account_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindAccountByAPIKeyParams.ProtoReflect.Descriptor instead.
func (*FindAccountByAPIKeyParams) Descriptor() ([]byte, []int) {
	return file_account_rpc_accountpb_account_proto_rawDescGZIP(), []int{2}
}

func (x *FindAccountByAPIKeyParams) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type FindAccountByAPIKeyReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Account *Account `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
}

func (x *FindAccountByAPIKeyReply) Reset() {
	*x = FindAccountByAPIKeyReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_rpc_accountpb_account_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindAccountByAPIKeyReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindAccountByAPIKeyReply) ProtoMessage() {}

func (x *FindAccountByAPIKeyReply) ProtoReflect() protoreflect.Message {
	mi := &file_account_rpc_accountpb_account_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindAccountByAPIKeyReply.ProtoReflect.Descriptor instead.
func (*FindAccountByAPIKeyReply) Descriptor() ([]byte, []int) {
	return file_account_rpc_accountpb_account_proto_rawDescGZIP(), []int{3}
}

func (x *FindAccountByAPIKeyReply) GetAccount() *Account {
	if x != nil {
		return x.Account
	}
	return nil
}

type FindAccountByIDParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *FindAccountByIDParams) Reset() {
	*x = FindAccountByIDParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_rpc_accountpb_account_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindAccountByIDParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindAccountByIDParams) ProtoMessage() {}

func (x *FindAccountByIDParams) ProtoReflect() protoreflect.Message {
	mi := &file_account_rpc_accountpb_account_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindAccountByIDParams.ProtoReflect.Descriptor instead.
func (*FindAccountByIDParams) Descriptor() ([]byte, []int) {
	return file_account_rpc_accountpb_account_proto_rawDescGZIP(), []int{4}
}

func (x *FindAccountByIDParams) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type FindAccountByIDReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Account *Account `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
}

func (x *FindAccountByIDReply) Reset() {
	*x = FindAccountByIDReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_rpc_accountpb_account_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindAccountByIDReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindAccountByIDReply) ProtoMessage() {}

func (x *FindAccountByIDReply) ProtoReflect() protoreflect.Message {
	mi := &file_account_rpc_accountpb_account_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindAccountByIDReply.ProtoReflect.Descriptor instead.
func (*FindAccountByIDReply) Descriptor() ([]byte, []int) {
	return file_account_rpc_accountpb_account_proto_rawDescGZIP(), []int{5}
}

func (x *FindAccountByIDReply) GetAccount() *Account {
	if x != nil {
		return x.Account
	}
	return nil
}

type Account struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name           string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	AlarisUsername string                 `protobuf:"bytes,3,opt,name=alaris_username,json=alarisUsername,proto3" json:"alaris_username,omitempty"`
	AlarisPassword string                 `protobuf:"bytes,4,opt,name=alaris_password,json=alarisPassword,proto3" json:"alaris_password,omitempty"`
	AlarisUrl      string                 `protobuf:"bytes,5,opt,name=alaris_url,json=alarisUrl,proto3" json:"alaris_url,omitempty"`
	AccountApiKey  []*AccountAPIKey       `protobuf:"bytes,6,rep,name=account_api_key,json=accountApiKey,proto3" json:"account_api_key,omitempty"`
	CreatedAt      *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=Created_at,json=CreatedAt,proto3" json:"Created_at,omitempty"`
	UpdatedAt      *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=Updated_at,json=UpdatedAt,proto3" json:"Updated_at,omitempty"`
}

func (x *Account) Reset() {
	*x = Account{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_rpc_accountpb_account_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Account) ProtoMessage() {}

func (x *Account) ProtoReflect() protoreflect.Message {
	mi := &file_account_rpc_accountpb_account_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Account.ProtoReflect.Descriptor instead.
func (*Account) Descriptor() ([]byte, []int) {
	return file_account_rpc_accountpb_account_proto_rawDescGZIP(), []int{6}
}

func (x *Account) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Account) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Account) GetAlarisUsername() string {
	if x != nil {
		return x.AlarisUsername
	}
	return ""
}

func (x *Account) GetAlarisPassword() string {
	if x != nil {
		return x.AlarisPassword
	}
	return ""
}

func (x *Account) GetAlarisUrl() string {
	if x != nil {
		return x.AlarisUrl
	}
	return ""
}

func (x *Account) GetAccountApiKey() []*AccountAPIKey {
	if x != nil {
		return x.AccountApiKey
	}
	return nil
}

func (x *Account) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Account) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type AccountAPIKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description string                 `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	Key         string                 `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=Created_at,json=CreatedAt,proto3" json:"Created_at,omitempty"`
	UpdatedAt   *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=Updated_at,json=UpdatedAt,proto3" json:"Updated_at,omitempty"`
}

func (x *AccountAPIKey) Reset() {
	*x = AccountAPIKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_rpc_accountpb_account_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccountAPIKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountAPIKey) ProtoMessage() {}

func (x *AccountAPIKey) ProtoReflect() protoreflect.Message {
	mi := &file_account_rpc_accountpb_account_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountAPIKey.ProtoReflect.Descriptor instead.
func (*AccountAPIKey) Descriptor() ([]byte, []int) {
	return file_account_rpc_accountpb_account_proto_rawDescGZIP(), []int{7}
}

func (x *AccountAPIKey) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *AccountAPIKey) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *AccountAPIKey) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *AccountAPIKey) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_account_rpc_accountpb_account_proto protoreflect.FileDescriptor

var file_account_rpc_accountpb_account_proto_rawDesc = []byte{
	0x0a, 0x23, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x70, 0x62, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x70, 0x62,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x0a, 0x0a, 0x08, 0x4e, 0x6f, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x22, 0x09, 0x0a,
	0x07, 0x4e, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x2d, 0x0a, 0x19, 0x46, 0x69, 0x6e, 0x64,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x48, 0x0a, 0x18, 0x46, 0x69, 0x6e, 0x64, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x2c, 0x0a, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x70, 0x62,
	0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x22, 0x27, 0x0a, 0x15, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x42, 0x79, 0x49, 0x44, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x44, 0x0a, 0x14, 0x46, 0x69,
	0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x2c, 0x0a, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x70, 0x62, 0x2e,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x22, 0xd6, 0x02, 0x0a, 0x07, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x27, 0x0a, 0x0f, 0x61, 0x6c, 0x61, 0x72, 0x69, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x6c, 0x61, 0x72, 0x69,
	0x73, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x6c, 0x61,
	0x72, 0x69, 0x73, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x61, 0x6c, 0x61, 0x72, 0x69, 0x73, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x6c, 0x61, 0x72, 0x69, 0x73, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x6c, 0x61, 0x72, 0x69, 0x73, 0x55, 0x72,
	0x6c, 0x12, 0x40, 0x0a, 0x0f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x61, 0x70, 0x69,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x50,
	0x49, 0x4b, 0x65, 0x79, 0x52, 0x0d, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x70, 0x69,
	0x4b, 0x65, 0x79, 0x12, 0x39, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39,
	0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0xb9, 0x01, 0x0a, 0x0d, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x39, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x32, 0xc1, 0x01, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x60, 0x0a, 0x13, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x42, 0x79, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x12, 0x24, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x42, 0x79, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x23,
	0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x54, 0x0a, 0x0f, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x42, 0x79, 0x49, 0x44, 0x12, 0x20, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79,
	0x49, 0x44, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x1f, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x72, 0x73, 0x74, 0x73, 0x6d, 0x73,
	0x2f, 0x6d, 0x74, 0x6d, 0x6f, 0x2d, 0x74, 0x70, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64,
	0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_account_rpc_accountpb_account_proto_rawDescOnce sync.Once
	file_account_rpc_accountpb_account_proto_rawDescData = file_account_rpc_accountpb_account_proto_rawDesc
)

func file_account_rpc_accountpb_account_proto_rawDescGZIP() []byte {
	file_account_rpc_accountpb_account_proto_rawDescOnce.Do(func() {
		file_account_rpc_accountpb_account_proto_rawDescData = protoimpl.X.CompressGZIP(file_account_rpc_accountpb_account_proto_rawDescData)
	})
	return file_account_rpc_accountpb_account_proto_rawDescData
}

var file_account_rpc_accountpb_account_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_account_rpc_accountpb_account_proto_goTypes = []interface{}{
	(*NoParams)(nil),                  // 0: accountpb.NoParams
	(*NoReply)(nil),                   // 1: accountpb.NoReply
	(*FindAccountByAPIKeyParams)(nil), // 2: accountpb.FindAccountByAPIKeyParams
	(*FindAccountByAPIKeyReply)(nil),  // 3: accountpb.FindAccountByAPIKeyReply
	(*FindAccountByIDParams)(nil),     // 4: accountpb.FindAccountByIDParams
	(*FindAccountByIDReply)(nil),      // 5: accountpb.FindAccountByIDReply
	(*Account)(nil),                   // 6: accountpb.Account
	(*AccountAPIKey)(nil),             // 7: accountpb.AccountAPIKey
	(*timestamppb.Timestamp)(nil),     // 8: google.protobuf.Timestamp
}
var file_account_rpc_accountpb_account_proto_depIdxs = []int32{
	6, // 0: accountpb.FindAccountByAPIKeyReply.account:type_name -> accountpb.Account
	6, // 1: accountpb.FindAccountByIDReply.account:type_name -> accountpb.Account
	7, // 2: accountpb.Account.account_api_key:type_name -> accountpb.AccountAPIKey
	8, // 3: accountpb.Account.Created_at:type_name -> google.protobuf.Timestamp
	8, // 4: accountpb.Account.Updated_at:type_name -> google.protobuf.Timestamp
	8, // 5: accountpb.AccountAPIKey.Created_at:type_name -> google.protobuf.Timestamp
	8, // 6: accountpb.AccountAPIKey.Updated_at:type_name -> google.protobuf.Timestamp
	2, // 7: accountpb.Service.FindAccountByAPIKey:input_type -> accountpb.FindAccountByAPIKeyParams
	4, // 8: accountpb.Service.FindAccountByID:input_type -> accountpb.FindAccountByIDParams
	3, // 9: accountpb.Service.FindAccountByAPIKey:output_type -> accountpb.FindAccountByAPIKeyReply
	5, // 10: accountpb.Service.FindAccountByID:output_type -> accountpb.FindAccountByIDReply
	9, // [9:11] is the sub-list for method output_type
	7, // [7:9] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_account_rpc_accountpb_account_proto_init() }
func file_account_rpc_accountpb_account_proto_init() {
	if File_account_rpc_accountpb_account_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_account_rpc_accountpb_account_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NoParams); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_account_rpc_accountpb_account_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NoReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_account_rpc_accountpb_account_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindAccountByAPIKeyParams); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_account_rpc_accountpb_account_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindAccountByAPIKeyReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_account_rpc_accountpb_account_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindAccountByIDParams); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_account_rpc_accountpb_account_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindAccountByIDReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_account_rpc_accountpb_account_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Account); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_account_rpc_accountpb_account_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccountAPIKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_account_rpc_accountpb_account_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_account_rpc_accountpb_account_proto_goTypes,
		DependencyIndexes: file_account_rpc_accountpb_account_proto_depIdxs,
		MessageInfos:      file_account_rpc_accountpb_account_proto_msgTypes,
	}.Build()
	File_account_rpc_accountpb_account_proto = out.File
	file_account_rpc_accountpb_account_proto_rawDesc = nil
	file_account_rpc_accountpb_account_proto_goTypes = nil
	file_account_rpc_accountpb_account_proto_depIdxs = nil
}
