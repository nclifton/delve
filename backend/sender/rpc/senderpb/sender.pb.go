//
// sender.proto
//
// definition of the sender rpc service

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: rpc/senderpb/sender.proto

package senderpb

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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
		mi := &file_rpc_senderpb_sender_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoParams) ProtoMessage() {}

func (x *NoParams) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_senderpb_sender_proto_msgTypes[0]
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
	return file_rpc_senderpb_sender_proto_rawDescGZIP(), []int{0}
}

type NoReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NoReply) Reset() {
	*x = NoReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_senderpb_sender_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoReply) ProtoMessage() {}

func (x *NoReply) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_senderpb_sender_proto_msgTypes[1]
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
	return file_rpc_senderpb_sender_proto_rawDescGZIP(), []int{1}
}

type Sender struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string               `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	AccountId      string               `protobuf:"bytes,2,opt,name=Account_id,json=AccountId,proto3" json:"Account_id,omitempty"`
	Address        string               `protobuf:"bytes,3,opt,name=Address,proto3" json:"Address,omitempty"`
	MMSProviderKey string               `protobuf:"bytes,4,opt,name=MMS_provider_key,json=MMSProviderKey,proto3" json:"MMS_provider_key,omitempty"`
	Channels       []string             `protobuf:"bytes,5,rep,name=Channels,proto3" json:"Channels,omitempty"`
	Country        string               `protobuf:"bytes,6,opt,name=Country,proto3" json:"Country,omitempty"`
	Comment        string               `protobuf:"bytes,7,opt,name=Comment,proto3" json:"Comment,omitempty"`
	CreatedAt      *timestamp.Timestamp `protobuf:"bytes,8,opt,name=Created_at,json=CreatedAt,proto3" json:"Created_at,omitempty"`
	UpdatedAt      *timestamp.Timestamp `protobuf:"bytes,9,opt,name=Updated_at,json=UpdatedAt,proto3" json:"Updated_at,omitempty"`
}

func (x *Sender) Reset() {
	*x = Sender{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_senderpb_sender_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sender) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sender) ProtoMessage() {}

func (x *Sender) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_senderpb_sender_proto_msgTypes[2]
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
	return file_rpc_senderpb_sender_proto_rawDescGZIP(), []int{2}
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

func (x *Sender) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Sender) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_rpc_senderpb_sender_proto protoreflect.FileDescriptor

var file_rpc_senderpb_sender_proto_rawDesc = []byte{
	0x0a, 0x19, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2f, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0a, 0x0a, 0x08, 0x4e, 0x6f, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x73, 0x22, 0x09, 0x0a, 0x07, 0x4e, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0xc1, 0x02,
	0x0a, 0x06, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x12, 0x28, 0x0a, 0x10, 0x4d, 0x4d, 0x53, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x4d, 0x4d, 0x53,
	0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x43,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x43,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x72, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x32, 0x09, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x42, 0x39, 0x5a, 0x37,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x72, 0x73, 0x74,
	0x73, 0x6d, 0x73, 0x2f, 0x6d, 0x74, 0x6d, 0x6f, 0x2d, 0x74, 0x70, 0x2f, 0x62, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_senderpb_sender_proto_rawDescOnce sync.Once
	file_rpc_senderpb_sender_proto_rawDescData = file_rpc_senderpb_sender_proto_rawDesc
)

func file_rpc_senderpb_sender_proto_rawDescGZIP() []byte {
	file_rpc_senderpb_sender_proto_rawDescOnce.Do(func() {
		file_rpc_senderpb_sender_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_senderpb_sender_proto_rawDescData)
	})
	return file_rpc_senderpb_sender_proto_rawDescData
}

var file_rpc_senderpb_sender_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_rpc_senderpb_sender_proto_goTypes = []interface{}{
	(*NoParams)(nil),            // 0: senderpb.NoParams
	(*NoReply)(nil),             // 1: senderpb.NoReply
	(*Sender)(nil),              // 2: senderpb.Sender
	(*timestamp.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_rpc_senderpb_sender_proto_depIdxs = []int32{
	3, // 0: senderpb.Sender.Created_at:type_name -> google.protobuf.Timestamp
	3, // 1: senderpb.Sender.Updated_at:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_rpc_senderpb_sender_proto_init() }
func file_rpc_senderpb_sender_proto_init() {
	if File_rpc_senderpb_sender_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_senderpb_sender_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_rpc_senderpb_sender_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_rpc_senderpb_sender_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_senderpb_sender_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_senderpb_sender_proto_goTypes,
		DependencyIndexes: file_rpc_senderpb_sender_proto_depIdxs,
		MessageInfos:      file_rpc_senderpb_sender_proto_msgTypes,
	}.Build()
	File_rpc_senderpb_sender_proto = out.File
	file_rpc_senderpb_sender_proto_rawDesc = nil
	file_rpc_senderpb_sender_proto_goTypes = nil
	file_rpc_senderpb_sender_proto_depIdxs = nil
}
