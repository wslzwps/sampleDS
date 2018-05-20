// Code generated by protoc-gen-go. DO NOT EDIT.
// source: raftlog.proto

/*
Package raftlog_pb is a generated protocol buffer package.

It is generated from these files:
	raftlog.proto

It has these top-level messages:
	LBO
*/
package raftlogpb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import raftpb "github.com/coreos/etcd/raft/raftpb"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
//const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type LBO struct {
	Offset           *int64        `protobuf:"varint,1,req,name=offset" json:"offset,omitempty"`
	Entry            *raftpb.Entry `protobuf:"bytes,2,req,name=entry" json:"entry,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *LBO) Reset()                    { *m = LBO{} }
func (m *LBO) String() string            { return proto.CompactTextString(m) }
func (*LBO) ProtoMessage()               {}
func (*LBO) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *LBO) GetOffset() int64 {
	if m != nil && m.Offset != nil {
		return *m.Offset
	}
	return 0
}

func (m *LBO) GetEntry() *raftpb.Entry {
	if m != nil {
		return m.Entry
	}
	return nil
}

func init() {
	proto.RegisterType((*LBO)(nil), "raftlog_pb.LBO")
}

func init() { proto.RegisterFile("raftlog.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 105 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4a, 0x4c, 0x2b,
	0xc9, 0xc9, 0x4f, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x82, 0x72, 0xe3, 0x0b, 0x92,
	0xa4, 0xc0, 0x6c, 0x88, 0xb8, 0x92, 0x13, 0x17, 0xb3, 0x8f, 0x93, 0xbf, 0x90, 0x18, 0x17, 0x5b,
	0x7e, 0x5a, 0x5a, 0x71, 0x6a, 0x89, 0x04, 0xa3, 0x02, 0x93, 0x06, 0x73, 0x10, 0x94, 0x27, 0xa4,
	0xcc, 0xc5, 0x9a, 0x9a, 0x57, 0x52, 0x54, 0x29, 0xc1, 0xa4, 0xc0, 0xa4, 0xc1, 0x6d, 0xc4, 0xab,
	0x07, 0xd2, 0x5a, 0x90, 0xa4, 0xe7, 0x0a, 0x12, 0x0c, 0x82, 0xc8, 0x01, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x89, 0xbf, 0x8e, 0x79, 0x6b, 0x00, 0x00, 0x00,
}
