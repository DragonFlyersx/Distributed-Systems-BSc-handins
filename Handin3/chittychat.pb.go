// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.0--rc1
// source: ChittyChat.proto

// FOR SERVER

package Handin3

import (
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

type ChatMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message   string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Timestamp int64  `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *ChatMessage) Reset() {
	*x = ChatMessage{}
	mi := &file_ChittyChat_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMessage) ProtoMessage() {}

func (x *ChatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_ChittyChat_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMessage.ProtoReflect.Descriptor instead.
func (*ChatMessage) Descriptor() ([]byte, []int) {
	return file_ChittyChat_proto_rawDescGZIP(), []int{0}
}

func (x *ChatMessage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ChatMessage) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_ChittyChat_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_ChittyChat_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_ChittyChat_proto_rawDescGZIP(), []int{1}
}

var File_ChittyChat_proto protoreflect.FileDescriptor

var file_ChittyChat_proto_rawDesc = []byte{
	0x0a, 0x10, 0x43, 0x68, 0x69, 0x74, 0x74, 0x79, 0x43, 0x68, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x6c, 0x6f, 0x67, 0x22, 0x45, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x74, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x07,
	0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x70, 0x0a, 0x0a, 0x43, 0x68, 0x69, 0x74, 0x74,
	0x79, 0x43, 0x68, 0x61, 0x74, 0x12, 0x32, 0x0a, 0x10, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61,
	0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0a, 0x2e, 0x6c, 0x6f, 0x67, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x10, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x43, 0x68, 0x61, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x30, 0x01, 0x12, 0x2e, 0x0a, 0x0e, 0x50, 0x75, 0x62,
	0x6c, 0x69, 0x73, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x10, 0x2e, 0x6c, 0x6f,
	0x67, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0a, 0x2e,
	0x6c, 0x6f, 0x67, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0a, 0x5a, 0x08, 0x2f, 0x48, 0x61,
	0x6e, 0x64, 0x69, 0x6e, 0x33, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ChittyChat_proto_rawDescOnce sync.Once
	file_ChittyChat_proto_rawDescData = file_ChittyChat_proto_rawDesc
)

func file_ChittyChat_proto_rawDescGZIP() []byte {
	file_ChittyChat_proto_rawDescOnce.Do(func() {
		file_ChittyChat_proto_rawDescData = protoimpl.X.CompressGZIP(file_ChittyChat_proto_rawDescData)
	})
	return file_ChittyChat_proto_rawDescData
}

var file_ChittyChat_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_ChittyChat_proto_goTypes = []any{
	(*ChatMessage)(nil), // 0: log.ChatMessage
	(*Empty)(nil),       // 1: log.Empty
}
var file_ChittyChat_proto_depIdxs = []int32{
	1, // 0: log.ChittyChat.BroadcastMessage:input_type -> log.Empty
	0, // 1: log.ChittyChat.PublishMessage:input_type -> log.ChatMessage
	0, // 2: log.ChittyChat.BroadcastMessage:output_type -> log.ChatMessage
	1, // 3: log.ChittyChat.PublishMessage:output_type -> log.Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ChittyChat_proto_init() }
func file_ChittyChat_proto_init() {
	if File_ChittyChat_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ChittyChat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ChittyChat_proto_goTypes,
		DependencyIndexes: file_ChittyChat_proto_depIdxs,
		MessageInfos:      file_ChittyChat_proto_msgTypes,
	}.Build()
	File_ChittyChat_proto = out.File
	file_ChittyChat_proto_rawDesc = nil
	file_ChittyChat_proto_goTypes = nil
	file_ChittyChat_proto_depIdxs = nil
}
