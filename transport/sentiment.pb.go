// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sentiment.proto

/*
Package Transport is a generated protocol buffer package.

It is generated from these files:
	sentiment.proto
	twitter.proto
	user.proto

It has these top-level messages:
	Sentiment
	Tweet
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

type Sentiment struct {
	Score int32 `protobuf:"varint,1,opt,name=score" json:"score,omitempty"`
}

func (m *Sentiment) Reset()                    { *m = Sentiment{} }
func (m *Sentiment) String() string            { return proto.CompactTextString(m) }
func (*Sentiment) ProtoMessage()               {}
func (*Sentiment) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Sentiment) GetScore() int32 {
	if m != nil {
		return m.Score
	}
	return 0
}

func init() {
	proto.RegisterType((*Sentiment)(nil), "Transport.Sentiment")
}

func init() { proto.RegisterFile("sentiment.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 82 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0x4e, 0xcd, 0x2b,
	0xc9, 0xcc, 0x4d, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x0c, 0x29, 0x4a,
	0xcc, 0x2b, 0x2e, 0xc8, 0x2f, 0x2a, 0x51, 0x52, 0xe4, 0xe2, 0x0c, 0x86, 0xc9, 0x0a, 0x89, 0x70,
	0xb1, 0x16, 0x27, 0xe7, 0x17, 0xa5, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0xb0, 0x06, 0x41, 0x38, 0x49,
	0x6c, 0x60, 0x4d, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x0d, 0x09, 0xb3, 0x1a, 0x47, 0x00,
	0x00, 0x00,
}