// Code generated by protoc-gen-go. DO NOT EDIT.
// source: coprocess_return_overrides.proto

package coprocess

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ReturnOverrides struct {
	ResponseCode         int32             `protobuf:"varint,1,opt,name=response_code,json=responseCode,proto3" json:"response_code,omitempty"`
	ResponseError        string            `protobuf:"bytes,2,opt,name=response_error,json=responseError,proto3" json:"response_error,omitempty"`
	Headers              map[string]string `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	OverrideError        bool              `protobuf:"varint,4,opt,name=override_error,json=overrideError,proto3" json:"override_error,omitempty"`
	ResponseBody         string            `protobuf:"bytes,5,opt,name=response_body,json=responseBody,proto3" json:"response_body,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ReturnOverrides) Reset()         { *m = ReturnOverrides{} }
func (m *ReturnOverrides) String() string { return proto.CompactTextString(m) }
func (*ReturnOverrides) ProtoMessage()    {}
func (*ReturnOverrides) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c6abd8ea4a81548, []int{0}
}

func (m *ReturnOverrides) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReturnOverrides.Unmarshal(m, b)
}
func (m *ReturnOverrides) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReturnOverrides.Marshal(b, m, deterministic)
}
func (m *ReturnOverrides) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReturnOverrides.Merge(m, src)
}
func (m *ReturnOverrides) XXX_Size() int {
	return xxx_messageInfo_ReturnOverrides.Size(m)
}
func (m *ReturnOverrides) XXX_DiscardUnknown() {
	xxx_messageInfo_ReturnOverrides.DiscardUnknown(m)
}

var xxx_messageInfo_ReturnOverrides proto.InternalMessageInfo

func (m *ReturnOverrides) GetResponseCode() int32 {
	if m != nil {
		return m.ResponseCode
	}
	return 0
}

func (m *ReturnOverrides) GetResponseError() string {
	if m != nil {
		return m.ResponseError
	}
	return ""
}

func (m *ReturnOverrides) GetHeaders() map[string]string {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *ReturnOverrides) GetOverrideError() bool {
	if m != nil {
		return m.OverrideError
	}
	return false
}

func (m *ReturnOverrides) GetResponseBody() string {
	if m != nil {
		return m.ResponseBody
	}
	return ""
}

func init() {
	proto.RegisterType((*ReturnOverrides)(nil), "coprocess.ReturnOverrides")
	proto.RegisterMapType((map[string]string)(nil), "coprocess.ReturnOverrides.HeadersEntry")
}

func init() { proto.RegisterFile("coprocess_return_overrides.proto", fileDescriptor_7c6abd8ea4a81548) }

var fileDescriptor_7c6abd8ea4a81548 = []byte{
	// 235 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xc1, 0x4a, 0xc4, 0x30,
	0x10, 0x86, 0x49, 0x6b, 0xd5, 0x8e, 0xbb, 0x2a, 0xc1, 0x43, 0xf1, 0x14, 0x14, 0xb1, 0xa7, 0x1c,
	0xf4, 0x22, 0x7b, 0x53, 0x59, 0xf0, 0x26, 0xe4, 0x05, 0xca, 0x6e, 0x33, 0xa0, 0x28, 0x9d, 0x32,
	0xe9, 0x2e, 0xe4, 0x8d, 0x7c, 0x4c, 0x69, 0xba, 0x09, 0x65, 0x6f, 0xed, 0xcf, 0x37, 0x7c, 0xff,
	0x1f, 0x50, 0x2d, 0xf5, 0x4c, 0x2d, 0x3a, 0xd7, 0x30, 0x0e, 0x3b, 0xee, 0x1a, 0xda, 0x23, 0xf3,
	0xb7, 0x45, 0xa7, 0x7b, 0xa6, 0x81, 0x64, 0x99, 0x88, 0xbb, 0xbf, 0x0c, 0xae, 0x4c, 0xa0, 0x3e,
	0x23, 0x24, 0xef, 0x61, 0xc9, 0xe8, 0x7a, 0xea, 0x1c, 0x36, 0x2d, 0x59, 0xac, 0x84, 0x12, 0x75,
	0x61, 0x16, 0x31, 0x7c, 0x27, 0x8b, 0xf2, 0x01, 0x2e, 0x13, 0x84, 0xcc, 0xc4, 0x55, 0xa6, 0x44,
	0x5d, 0x9a, 0x74, 0xba, 0x1e, 0x43, 0xf9, 0x0a, 0x67, 0x5f, 0xb8, 0xb1, 0xc8, 0xae, 0xca, 0x55,
	0x5e, 0x5f, 0x3c, 0x3d, 0xea, 0x24, 0xd7, 0x47, 0x62, 0xfd, 0x31, 0x91, 0xeb, 0x6e, 0x60, 0x6f,
	0xe2, 0xdd, 0x68, 0x8a, 0x03, 0x0e, 0xa6, 0x13, 0x25, 0xea, 0x73, 0xb3, 0x8c, 0xe9, 0x64, 0x9a,
	0xb7, 0xde, 0x92, 0xf5, 0x55, 0x11, 0xfa, 0xa4, 0xd6, 0x6f, 0x64, 0xfd, 0xed, 0x0a, 0x16, 0x73,
	0x89, 0xbc, 0x86, 0xfc, 0x07, 0x7d, 0x18, 0x58, 0x9a, 0xf1, 0x53, 0xde, 0x40, 0xb1, 0xdf, 0xfc,
	0xee, 0xf0, 0x30, 0x67, 0xfa, 0x59, 0x65, 0x2f, 0x62, 0x7b, 0x1a, 0x1e, 0xef, 0xf9, 0x3f, 0x00,
	0x00, 0xff, 0xff, 0x21, 0x7a, 0xdb, 0xef, 0x60, 0x01, 0x00, 0x00,
}
