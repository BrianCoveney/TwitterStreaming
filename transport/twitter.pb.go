// Code generated by protoc-gen-go. DO NOT EDIT.
// source: twitter.proto

package Transport

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Tweet struct {
	Text string `protobuf:"bytes,1,opt,name=Text" json:"Text,omitempty"`
}

func (m *Tweet) Reset()                    { *m = Tweet{} }
func (m *Tweet) String() string            { return proto.CompactTextString(m) }
func (*Tweet) ProtoMessage()               {}
func (*Tweet) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *Tweet) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type Twitter struct {
	TwitterText []string `protobuf:"bytes,1,rep,name=TwitterText" json:"TwitterText,omitempty"`
}

func (m *Twitter) Reset()                    { *m = Twitter{} }
func (m *Twitter) String() string            { return proto.CompactTextString(m) }
func (*Twitter) ProtoMessage()               {}
func (*Twitter) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *Twitter) GetTwitterText() []string {
	if m != nil {
		return m.TwitterText
	}
	return nil
}

func init() {
	proto.RegisterType((*Tweet)(nil), "Transport.Tweet")
	proto.RegisterType((*Twitter)(nil), "Transport.Twitter")
}

func init() { proto.RegisterFile("twitter.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 102 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x29, 0xcf, 0x2c,
	0x29, 0x49, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x0c, 0x29, 0x4a, 0xcc, 0x2b,
	0x2e, 0xc8, 0x2f, 0x2a, 0x51, 0x92, 0xe6, 0x62, 0x0d, 0x29, 0x4f, 0x4d, 0x2d, 0x11, 0x12, 0xe2,
	0x62, 0x09, 0x49, 0xad, 0x28, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x95, 0xb4,
	0xb9, 0xd8, 0x43, 0x20, 0x1a, 0x85, 0x14, 0xb8, 0xb8, 0xa1, 0x4c, 0xa8, 0x2a, 0x66, 0x0d, 0xce,
	0x20, 0x64, 0xa1, 0x24, 0x36, 0xb0, 0xd9, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2b, 0x9a,
	0x27, 0x8d, 0x6c, 0x00, 0x00, 0x00,
}
