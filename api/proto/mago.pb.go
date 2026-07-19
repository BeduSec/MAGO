// Copyright (c) BeduSec. All rights reserved.
package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const _ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
const _ = protoimpl.MinVersion

type HealthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HealthRequest) Reset() { *x = HealthRequest{} }
func (x *HealthRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*HealthRequest) ProtoMessage() {}
func (x *HealthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mago_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*HealthRequest) Descriptor() ([]byte, []int) { return file_mago_proto_rawDescGZIP(), []int{0} }

type HealthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	Status        string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *HealthResponse) Reset() { *x = HealthResponse{} }
func (x *HealthResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*HealthResponse) ProtoMessage() {}
func (x *HealthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mago_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*HealthResponse) Descriptor() ([]byte, []int) { return file_mago_proto_rawDescGZIP(), []int{1} }
func (x *HealthResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type ReloadRulesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	Token         string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *ReloadRulesRequest) Reset() { *x = ReloadRulesRequest{} }
func (x *ReloadRulesRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*ReloadRulesRequest) ProtoMessage() {}
func (x *ReloadRulesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mago_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*ReloadRulesRequest) Descriptor() ([]byte, []int) { return file_mago_proto_rawDescGZIP(), []int{2} }
func (x *ReloadRulesRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type ReloadRulesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	Message       string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ReloadRulesResponse) Reset() { *x = ReloadRulesResponse{} }
func (x *ReloadRulesResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*ReloadRulesResponse) ProtoMessage() {}
func (x *ReloadRulesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mago_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*ReloadRulesResponse) Descriptor() ([]byte, []int) { return file_mago_proto_rawDescGZIP(), []int{3} }
func (x *ReloadRulesResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_mago_proto protoreflect.FileDescriptor

var file_mago_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6d, 0x61, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6d, 0x61,
	0x67, 0x6f, 0x22, 0x0f, 0x0a, 0x0d, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x28, 0x0a, 0x0e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x2a, 0x0a,
	0x12, 0x52, 0x65, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x2f, 0x0a, 0x13, 0x52, 0x65, 0x6c,
	0x6f, 0x61, 0x64, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x8d, 0x01, 0x0a, 0x04, 0x4d,
	0x61, 0x67, 0x6f, 0x12, 0x35, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x13, 0x2e,
	0x6d, 0x61, 0x67, 0x6f, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x14, 0x2e, 0x6d, 0x61, 0x67, 0x6f, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4e, 0x0a, 0x0b, 0x52, 0x65,
	0x6c, 0x6f, 0x61, 0x64, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x1d, 0x2e, 0x6d, 0x61, 0x67, 0x6f,
	0x2e, 0x52, 0x65, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x6d, 0x61, 0x67, 0x6f, 0x2e, 0x52, 0x65, 0x6c, 0x6f, 0x61,
	0x64, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62,
	0x65, 0x64, 0x75, 0x73, 0x65, 0x63, 0x2f, 0x6d, 0x61, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mago_proto_rawDescOnce sync.Once
	file_mago_proto_rawDescData = file_mago_proto_rawDesc
)

func file_mago_proto_rawDescGZIP() []byte {
	file_mago_proto_rawDescOnce.Do(func() {
		file_mago_proto_rawDescData = protoimpl.X.CompressGZIP(file_mago_proto_rawDescData)
	})
	return file_mago_proto_rawDescData
}

var file_mago_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_mago_proto_goTypes = []interface{}{
	(*HealthRequest)(nil),       // 0: mago.HealthRequest
	(*HealthResponse)(nil),      // 1: mago.HealthResponse
	(*ReloadRulesRequest)(nil),  // 2: mago.ReloadRulesRequest
	(*ReloadRulesResponse)(nil), // 3: mago.ReloadRulesResponse
}
var file_mago_proto_depIdxs = []int32{
	0, // 0: mago.Mago.Health:input_type -> mago.HealthRequest
	2, // 1: mago.Mago.ReloadRules:input_type -> mago.ReloadRulesRequest
	1, // 2: mago.Mago.Health:output_type -> mago.HealthResponse
	3, // 3: mago.Mago.ReloadRules:output_type -> mago.ReloadRulesResponse
}

func init() { file_mago_proto_init() }
func file_mago_proto_init() {
	if File_mago_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_mago_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mago_proto_goTypes,
		DependencyIndexes: file_mago_proto_depIdxs,
		MessageInfos:      file_mago_proto_msgTypes,
	}.Build()
	File_mago_proto = out.File
	file_mago_proto_rawDesc = nil
	file_mago_proto_goTypes = nil
	file_mago_proto_depIdxs = nil
}