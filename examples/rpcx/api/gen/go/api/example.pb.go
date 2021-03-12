// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.15.5
// source: api/example.proto

package api

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type EchoReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *EchoReq) Reset() {
	*x = EchoReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_example_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoReq) ProtoMessage() {}

func (x *EchoReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_example_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoReq.ProtoReflect.Descriptor instead.
func (*EchoReq) Descriptor() ([]byte, []int) {
	return file_api_example_proto_rawDescGZIP(), []int{0}
}

func (x *EchoReq) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type EchoRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *EchoRes) Reset() {
	*x = EchoRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_example_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoRes) ProtoMessage() {}

func (x *EchoRes) ProtoReflect() protoreflect.Message {
	mi := &file_api_example_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoRes.ProtoReflect.Descriptor instead.
func (*EchoRes) Descriptor() ([]byte, []int) {
	return file_api_example_proto_rawDescGZIP(), []int{1}
}

func (x *EchoRes) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type AddReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	I1 int64 `protobuf:"varint,1,opt,name=i1,proto3" json:"i1,omitempty"`
	I2 int64 `protobuf:"varint,2,opt,name=i2,proto3" json:"i2,omitempty"`
}

func (x *AddReq) Reset() {
	*x = AddReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_example_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddReq) ProtoMessage() {}

func (x *AddReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_example_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddReq.ProtoReflect.Descriptor instead.
func (*AddReq) Descriptor() ([]byte, []int) {
	return file_api_example_proto_rawDescGZIP(), []int{2}
}

func (x *AddReq) GetI1() int64 {
	if x != nil {
		return x.I1
	}
	return 0
}

func (x *AddReq) GetI2() int64 {
	if x != nil {
		return x.I2
	}
	return 0
}

type AddRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Val int64 `protobuf:"varint,1,opt,name=val,proto3" json:"val,omitempty"`
}

func (x *AddRes) Reset() {
	*x = AddRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_example_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddRes) ProtoMessage() {}

func (x *AddRes) ProtoReflect() protoreflect.Message {
	mi := &file_api_example_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddRes.ProtoReflect.Descriptor instead.
func (*AddRes) Descriptor() ([]byte, []int) {
	return file_api_example_proto_rawDescGZIP(), []int{3}
}

func (x *AddRes) GetVal() int64 {
	if x != nil {
		return x.Val
	}
	return 0
}

var File_api_example_proto protoreflect.FileDescriptor

var file_api_example_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x23, 0x0a, 0x07, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65,
	0x71, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x23, 0x0a, 0x07, 0x45,
	0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0x28, 0x0a, 0x06, 0x41, 0x64, 0x64, 0x52, 0x65, 0x71, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x31,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x31, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x32,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x32, 0x22, 0x1a, 0x0a, 0x06, 0x41, 0x64,
	0x64, 0x52, 0x65, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x32, 0x7b, 0x0a, 0x0e, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x34, 0x0a, 0x04, 0x45, 0x63, 0x68, 0x6f,
	0x12, 0x0c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x1a, 0x0c,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x22, 0x10, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x0a, 0x12, 0x08, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x63, 0x68, 0x6f, 0x12, 0x33,
	0x0a, 0x03, 0x41, 0x64, 0x64, 0x12, 0x0b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x64, 0x64, 0x52,
	0x65, 0x71, 0x1a, 0x0b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65, 0x73, 0x22,
	0x12, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x22, 0x07, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x64, 0x64,
	0x3a, 0x01, 0x2a, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x68, 0x61, 0x74, 0x6c, 0x6f, 0x6e, 0x65, 0x6c, 0x79, 0x2f, 0x67, 0x6f, 0x2d, 0x6b,
	0x69, 0x74, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x72, 0x70, 0x63, 0x78, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_example_proto_rawDescOnce sync.Once
	file_api_example_proto_rawDescData = file_api_example_proto_rawDesc
)

func file_api_example_proto_rawDescGZIP() []byte {
	file_api_example_proto_rawDescOnce.Do(func() {
		file_api_example_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_example_proto_rawDescData)
	})
	return file_api_example_proto_rawDescData
}

var file_api_example_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_example_proto_goTypes = []interface{}{
	(*EchoReq)(nil), // 0: api.EchoReq
	(*EchoRes)(nil), // 1: api.EchoRes
	(*AddReq)(nil),  // 2: api.AddReq
	(*AddRes)(nil),  // 3: api.AddRes
}
var file_api_example_proto_depIdxs = []int32{
	0, // 0: api.ExampleService.Echo:input_type -> api.EchoReq
	2, // 1: api.ExampleService.Add:input_type -> api.AddReq
	1, // 2: api.ExampleService.Echo:output_type -> api.EchoRes
	3, // 3: api.ExampleService.Add:output_type -> api.AddRes
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_example_proto_init() }
func file_api_example_proto_init() {
	if File_api_example_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_example_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoReq); i {
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
		file_api_example_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoRes); i {
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
		file_api_example_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddReq); i {
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
		file_api_example_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddRes); i {
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
			RawDescriptor: file_api_example_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_example_proto_goTypes,
		DependencyIndexes: file_api_example_proto_depIdxs,
		MessageInfos:      file_api_example_proto_msgTypes,
	}.Build()
	File_api_example_proto = out.File
	file_api_example_proto_rawDesc = nil
	file_api_example_proto_goTypes = nil
	file_api_example_proto_depIdxs = nil
}