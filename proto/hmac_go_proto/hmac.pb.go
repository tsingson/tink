// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
////////////////////////////////////////////////////////////////////////////////

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: third_party/tink/proto/hmac.proto

package hmac_go_proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	common_go_proto "github.com/tsingson/tink/proto/common_go_proto"
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

type HmacParams struct {
	Hash                 common_go_proto.HashType `protobuf:"varint,1,opt,name=hash,proto3,enum=google.crypto.tink.HashType" json:"hash,omitempty"`
	TagSize              uint32                   `protobuf:"varint,2,opt,name=tag_size,json=tagSize,proto3" json:"tag_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *HmacParams) Reset()         { *m = HmacParams{} }
func (m *HmacParams) String() string { return proto.CompactTextString(m) }
func (*HmacParams) ProtoMessage()    {}
func (*HmacParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_310803c785e2f4dc, []int{0}
}

func (m *HmacParams) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HmacParams.Unmarshal(m, b)
}
func (m *HmacParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HmacParams.Marshal(b, m, deterministic)
}
func (m *HmacParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HmacParams.Merge(m, src)
}
func (m *HmacParams) XXX_Size() int {
	return xxx_messageInfo_HmacParams.Size(m)
}
func (m *HmacParams) XXX_DiscardUnknown() {
	xxx_messageInfo_HmacParams.DiscardUnknown(m)
}

var xxx_messageInfo_HmacParams proto.InternalMessageInfo

func (m *HmacParams) GetHash() common_go_proto.HashType {
	if m != nil {
		return m.Hash
	}
	return common_go_proto.HashType_UNKNOWN_HASH
}

func (m *HmacParams) GetTagSize() uint32 {
	if m != nil {
		return m.TagSize
	}
	return 0
}

// key_type: type.googleapis.com/google.crypto.tink.HmacKey
type HmacKey struct {
	Version              uint32      `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Params               *HmacParams `protobuf:"bytes,2,opt,name=params,proto3" json:"params,omitempty"`
	KeyValue             []byte      `protobuf:"bytes,3,opt,name=key_value,json=keyValue,proto3" json:"key_value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *HmacKey) Reset()         { *m = HmacKey{} }
func (m *HmacKey) String() string { return proto.CompactTextString(m) }
func (*HmacKey) ProtoMessage()    {}
func (*HmacKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_310803c785e2f4dc, []int{1}
}

func (m *HmacKey) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HmacKey.Unmarshal(m, b)
}
func (m *HmacKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HmacKey.Marshal(b, m, deterministic)
}
func (m *HmacKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HmacKey.Merge(m, src)
}
func (m *HmacKey) XXX_Size() int {
	return xxx_messageInfo_HmacKey.Size(m)
}
func (m *HmacKey) XXX_DiscardUnknown() {
	xxx_messageInfo_HmacKey.DiscardUnknown(m)
}

var xxx_messageInfo_HmacKey proto.InternalMessageInfo

func (m *HmacKey) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *HmacKey) GetParams() *HmacParams {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *HmacKey) GetKeyValue() []byte {
	if m != nil {
		return m.KeyValue
	}
	return nil
}

type HmacKeyFormat struct {
	Params               *HmacParams `protobuf:"bytes,1,opt,name=params,proto3" json:"params,omitempty"`
	KeySize              uint32      `protobuf:"varint,2,opt,name=key_size,json=keySize,proto3" json:"key_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *HmacKeyFormat) Reset()         { *m = HmacKeyFormat{} }
func (m *HmacKeyFormat) String() string { return proto.CompactTextString(m) }
func (*HmacKeyFormat) ProtoMessage()    {}
func (*HmacKeyFormat) Descriptor() ([]byte, []int) {
	return fileDescriptor_310803c785e2f4dc, []int{2}
}

func (m *HmacKeyFormat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HmacKeyFormat.Unmarshal(m, b)
}
func (m *HmacKeyFormat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HmacKeyFormat.Marshal(b, m, deterministic)
}
func (m *HmacKeyFormat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HmacKeyFormat.Merge(m, src)
}
func (m *HmacKeyFormat) XXX_Size() int {
	return xxx_messageInfo_HmacKeyFormat.Size(m)
}
func (m *HmacKeyFormat) XXX_DiscardUnknown() {
	xxx_messageInfo_HmacKeyFormat.DiscardUnknown(m)
}

var xxx_messageInfo_HmacKeyFormat proto.InternalMessageInfo

func (m *HmacKeyFormat) GetParams() *HmacParams {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *HmacKeyFormat) GetKeySize() uint32 {
	if m != nil {
		return m.KeySize
	}
	return 0
}

func init() {
	proto.RegisterType((*HmacParams)(nil), "google.crypto.tink.HmacParams")
	proto.RegisterType((*HmacKey)(nil), "google.crypto.tink.HmacKey")
	proto.RegisterType((*HmacKeyFormat)(nil), "google.crypto.tink.HmacKeyFormat")
}

func init() { proto.RegisterFile("proto/hmac.proto", fileDescriptor_310803c785e2f4dc) }

var fileDescriptor_310803c785e2f4dc = []byte{
	// 303 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xcd, 0x4e, 0x32, 0x31,
	0x14, 0x86, 0x33, 0x7c, 0x5f, 0x00, 0x8f, 0xe2, 0xa2, 0xab, 0x41, 0x89, 0x41, 0xdc, 0x10, 0x17,
	0x1d, 0x83, 0x89, 0x17, 0xc0, 0xc2, 0x60, 0x48, 0x0c, 0x19, 0xd1, 0x44, 0x37, 0x93, 0x52, 0x9b,
	0xb6, 0x19, 0xca, 0x99, 0x74, 0x0a, 0xb1, 0x5c, 0x8e, 0x57, 0x6a, 0xa6, 0xa0, 0xf1, 0x07, 0x17,
	0xee, 0xfa, 0x26, 0xef, 0x79, 0x9e, 0x9e, 0x16, 0x4e, 0x9d, 0xd2, 0xf6, 0x39, 0x2b, 0x98, 0x75,
	0x3e, 0x71, 0x7a, 0x91, 0x27, 0x85, 0x45, 0x87, 0x89, 0x32, 0x8c, 0xd3, 0x70, 0x24, 0x44, 0x22,
	0xca, 0xb9, 0xa0, 0xdc, 0xfa, 0xc2, 0x21, 0xad, 0x4a, 0x47, 0x67, 0xbf, 0x8c, 0x71, 0x34, 0x06,
	0x17, 0x9b, 0xc1, 0xde, 0x23, 0xc0, 0xc8, 0x30, 0x3e, 0x61, 0x96, 0x99, 0x92, 0x5c, 0xc0, 0x7f,
	0xc5, 0x4a, 0x15, 0x47, 0xdd, 0xa8, 0x7f, 0x38, 0xe8, 0xd0, 0x9f, 0x54, 0x3a, 0x62, 0xa5, 0x9a,
	0xfa, 0x42, 0xa4, 0xa1, 0x49, 0xda, 0xd0, 0x74, 0x4c, 0x66, 0xa5, 0x5e, 0x8b, 0xb8, 0xd6, 0x8d,
	0xfa, 0xad, 0xb4, 0xe1, 0x98, 0xbc, 0xd3, 0x6b, 0xd1, 0x7b, 0x81, 0x46, 0x85, 0x1e, 0x0b, 0x4f,
	0x62, 0x68, 0xac, 0x84, 0x2d, 0x35, 0x2e, 0x02, 0xba, 0x95, 0xbe, 0x47, 0x72, 0x05, 0xf5, 0x22,
	0xb8, 0xc3, 0xf4, 0xfe, 0xe0, 0x64, 0xa7, 0xf3, 0xe3, 0x86, 0xe9, 0xb6, 0x4d, 0x8e, 0x61, 0x2f,
	0x17, 0x3e, 0x5b, 0xb1, 0xf9, 0x52, 0xc4, 0xff, 0xba, 0x51, 0xff, 0x20, 0x6d, 0xe6, 0xc2, 0x3f,
	0x54, 0xb9, 0x37, 0x83, 0xd6, 0xd6, 0x7c, 0x8d, 0xd6, 0x30, 0xf7, 0xc9, 0x12, 0xfd, 0xc9, 0xd2,
	0x86, 0x0a, 0xfa, 0x65, 0xbb, 0x5c, 0xf8, 0x6a, 0xbb, 0xe1, 0x3d, 0x74, 0x38, 0x9a, 0x5d, 0x9c,
	0xf0, 0xb0, 0x93, 0xe8, 0xe9, 0x5c, 0x6a, 0xa7, 0x96, 0x33, 0xca, 0xd1, 0x24, 0x9b, 0xda, 0xf7,
	0xcf, 0xcb, 0x24, 0x66, 0x21, 0xbd, 0xd6, 0xea, 0xd3, 0x9b, 0xdb, 0xf1, 0x64, 0x38, 0xab, 0x87,
	0x7c, 0xf9, 0x16, 0x00, 0x00, 0xff, 0xff, 0x40, 0x5f, 0x69, 0xd2, 0xf4, 0x01, 0x00, 0x00,
}
