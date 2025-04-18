// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: internal/grpc/ipcheck.proto

package ipcheckpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CheckRequest struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	Ip               string                 `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	AllowedCountries []string               `protobuf:"bytes,2,rep,name=allowed_countries,json=allowedCountries,proto3" json:"allowed_countries,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *CheckRequest) Reset() {
	*x = CheckRequest{}
	mi := &file_internal_grpc_ipcheck_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRequest) ProtoMessage() {}

func (x *CheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpc_ipcheck_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRequest.ProtoReflect.Descriptor instead.
func (*CheckRequest) Descriptor() ([]byte, []int) {
	return file_internal_grpc_ipcheck_proto_rawDescGZIP(), []int{0}
}

func (x *CheckRequest) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *CheckRequest) GetAllowedCountries() []string {
	if x != nil {
		return x.AllowedCountries
	}
	return nil
}

type CheckResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Allowed       bool                   `protobuf:"varint,1,opt,name=allowed,proto3" json:"allowed,omitempty"`
	Country       string                 `protobuf:"bytes,2,opt,name=country,proto3" json:"country,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckResponse) Reset() {
	*x = CheckResponse{}
	mi := &file_internal_grpc_ipcheck_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse) ProtoMessage() {}

func (x *CheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpc_ipcheck_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckResponse.ProtoReflect.Descriptor instead.
func (*CheckResponse) Descriptor() ([]byte, []int) {
	return file_internal_grpc_ipcheck_proto_rawDescGZIP(), []int{1}
}

func (x *CheckResponse) GetAllowed() bool {
	if x != nil {
		return x.Allowed
	}
	return false
}

func (x *CheckResponse) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

var File_internal_grpc_ipcheck_proto protoreflect.FileDescriptor

const file_internal_grpc_ipcheck_proto_rawDesc = "" +
	"\n" +
	"\x1binternal/grpc/ipcheck.proto\x12\tipcheckpb\"K\n" +
	"\fCheckRequest\x12\x0e\n" +
	"\x02ip\x18\x01 \x01(\tR\x02ip\x12+\n" +
	"\x11allowed_countries\x18\x02 \x03(\tR\x10allowedCountries\"C\n" +
	"\rCheckResponse\x12\x18\n" +
	"\aallowed\x18\x01 \x01(\bR\aallowed\x12\x18\n" +
	"\acountry\x18\x02 \x01(\tR\acountry2I\n" +
	"\tIpChecker\x12<\n" +
	"\aCheckIP\x12\x17.ipcheckpb.CheckRequest\x1a\x18.ipcheckpb.CheckResponseB\x19Z\x17internal/grpc/ipcheckpbb\x06proto3"

var (
	file_internal_grpc_ipcheck_proto_rawDescOnce sync.Once
	file_internal_grpc_ipcheck_proto_rawDescData []byte
)

func file_internal_grpc_ipcheck_proto_rawDescGZIP() []byte {
	file_internal_grpc_ipcheck_proto_rawDescOnce.Do(func() {
		file_internal_grpc_ipcheck_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_internal_grpc_ipcheck_proto_rawDesc), len(file_internal_grpc_ipcheck_proto_rawDesc)))
	})
	return file_internal_grpc_ipcheck_proto_rawDescData
}

var file_internal_grpc_ipcheck_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_internal_grpc_ipcheck_proto_goTypes = []any{
	(*CheckRequest)(nil),  // 0: ipcheckpb.CheckRequest
	(*CheckResponse)(nil), // 1: ipcheckpb.CheckResponse
}
var file_internal_grpc_ipcheck_proto_depIdxs = []int32{
	0, // 0: ipcheckpb.IpChecker.CheckIP:input_type -> ipcheckpb.CheckRequest
	1, // 1: ipcheckpb.IpChecker.CheckIP:output_type -> ipcheckpb.CheckResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_grpc_ipcheck_proto_init() }
func file_internal_grpc_ipcheck_proto_init() {
	if File_internal_grpc_ipcheck_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_internal_grpc_ipcheck_proto_rawDesc), len(file_internal_grpc_ipcheck_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_grpc_ipcheck_proto_goTypes,
		DependencyIndexes: file_internal_grpc_ipcheck_proto_depIdxs,
		MessageInfos:      file_internal_grpc_ipcheck_proto_msgTypes,
	}.Build()
	File_internal_grpc_ipcheck_proto = out.File
	file_internal_grpc_ipcheck_proto_goTypes = nil
	file_internal_grpc_ipcheck_proto_depIdxs = nil
}
