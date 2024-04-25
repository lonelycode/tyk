// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v5.26.1
// source: coprocess_mini_request_object.proto

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

// Used for middleware calls and contains important fields like headers, parameters, body and URL.
type MiniRequestObject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A read-only field for reading headers injected by previous middleware
	Headers map[string]string `protobuf:"bytes,1,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Map of header key values to append to the request
	SetHeaders map[string]string `protobuf:"bytes,2,rep,name=set_headers,json=setHeaders,proto3" json:"set_headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// List of header names to be removed from the request
	DeleteHeaders []string `protobuf:"bytes,3,rep,name=delete_headers,json=deleteHeaders,proto3" json:"delete_headers,omitempty"`
	// Request body
	Body string `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	// Request URL
	Url string `protobuf:"bytes,5,opt,name=url,proto3" json:"url,omitempty"`
	// Read only map of request params
	Params map[string]string `protobuf:"bytes,6,rep,name=params,proto3" json:"params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Map of parameter keys and values to add to the request
	AddParams map[string]string `protobuf:"bytes,7,rep,name=add_params,json=addParams,proto3" json:"add_params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Allows a parameter to have multiple values, currently unsupported
	ExtendedParams map[string]string `protobuf:"bytes,8,rep,name=extended_params,json=extendedParams,proto3" json:"extended_params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// List of parameter keys to be removed from the request
	DeleteParams []string `protobuf:"bytes,9,rep,name=delete_params,json=deleteParams,proto3" json:"delete_params,omitempty"`
	// Override the response for the request, see ReturnOverrides
	ReturnOverrides *ReturnOverrides `protobuf:"bytes,10,opt,name=return_overrides,json=returnOverrides,proto3" json:"return_overrides,omitempty"`
	// Request method, eg GET, POST, etc
	Method string `protobuf:"bytes,11,opt,name=method,proto3" json:"method,omitempty"`
	// The raw unprocessed request URL, including query string and fragments
	RequestUri string `protobuf:"bytes,12,opt,name=request_uri,json=requestUri,proto3" json:"request_uri,omitempty"`
	// The URL scheme, e.g. http or https
	Scheme string `protobuf:"bytes,13,opt,name=scheme,proto3" json:"scheme,omitempty"`
	// The raw request body
	RawBody []byte `protobuf:"bytes,14,opt,name=raw_body,json=rawBody,proto3" json:"raw_body,omitempty"`
}

func (x *MiniRequestObject) Reset() {
	*x = MiniRequestObject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coprocess_mini_request_object_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MiniRequestObject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MiniRequestObject) ProtoMessage() {}

func (x *MiniRequestObject) ProtoReflect() protoreflect.Message {
	mi := &file_coprocess_mini_request_object_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MiniRequestObject.ProtoReflect.Descriptor instead.
func (*MiniRequestObject) Descriptor() ([]byte, []int) {
	return file_coprocess_mini_request_object_proto_rawDescGZIP(), []int{0}
}

func (x *MiniRequestObject) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *MiniRequestObject) GetSetHeaders() map[string]string {
	if x != nil {
		return x.SetHeaders
	}
	return nil
}

func (x *MiniRequestObject) GetDeleteHeaders() []string {
	if x != nil {
		return x.DeleteHeaders
	}
	return nil
}

func (x *MiniRequestObject) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *MiniRequestObject) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *MiniRequestObject) GetParams() map[string]string {
	if x != nil {
		return x.Params
	}
	return nil
}

func (x *MiniRequestObject) GetAddParams() map[string]string {
	if x != nil {
		return x.AddParams
	}
	return nil
}

func (x *MiniRequestObject) GetExtendedParams() map[string]string {
	if x != nil {
		return x.ExtendedParams
	}
	return nil
}

func (x *MiniRequestObject) GetDeleteParams() []string {
	if x != nil {
		return x.DeleteParams
	}
	return nil
}

func (x *MiniRequestObject) GetReturnOverrides() *ReturnOverrides {
	if x != nil {
		return x.ReturnOverrides
	}
	return nil
}

func (x *MiniRequestObject) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *MiniRequestObject) GetRequestUri() string {
	if x != nil {
		return x.RequestUri
	}
	return ""
}

func (x *MiniRequestObject) GetScheme() string {
	if x != nil {
		return x.Scheme
	}
	return ""
}

func (x *MiniRequestObject) GetRawBody() []byte {
	if x != nil {
		return x.RawBody
	}
	return nil
}

var File_coprocess_mini_request_object_proto protoreflect.FileDescriptor

var file_coprocess_mini_request_object_proto_rawDesc = []byte{
	0x0a, 0x23, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x6d, 0x69, 0x6e, 0x69,
	0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x1a, 0x20, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x74, 0x75,
	0x72, 0x6e, 0x5f, 0x6f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xec, 0x07, 0x0a, 0x11, 0x4d, 0x69, 0x6e, 0x69, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x43, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x63, 0x6f, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x4d, 0x69, 0x6e, 0x69, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x4d, 0x0a,
	0x0b, 0x73, 0x65, 0x74, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x4d,
	0x69, 0x6e, 0x69, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x2e, 0x53, 0x65, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x0a, 0x73, 0x65, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x25, 0x0a, 0x0e,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x40, 0x0a, 0x06, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x63, 0x6f, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x4d, 0x69, 0x6e, 0x69, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x4a, 0x0a, 0x0a, 0x61,
	0x64, 0x64, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x2b, 0x2e, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x4d, 0x69, 0x6e, 0x69,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x41, 0x64,
	0x64, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x61, 0x64,
	0x64, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x59, 0x0a, 0x0f, 0x65, 0x78, 0x74, 0x65, 0x6e,
	0x64, 0x65, 0x64, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x30, 0x2e, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x4d, 0x69, 0x6e,
	0x69, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x45,
	0x78, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x0e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x45, 0x0a, 0x10, 0x72, 0x65, 0x74, 0x75, 0x72,
	0x6e, 0x5f, 0x6f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x52, 0x65,
	0x74, 0x75, 0x72, 0x6e, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x73, 0x52, 0x0f, 0x72,
	0x65, 0x74, 0x75, 0x72, 0x6e, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x55, 0x72, 0x69, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d,
	0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x72, 0x61, 0x77, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x72, 0x61, 0x77, 0x42, 0x6f, 0x64, 0x79, 0x1a, 0x3a, 0x0a, 0x0c, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3d, 0x0a, 0x0f, 0x53, 0x65, 0x74, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x39, 0x0a, 0x0b, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x1a, 0x3c, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x41,
	0x0a, 0x13, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x42, 0x0c, 0x5a, 0x0a, 0x2f, 0x63, 0x6f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_coprocess_mini_request_object_proto_rawDescOnce sync.Once
	file_coprocess_mini_request_object_proto_rawDescData = file_coprocess_mini_request_object_proto_rawDesc
)

func file_coprocess_mini_request_object_proto_rawDescGZIP() []byte {
	file_coprocess_mini_request_object_proto_rawDescOnce.Do(func() {
		file_coprocess_mini_request_object_proto_rawDescData = protoimpl.X.CompressGZIP(file_coprocess_mini_request_object_proto_rawDescData)
	})
	return file_coprocess_mini_request_object_proto_rawDescData
}

var file_coprocess_mini_request_object_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_coprocess_mini_request_object_proto_goTypes = []interface{}{
	(*MiniRequestObject)(nil), // 0: coprocess.MiniRequestObject
	nil,                       // 1: coprocess.MiniRequestObject.HeadersEntry
	nil,                       // 2: coprocess.MiniRequestObject.SetHeadersEntry
	nil,                       // 3: coprocess.MiniRequestObject.ParamsEntry
	nil,                       // 4: coprocess.MiniRequestObject.AddParamsEntry
	nil,                       // 5: coprocess.MiniRequestObject.ExtendedParamsEntry
	(*ReturnOverrides)(nil),   // 6: coprocess.ReturnOverrides
}
var file_coprocess_mini_request_object_proto_depIdxs = []int32{
	1, // 0: coprocess.MiniRequestObject.headers:type_name -> coprocess.MiniRequestObject.HeadersEntry
	2, // 1: coprocess.MiniRequestObject.set_headers:type_name -> coprocess.MiniRequestObject.SetHeadersEntry
	3, // 2: coprocess.MiniRequestObject.params:type_name -> coprocess.MiniRequestObject.ParamsEntry
	4, // 3: coprocess.MiniRequestObject.add_params:type_name -> coprocess.MiniRequestObject.AddParamsEntry
	5, // 4: coprocess.MiniRequestObject.extended_params:type_name -> coprocess.MiniRequestObject.ExtendedParamsEntry
	6, // 5: coprocess.MiniRequestObject.return_overrides:type_name -> coprocess.ReturnOverrides
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_coprocess_mini_request_object_proto_init() }
func file_coprocess_mini_request_object_proto_init() {
	if File_coprocess_mini_request_object_proto != nil {
		return
	}
	file_coprocess_return_overrides_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_coprocess_mini_request_object_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MiniRequestObject); i {
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
			RawDescriptor: file_coprocess_mini_request_object_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_coprocess_mini_request_object_proto_goTypes,
		DependencyIndexes: file_coprocess_mini_request_object_proto_depIdxs,
		MessageInfos:      file_coprocess_mini_request_object_proto_msgTypes,
	}.Build()
	File_coprocess_mini_request_object_proto = out.File
	file_coprocess_mini_request_object_proto_rawDesc = nil
	file_coprocess_mini_request_object_proto_goTypes = nil
	file_coprocess_mini_request_object_proto_depIdxs = nil
}
