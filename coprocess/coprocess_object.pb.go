// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v5.26.1
// source: coprocess_object.proto

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

// Wraps a MiniRequestObject and contains additional fields that are useful for users that implement
// their own request dispatchers, like the middleware hook type and name
type Object struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The plugin hook type
	HookType HookType `protobuf:"varint,1,opt,name=hook_type,json=hookType,proto3,enum=coprocess.HookType" json:"hook_type,omitempty"`
	// The plugin name
	HookName string `protobuf:"bytes,2,opt,name=hook_name,json=hookName,proto3" json:"hook_name,omitempty"`
	// The main request data structure used by rich plugins. It’s used for middleware calls
	// and contains important fields like headers, parameters, body and URL
	Request *MiniRequestObject `protobuf:"bytes,3,opt,name=request,proto3" json:"request,omitempty"`
	// Stores information about the current key/user that’s used for authentication
	Session *SessionState `protobuf:"bytes,4,opt,name=session,proto3" json:"session,omitempty"`
	// Contains the metadata. This is a dynamic field
	Metadata map[string]string `protobuf:"bytes,5,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Contains information about API definition, including APIID, OrgID and config_data
	Spec map[string]string `protobuf:"bytes,6,rep,name=spec,proto3" json:"spec,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// The ResponseObject is used by response hooks. The fields are populated with the upstream HTTP
	// response data. All the field contents can be modified.
	Response *ResponseObject `protobuf:"bytes,7,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *Object) Reset() {
	*x = Object{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coprocess_object_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Object) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Object) ProtoMessage() {}

func (x *Object) ProtoReflect() protoreflect.Message {
	mi := &file_coprocess_object_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Object.ProtoReflect.Descriptor instead.
func (*Object) Descriptor() ([]byte, []int) {
	return file_coprocess_object_proto_rawDescGZIP(), []int{0}
}

func (x *Object) GetHookType() HookType {
	if x != nil {
		return x.HookType
	}
	return HookType_Unknown
}

func (x *Object) GetHookName() string {
	if x != nil {
		return x.HookName
	}
	return ""
}

func (x *Object) GetRequest() *MiniRequestObject {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *Object) GetSession() *SessionState {
	if x != nil {
		return x.Session
	}
	return nil
}

func (x *Object) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Object) GetSpec() map[string]string {
	if x != nil {
		return x.Spec
	}
	return nil
}

func (x *Object) GetResponse() *ResponseObject {
	if x != nil {
		return x.Response
	}
	return nil
}

// An event is represented as a JSON payload
type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The JSON payload
	Payload string `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coprocess_object_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_coprocess_object_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_coprocess_object_proto_rawDescGZIP(), []int{1}
}

func (x *Event) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

// Response for event
type EventReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EventReply) Reset() {
	*x = EventReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coprocess_object_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventReply) ProtoMessage() {}

func (x *EventReply) ProtoReflect() protoreflect.Message {
	mi := &file_coprocess_object_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventReply.ProtoReflect.Descriptor instead.
func (*EventReply) Descriptor() ([]byte, []int) {
	return file_coprocess_object_proto_rawDescGZIP(), []int{2}
}

var File_coprocess_object_proto protoreflect.FileDescriptor

var file_coprocess_object_proto_rawDesc = []byte{
	0x0a, 0x16, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x1a, 0x23, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x6d,
	0x69, 0x6e, 0x69, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x6f, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x63, 0x6f, 0x70, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xdd, 0x03, 0x0a, 0x06, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x30, 0x0a, 0x09, 0x68,
	0x6f, 0x6f, 0x6b, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13,
	0x2e, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x48, 0x6f, 0x6f, 0x6b, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x08, 0x68, 0x6f, 0x6f, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a,
	0x09, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x68, 0x6f, 0x6f, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x6f,
	0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x4d, 0x69, 0x6e, 0x69, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x31, 0x0a, 0x07, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x07, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3b, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x2f, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x4f, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x2e, 0x53, 0x70, 0x65, 0x63, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x73,
	0x70, 0x65, 0x63, 0x12, 0x35, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x37, 0x0a, 0x09, 0x53, 0x70, 0x65, 0x63, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x21, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x22, 0x0c, 0x0a, 0x0a, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x32, 0x7c, 0x0a, 0x0a, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x12,
	0x32, 0x0a, 0x08, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x12, 0x11, 0x2e, 0x63, 0x6f,
	0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x1a, 0x11,
	0x2e, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x0d, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x2e, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x15, 0x2e, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42,
	0x0c, 0x5a, 0x0a, 0x2f, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_coprocess_object_proto_rawDescOnce sync.Once
	file_coprocess_object_proto_rawDescData = file_coprocess_object_proto_rawDesc
)

func file_coprocess_object_proto_rawDescGZIP() []byte {
	file_coprocess_object_proto_rawDescOnce.Do(func() {
		file_coprocess_object_proto_rawDescData = protoimpl.X.CompressGZIP(file_coprocess_object_proto_rawDescData)
	})
	return file_coprocess_object_proto_rawDescData
}

var file_coprocess_object_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_coprocess_object_proto_goTypes = []interface{}{
	(*Object)(nil),            // 0: coprocess.Object
	(*Event)(nil),             // 1: coprocess.Event
	(*EventReply)(nil),        // 2: coprocess.EventReply
	nil,                       // 3: coprocess.Object.MetadataEntry
	nil,                       // 4: coprocess.Object.SpecEntry
	(HookType)(0),             // 5: coprocess.HookType
	(*MiniRequestObject)(nil), // 6: coprocess.MiniRequestObject
	(*SessionState)(nil),      // 7: coprocess.SessionState
	(*ResponseObject)(nil),    // 8: coprocess.ResponseObject
}
var file_coprocess_object_proto_depIdxs = []int32{
	5, // 0: coprocess.Object.hook_type:type_name -> coprocess.HookType
	6, // 1: coprocess.Object.request:type_name -> coprocess.MiniRequestObject
	7, // 2: coprocess.Object.session:type_name -> coprocess.SessionState
	3, // 3: coprocess.Object.metadata:type_name -> coprocess.Object.MetadataEntry
	4, // 4: coprocess.Object.spec:type_name -> coprocess.Object.SpecEntry
	8, // 5: coprocess.Object.response:type_name -> coprocess.ResponseObject
	0, // 6: coprocess.Dispatcher.Dispatch:input_type -> coprocess.Object
	1, // 7: coprocess.Dispatcher.DispatchEvent:input_type -> coprocess.Event
	0, // 8: coprocess.Dispatcher.Dispatch:output_type -> coprocess.Object
	2, // 9: coprocess.Dispatcher.DispatchEvent:output_type -> coprocess.EventReply
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_coprocess_object_proto_init() }
func file_coprocess_object_proto_init() {
	if File_coprocess_object_proto != nil {
		return
	}
	file_coprocess_mini_request_object_proto_init()
	file_coprocess_response_object_proto_init()
	file_coprocess_session_state_proto_init()
	file_coprocess_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_coprocess_object_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Object); i {
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
		file_coprocess_object_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
		file_coprocess_object_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventReply); i {
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
			RawDescriptor: file_coprocess_object_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_coprocess_object_proto_goTypes,
		DependencyIndexes: file_coprocess_object_proto_depIdxs,
		MessageInfos:      file_coprocess_object_proto_msgTypes,
	}.Build()
	File_coprocess_object_proto = out.File
	file_coprocess_object_proto_rawDesc = nil
	file_coprocess_object_proto_goTypes = nil
	file_coprocess_object_proto_depIdxs = nil
}
