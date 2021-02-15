//
// sender.proto
//
// definition of the sender rpc service

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: sender/rpc/senderpb/sender.proto

package senderpb

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
		mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoParams) ProtoMessage() {}

func (x *NoParams) ProtoReflect() protoreflect.Message {
	mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[0]
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
	return file_sender_rpc_senderpb_sender_proto_rawDescGZIP(), []int{0}
}

type NoReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NoReply) Reset() {
	*x = NoReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoReply) ProtoMessage() {}

func (x *NoReply) ProtoReflect() protoreflect.Message {
	mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[1]
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
	return file_sender_rpc_senderpb_sender_proto_rawDescGZIP(), []int{1}
}

type Sender struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                 `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	AccountId      string                 `protobuf:"bytes,2,opt,name=Account_id,json=AccountId,proto3" json:"Account_id,omitempty"`
	Address        string                 `protobuf:"bytes,3,opt,name=Address,proto3" json:"Address,omitempty"`
	MMSProviderKey string                 `protobuf:"bytes,4,opt,name=MMS_provider_key,json=MMSProviderKey,proto3" json:"MMS_provider_key,omitempty"`
	Channels       []string               `protobuf:"bytes,5,rep,name=Channels,proto3" json:"Channels,omitempty"`
	Country        string                 `protobuf:"bytes,6,opt,name=Country,proto3" json:"Country,omitempty"`
	Comment        string                 `protobuf:"bytes,7,opt,name=Comment,proto3" json:"Comment,omitempty"`
	CreatedAt      *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=Created_at,json=CreatedAt,proto3" json:"Created_at,omitempty"`
	UpdatedAt      *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=Updated_at,json=UpdatedAt,proto3" json:"Updated_at,omitempty"`
}

func (x *Sender) Reset() {
	*x = Sender{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sender) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sender) ProtoMessage() {}

func (x *Sender) ProtoReflect() protoreflect.Message {
	mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sender.ProtoReflect.Descriptor instead.
func (*Sender) Descriptor() ([]byte, []int) {
	return file_sender_rpc_senderpb_sender_proto_rawDescGZIP(), []int{2}
}

func (x *Sender) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Sender) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *Sender) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Sender) GetMMSProviderKey() string {
	if x != nil {
		return x.MMSProviderKey
	}
	return ""
}

func (x *Sender) GetChannels() []string {
	if x != nil {
		return x.Channels
	}
	return nil
}

func (x *Sender) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

func (x *Sender) GetComment() string {
	if x != nil {
		return x.Comment
	}
	return ""
}

func (x *Sender) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Sender) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type NewSender struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId      string   `protobuf:"bytes,1,opt,name=Account_id,json=AccountId,proto3" json:"Account_id,omitempty"`
	Address        string   `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
	MMSProviderKey string   `protobuf:"bytes,3,opt,name=MMS_provider_key,json=MMSProviderKey,proto3" json:"MMS_provider_key,omitempty"`
	Channels       []string `protobuf:"bytes,4,rep,name=Channels,proto3" json:"Channels,omitempty"`
	Country        string   `protobuf:"bytes,5,opt,name=Country,proto3" json:"Country,omitempty"`
	Comment        string   `protobuf:"bytes,6,opt,name=Comment,proto3" json:"Comment,omitempty"`
}

func (x *NewSender) Reset() {
	*x = NewSender{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewSender) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewSender) ProtoMessage() {}

func (x *NewSender) ProtoReflect() protoreflect.Message {
	mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewSender.ProtoReflect.Descriptor instead.
func (*NewSender) Descriptor() ([]byte, []int) {
	return file_sender_rpc_senderpb_sender_proto_rawDescGZIP(), []int{3}
}

func (x *NewSender) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *NewSender) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *NewSender) GetMMSProviderKey() string {
	if x != nil {
		return x.MMSProviderKey
	}
	return ""
}

func (x *NewSender) GetChannels() []string {
	if x != nil {
		return x.Channels
	}
	return nil
}

func (x *NewSender) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

func (x *NewSender) GetComment() string {
	if x != nil {
		return x.Comment
	}
	return ""
}

type FindSenderByAddressAndAccountIDParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId string `protobuf:"bytes,1,opt,name=Account_id,json=AccountId,proto3" json:"Account_id,omitempty"`
	Address   string `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
}

func (x *FindSenderByAddressAndAccountIDParams) Reset() {
	*x = FindSenderByAddressAndAccountIDParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindSenderByAddressAndAccountIDParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindSenderByAddressAndAccountIDParams) ProtoMessage() {}

func (x *FindSenderByAddressAndAccountIDParams) ProtoReflect() protoreflect.Message {
	mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindSenderByAddressAndAccountIDParams.ProtoReflect.Descriptor instead.
func (*FindSenderByAddressAndAccountIDParams) Descriptor() ([]byte, []int) {
	return file_sender_rpc_senderpb_sender_proto_rawDescGZIP(), []int{4}
}

func (x *FindSenderByAddressAndAccountIDParams) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *FindSenderByAddressAndAccountIDParams) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type FindSenderByAddressAndAccountIDReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender *Sender `protobuf:"bytes,1,opt,name=Sender,proto3" json:"Sender,omitempty"`
}

func (x *FindSenderByAddressAndAccountIDReply) Reset() {
	*x = FindSenderByAddressAndAccountIDReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindSenderByAddressAndAccountIDReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindSenderByAddressAndAccountIDReply) ProtoMessage() {}

func (x *FindSenderByAddressAndAccountIDReply) ProtoReflect() protoreflect.Message {
	mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindSenderByAddressAndAccountIDReply.ProtoReflect.Descriptor instead.
func (*FindSenderByAddressAndAccountIDReply) Descriptor() ([]byte, []int) {
	return file_sender_rpc_senderpb_sender_proto_rawDescGZIP(), []int{5}
}

func (x *FindSenderByAddressAndAccountIDReply) GetSender() *Sender {
	if x != nil {
		return x.Sender
	}
	return nil
}

type FindSendersByAccountIdParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId string `protobuf:"bytes,1,opt,name=Account_id,json=AccountId,proto3" json:"Account_id,omitempty"`
}

func (x *FindSendersByAccountIdParams) Reset() {
	*x = FindSendersByAccountIdParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindSendersByAccountIdParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindSendersByAccountIdParams) ProtoMessage() {}

func (x *FindSendersByAccountIdParams) ProtoReflect() protoreflect.Message {
	mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindSendersByAccountIdParams.ProtoReflect.Descriptor instead.
func (*FindSendersByAccountIdParams) Descriptor() ([]byte, []int) {
	return file_sender_rpc_senderpb_sender_proto_rawDescGZIP(), []int{6}
}

func (x *FindSendersByAccountIdParams) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

type FindSendersByAccountIdReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Senders []*Sender `protobuf:"bytes,1,rep,name=Senders,proto3" json:"Senders,omitempty"`
}

func (x *FindSendersByAccountIdReply) Reset() {
	*x = FindSendersByAccountIdReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindSendersByAccountIdReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindSendersByAccountIdReply) ProtoMessage() {}

func (x *FindSendersByAccountIdReply) ProtoReflect() protoreflect.Message {
	mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindSendersByAccountIdReply.ProtoReflect.Descriptor instead.
func (*FindSendersByAccountIdReply) Descriptor() ([]byte, []int) {
	return file_sender_rpc_senderpb_sender_proto_rawDescGZIP(), []int{7}
}

func (x *FindSendersByAccountIdReply) GetSenders() []*Sender {
	if x != nil {
		return x.Senders
	}
	return nil
}

type FindSendersByAddressParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
}

func (x *FindSendersByAddressParams) Reset() {
	*x = FindSendersByAddressParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindSendersByAddressParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindSendersByAddressParams) ProtoMessage() {}

func (x *FindSendersByAddressParams) ProtoReflect() protoreflect.Message {
	mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindSendersByAddressParams.ProtoReflect.Descriptor instead.
func (*FindSendersByAddressParams) Descriptor() ([]byte, []int) {
	return file_sender_rpc_senderpb_sender_proto_rawDescGZIP(), []int{8}
}

func (x *FindSendersByAddressParams) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type FindSendersByAddressReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Senders []*Sender `protobuf:"bytes,1,rep,name=Senders,proto3" json:"Senders,omitempty"`
}

func (x *FindSendersByAddressReply) Reset() {
	*x = FindSendersByAddressReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindSendersByAddressReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindSendersByAddressReply) ProtoMessage() {}

func (x *FindSendersByAddressReply) ProtoReflect() protoreflect.Message {
	mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindSendersByAddressReply.ProtoReflect.Descriptor instead.
func (*FindSendersByAddressReply) Descriptor() ([]byte, []int) {
	return file_sender_rpc_senderpb_sender_proto_rawDescGZIP(), []int{9}
}

func (x *FindSendersByAddressReply) GetSenders() []*Sender {
	if x != nil {
		return x.Senders
	}
	return nil
}

type CreateSendersParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Senders []*NewSender `protobuf:"bytes,1,rep,name=Senders,proto3" json:"Senders,omitempty"`
}

func (x *CreateSendersParams) Reset() {
	*x = CreateSendersParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSendersParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSendersParams) ProtoMessage() {}

func (x *CreateSendersParams) ProtoReflect() protoreflect.Message {
	mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSendersParams.ProtoReflect.Descriptor instead.
func (*CreateSendersParams) Descriptor() ([]byte, []int) {
	return file_sender_rpc_senderpb_sender_proto_rawDescGZIP(), []int{10}
}

func (x *CreateSendersParams) GetSenders() []*NewSender {
	if x != nil {
		return x.Senders
	}
	return nil
}

type CreateSendersReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Senders []*Sender `protobuf:"bytes,1,rep,name=Senders,proto3" json:"Senders,omitempty"`
}

func (x *CreateSendersReply) Reset() {
	*x = CreateSendersReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSendersReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSendersReply) ProtoMessage() {}

func (x *CreateSendersReply) ProtoReflect() protoreflect.Message {
	mi := &file_sender_rpc_senderpb_sender_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSendersReply.ProtoReflect.Descriptor instead.
func (*CreateSendersReply) Descriptor() ([]byte, []int) {
	return file_sender_rpc_senderpb_sender_proto_rawDescGZIP(), []int{11}
}

func (x *CreateSendersReply) GetSenders() []*Sender {
	if x != nil {
		return x.Senders
	}
	return nil
}

var File_sender_rpc_senderpb_sender_proto protoreflect.FileDescriptor

var file_sender_rpc_senderpb_sender_proto_rawDesc = []byte{
	0x0a, 0x20, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x70, 0x62, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0a, 0x0a,
	0x08, 0x4e, 0x6f, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x22, 0x09, 0x0a, 0x07, 0x4e, 0x6f, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0xc1, 0x02, 0x0a, 0x06, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12,
	0x1d, 0x0a, 0x0a, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x28, 0x0a, 0x10, 0x4d, 0x4d, 0x53, 0x5f,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x4d, 0x4d, 0x53, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x4b,
	0x65, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x12, 0x18,
	0x0a, 0x07, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a,
	0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0xbe, 0x01, 0x0a, 0x09, 0x4e, 0x65, 0x77,
	0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12,
	0x28, 0x0a, 0x10, 0x4d, 0x4d, 0x53, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f,
	0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x4d, 0x4d, 0x53, 0x50, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x60, 0x0a, 0x25, 0x46, 0x69, 0x6e,
	0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x42, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x41, 0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x50, 0x0a, 0x24, 0x46,
	0x69, 0x6e, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x42, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x41, 0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x28, 0x0a, 0x06, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x53,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x06, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x22, 0x3d, 0x0a,
	0x1c, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x42, 0x79, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x1d, 0x0a,
	0x0a, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x49, 0x0a, 0x1b,
	0x46, 0x69, 0x6e, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x42, 0x79, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x2a, 0x0a, 0x07, 0x53,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x07,
	0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x22, 0x36, 0x0a, 0x1a, 0x46, 0x69, 0x6e, 0x64, 0x53,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x42, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22,
	0x47, 0x0a, 0x19, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x42, 0x79,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x2a, 0x0a, 0x07,
	0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x52,
	0x07, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x22, 0x44, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12,
	0x2d, 0x0a, 0x07, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x4e, 0x65, 0x77, 0x53,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x07, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x22, 0x40,
	0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x2a, 0x0a, 0x07, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62,
	0x2e, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x07, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73,
	0x32, 0xa8, 0x03, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x82, 0x01, 0x0a,
	0x1f, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x42, 0x79, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x41, 0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44,
	0x12, 0x2f, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x64,
	0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x42, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x41,
	0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x1a, 0x2e, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e,
	0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x42, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x41, 0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x67, 0x0a, 0x16, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73,
	0x42, 0x79, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x26, 0x2e, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x73, 0x42, 0x79, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x1a, 0x25, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x46,
	0x69, 0x6e, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x42, 0x79, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x49, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x61, 0x0a, 0x14, 0x46, 0x69,
	0x6e, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x42, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x24, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x46, 0x69,
	0x6e, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x42, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x23, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x42,
	0x79, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4c, 0x0a,
	0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x12, 0x1d,
	0x2e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x1c, 0x2e,
	0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x39, 0x5a, 0x37, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x72, 0x73, 0x74, 0x73,
	0x6d, 0x73, 0x2f, 0x6d, 0x74, 0x6d, 0x6f, 0x2d, 0x74, 0x70, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65,
	0x6e, 0x64, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sender_rpc_senderpb_sender_proto_rawDescOnce sync.Once
	file_sender_rpc_senderpb_sender_proto_rawDescData = file_sender_rpc_senderpb_sender_proto_rawDesc
)

func file_sender_rpc_senderpb_sender_proto_rawDescGZIP() []byte {
	file_sender_rpc_senderpb_sender_proto_rawDescOnce.Do(func() {
		file_sender_rpc_senderpb_sender_proto_rawDescData = protoimpl.X.CompressGZIP(file_sender_rpc_senderpb_sender_proto_rawDescData)
	})
	return file_sender_rpc_senderpb_sender_proto_rawDescData
}

var file_sender_rpc_senderpb_sender_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_sender_rpc_senderpb_sender_proto_goTypes = []interface{}{
	(*NoParams)(nil),  // 0: senderpb.NoParams
	(*NoReply)(nil),   // 1: senderpb.NoReply
	(*Sender)(nil),    // 2: senderpb.Sender
	(*NewSender)(nil), // 3: senderpb.NewSender
	(*FindSenderByAddressAndAccountIDParams)(nil), // 4: senderpb.FindSenderByAddressAndAccountIDParams
	(*FindSenderByAddressAndAccountIDReply)(nil),  // 5: senderpb.FindSenderByAddressAndAccountIDReply
	(*FindSendersByAccountIdParams)(nil),          // 6: senderpb.FindSendersByAccountIdParams
	(*FindSendersByAccountIdReply)(nil),           // 7: senderpb.FindSendersByAccountIdReply
	(*FindSendersByAddressParams)(nil),            // 8: senderpb.FindSendersByAddressParams
	(*FindSendersByAddressReply)(nil),             // 9: senderpb.FindSendersByAddressReply
	(*CreateSendersParams)(nil),                   // 10: senderpb.CreateSendersParams
	(*CreateSendersReply)(nil),                    // 11: senderpb.CreateSendersReply
	(*timestamppb.Timestamp)(nil),                 // 12: google.protobuf.Timestamp
}
var file_sender_rpc_senderpb_sender_proto_depIdxs = []int32{
	12, // 0: senderpb.Sender.Created_at:type_name -> google.protobuf.Timestamp
	12, // 1: senderpb.Sender.Updated_at:type_name -> google.protobuf.Timestamp
	2,  // 2: senderpb.FindSenderByAddressAndAccountIDReply.Sender:type_name -> senderpb.Sender
	2,  // 3: senderpb.FindSendersByAccountIdReply.Senders:type_name -> senderpb.Sender
	2,  // 4: senderpb.FindSendersByAddressReply.Senders:type_name -> senderpb.Sender
	3,  // 5: senderpb.CreateSendersParams.Senders:type_name -> senderpb.NewSender
	2,  // 6: senderpb.CreateSendersReply.Senders:type_name -> senderpb.Sender
	4,  // 7: senderpb.Service.FindSenderByAddressAndAccountID:input_type -> senderpb.FindSenderByAddressAndAccountIDParams
	6,  // 8: senderpb.Service.FindSendersByAccountId:input_type -> senderpb.FindSendersByAccountIdParams
	8,  // 9: senderpb.Service.FindSendersByAddress:input_type -> senderpb.FindSendersByAddressParams
	10, // 10: senderpb.Service.CreateSenders:input_type -> senderpb.CreateSendersParams
	5,  // 11: senderpb.Service.FindSenderByAddressAndAccountID:output_type -> senderpb.FindSenderByAddressAndAccountIDReply
	7,  // 12: senderpb.Service.FindSendersByAccountId:output_type -> senderpb.FindSendersByAccountIdReply
	9,  // 13: senderpb.Service.FindSendersByAddress:output_type -> senderpb.FindSendersByAddressReply
	11, // 14: senderpb.Service.CreateSenders:output_type -> senderpb.CreateSendersReply
	11, // [11:15] is the sub-list for method output_type
	7,  // [7:11] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_sender_rpc_senderpb_sender_proto_init() }
func file_sender_rpc_senderpb_sender_proto_init() {
	if File_sender_rpc_senderpb_sender_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sender_rpc_senderpb_sender_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_sender_rpc_senderpb_sender_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_sender_rpc_senderpb_sender_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sender); i {
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
		file_sender_rpc_senderpb_sender_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewSender); i {
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
		file_sender_rpc_senderpb_sender_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindSenderByAddressAndAccountIDParams); i {
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
		file_sender_rpc_senderpb_sender_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindSenderByAddressAndAccountIDReply); i {
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
		file_sender_rpc_senderpb_sender_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindSendersByAccountIdParams); i {
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
		file_sender_rpc_senderpb_sender_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindSendersByAccountIdReply); i {
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
		file_sender_rpc_senderpb_sender_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindSendersByAddressParams); i {
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
		file_sender_rpc_senderpb_sender_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindSendersByAddressReply); i {
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
		file_sender_rpc_senderpb_sender_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSendersParams); i {
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
		file_sender_rpc_senderpb_sender_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSendersReply); i {
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
			RawDescriptor: file_sender_rpc_senderpb_sender_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sender_rpc_senderpb_sender_proto_goTypes,
		DependencyIndexes: file_sender_rpc_senderpb_sender_proto_depIdxs,
		MessageInfos:      file_sender_rpc_senderpb_sender_proto_msgTypes,
	}.Build()
	File_sender_rpc_senderpb_sender_proto = out.File
	file_sender_rpc_senderpb_sender_proto_rawDesc = nil
	file_sender_rpc_senderpb_sender_proto_goTypes = nil
	file_sender_rpc_senderpb_sender_proto_depIdxs = nil
}
