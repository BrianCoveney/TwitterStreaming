// Code generated by protoc-gen-go. DO NOT EDIT.
// source: twitter.proto

/*
Package Transport is a generated protocol buffer package.

It is generated from these files:
	twitter.proto
	user.proto

It has these top-level messages:
	Text
	User
*/
package Transport

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Text struct {
	Time string `protobuf:"bytes,1,opt,name=time" json:"time,omitempty"`
}

func (m *Text) Reset()                    { *m = Text{} }
func (m *Text) String() string            { return proto.CompactTextString(m) }
func (*Text) ProtoMessage()               {}
func (*Text) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Text) GetTime() string {
	if m != nil {
		return m.Time
	}
	return ""
}

func init() {
	proto.RegisterType((*Text)(nil), "Transport.Text")
}

func init() { proto.RegisterFile("twitter.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 80 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x29, 0xcf, 0x2c,
	0x29, 0x49, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x0c, 0x29, 0x4a, 0xcc, 0x2b,
	0x2e, 0xc8, 0x2f, 0x2a, 0x51, 0x92, 0xe2, 0x62, 0x09, 0x49, 0xad, 0x28, 0x11, 0x12, 0xe2, 0x62,
	0x29, 0xc9, 0xcc, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x93, 0xd8, 0xc0,
	0xaa, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x01, 0xf2, 0x8c, 0x06, 0x3e, 0x00, 0x00, 0x00,
}
