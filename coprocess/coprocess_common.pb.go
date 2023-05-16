// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.1
// 	protoc        v3.21.12
// source: coprocess_common.proto

package coprocess

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

type HookType int32

const (
	HookType_Unknown        HookType = 0
	HookType_Pre            HookType = 1
	HookType_Post           HookType = 2
	HookType_PostKeyAuth    HookType = 3
	HookType_CustomKeyCheck HookType = 4
	HookType_Response       HookType = 5
)

// Enum value maps for HookType.
var (
	HookType_name = map[int32]string{
		0: "Unknown",
		1: "Pre",
		2: "Post",
		3: "PostKeyAuth",
		4: "CustomKeyCheck",
		5: "Response",
	}
	HookType_value = map[string]int32{
		"Unknown":        0,
		"Pre":            1,
		"Post":           2,
		"PostKeyAuth":    3,
		"CustomKeyCheck": 4,
		"Response":       5,
	}
)

func (x HookType) Enum() *HookType {
	p := new(HookType)
	*p = x
	return p
}

func (x HookType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HookType) Descriptor() protoreflect.EnumDescriptor {
	return file_coprocess_common_proto_enumTypes[0].Descriptor()
}

func (HookType) Type() protoreflect.EnumType {
	return &file_coprocess_common_proto_enumTypes[0]
}

func (x HookType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HookType.Descriptor instead.
func (HookType) EnumDescriptor() ([]byte, []int) {
	return file_coprocess_common_proto_rawDescGZIP(), []int{0}
}

type StringSlice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []string `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *StringSlice) Reset() {
	*x = StringSlice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coprocess_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringSlice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringSlice) ProtoMessage() {}

func (x *StringSlice) ProtoReflect() protoreflect.Message {
	mi := &file_coprocess_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringSlice.ProtoReflect.Descriptor instead.
func (*StringSlice) Descriptor() ([]byte, []int) {
	return file_coprocess_common_proto_rawDescGZIP(), []int{0}
}

func (x *StringSlice) GetItems() []string {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_coprocess_common_proto protoreflect.FileDescriptor

var file_coprocess_common_proto_rawDesc = []byte{
	0x0a, 0x16, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x22, 0x23, 0x0a, 0x0b, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x6c, 0x69,
	0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x2a, 0x5d, 0x0a, 0x08, 0x48, 0x6f, 0x6f, 0x6b,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10,
	0x00, 0x12, 0x07, 0x0a, 0x03, 0x50, 0x72, 0x65, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x6f,
	0x73, 0x74, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x50, 0x6f, 0x73, 0x74, 0x4b, 0x65, 0x79, 0x41,
	0x75, 0x74, 0x68, 0x10, 0x03, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x4b,
	0x65, 0x79, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x10, 0x04, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x10, 0x05, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_coprocess_common_proto_rawDescOnce sync.Once
	file_coprocess_common_proto_rawDescData = file_coprocess_common_proto_rawDesc
)

func file_coprocess_common_proto_rawDescGZIP() []byte {
	file_coprocess_common_proto_rawDescOnce.Do(func() {
		file_coprocess_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_coprocess_common_proto_rawDescData)
	})
	return file_coprocess_common_proto_rawDescData
}

var file_coprocess_common_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_coprocess_common_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_coprocess_common_proto_goTypes = []interface{}{
	(HookType)(0),       // 0: coprocess.HookType
	(*StringSlice)(nil), // 1: coprocess.StringSlice
}
var file_coprocess_common_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_coprocess_common_proto_init() }
func file_coprocess_common_proto_init() {
	if File_coprocess_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_coprocess_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringSlice); i {
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
			RawDescriptor: file_coprocess_common_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_coprocess_common_proto_goTypes,
		DependencyIndexes: file_coprocess_common_proto_depIdxs,
		EnumInfos:         file_coprocess_common_proto_enumTypes,
		MessageInfos:      file_coprocess_common_proto_msgTypes,
	}.Build()
	File_coprocess_common_proto = out.File
	file_coprocess_common_proto_rawDesc = nil
	file_coprocess_common_proto_goTypes = nil
	file_coprocess_common_proto_depIdxs = nil
}
