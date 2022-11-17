// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: mediawatch/sessions/v2/session.proto

package sessionsv2

import (
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Session struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" bson:"_id,omitempty"`
	// creation datetime in RFC3339 format
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty" bson:"created_at,omitempty"`
	// update datetime in RFC3339 format
	ExpiresAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty" bson:"expires_at,omitempty"`
	// key to check against
	Key string `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
	// assign value
	Value string `protobuf:"bytes,5,opt,name=value,proto3" json:"value,omitempty"`
	// assign message
	Message string `protobuf:"bytes,6,opt,name=message,proto3" json:"message,omitempty"`
	// client's agent
	Agent string `protobuf:"bytes,7,opt,name=agent,proto3" json:"-"`
	// client's ip
	Ip string `protobuf:"bytes,8,opt,name=ip,proto3" json:"-"`
	// one time password
	Otp string `protobuf:"bytes,9,opt,name=otp,proto3" json:"otp,omitempty" bson:"otp,omitempty"`
}

func (x *Session) Reset() {
	*x = Session{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_sessions_v2_session_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Session) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Session) ProtoMessage() {}

func (x *Session) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_sessions_v2_session_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Session.ProtoReflect.Descriptor instead.
func (*Session) Descriptor() ([]byte, []int) {
	return file_mediawatch_sessions_v2_session_proto_rawDescGZIP(), []int{0}
}

func (x *Session) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Session) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Session) GetExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiresAt
	}
	return nil
}

func (x *Session) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Session) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Session) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Session) GetAgent() string {
	if x != nil {
		return x.Agent
	}
	return ""
}

func (x *Session) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *Session) GetOtp() string {
	if x != nil {
		return x.Otp
	}
	return ""
}

var File_mediawatch_sessions_v2_session_proto protoreflect.FileDescriptor

var file_mediawatch_sessions_v2_session_proto_rawDesc = []byte{
	0x0a, 0x24, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x76, 0x32, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74,
	0x63, 0x68, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x32, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x13, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa1, 0x03, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x29, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x19, 0x9a, 0x84,
	0x9e, 0x03, 0x14, 0x62, 0x73, 0x6f, 0x6e, 0x3a, 0x22, 0x5f, 0x69, 0x64, 0x2c, 0x6f, 0x6d, 0x69,
	0x74, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x52, 0x02, 0x69, 0x64, 0x12, 0x5b, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x20, 0x9a, 0x84, 0x9e,
	0x03, 0x1b, 0x62, 0x73, 0x6f, 0x6e, 0x3a, 0x22, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x2c, 0x6f, 0x6d, 0x69, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x5b, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69,
	0x72, 0x65, 0x73, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x20, 0x9a, 0x84, 0x9e, 0x03, 0x1b, 0x62,
	0x73, 0x6f, 0x6e, 0x3a, 0x22, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x61, 0x74, 0x2c,
	0x6f, 0x6d, 0x69, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69,
	0x72, 0x65, 0x73, 0x41, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x23, 0x0a, 0x05, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0d, 0x9a, 0x84, 0x9e, 0x03, 0x08, 0x6a, 0x73, 0x6f,
	0x6e, 0x3a, 0x22, 0x2d, 0x22, 0x52, 0x05, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x02,
	0x69, 0x70, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0d, 0x9a, 0x84, 0x9e, 0x03, 0x08, 0x6a,
	0x73, 0x6f, 0x6e, 0x3a, 0x22, 0x2d, 0x22, 0x52, 0x02, 0x69, 0x70, 0x12, 0x2b, 0x0a, 0x03, 0x6f,
	0x74, 0x70, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x42, 0x19, 0x9a, 0x84, 0x9e, 0x03, 0x14, 0x62,
	0x73, 0x6f, 0x6e, 0x3a, 0x22, 0x6f, 0x74, 0x70, 0x2c, 0x6f, 0x6d, 0x69, 0x74, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x52, 0x03, 0x6f, 0x74, 0x70, 0x42, 0xec, 0x01, 0x0a, 0x1a, 0x63, 0x6f, 0x6d,
	0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x32, 0x42, 0x0c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x46, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x76, 0x63, 0x69, 0x6f, 0x2f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77,
	0x61, 0x74, 0x63, 0x68, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x65,
	0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x76, 0x32, 0x3b, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x76, 0x32, 0xa2,
	0x02, 0x03, 0x4d, 0x53, 0x58, 0xaa, 0x02, 0x16, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74,
	0x63, 0x68, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x56, 0x32, 0xca, 0x02,
	0x16, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x5c, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x5c, 0x56, 0x32, 0xe2, 0x02, 0x22, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x77,
	0x61, 0x74, 0x63, 0x68, 0x5c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x5c, 0x56, 0x32,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x18, 0x4d,
	0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x3a, 0x3a, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x3a, 0x3a, 0x56, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mediawatch_sessions_v2_session_proto_rawDescOnce sync.Once
	file_mediawatch_sessions_v2_session_proto_rawDescData = file_mediawatch_sessions_v2_session_proto_rawDesc
)

func file_mediawatch_sessions_v2_session_proto_rawDescGZIP() []byte {
	file_mediawatch_sessions_v2_session_proto_rawDescOnce.Do(func() {
		file_mediawatch_sessions_v2_session_proto_rawDescData = protoimpl.X.CompressGZIP(file_mediawatch_sessions_v2_session_proto_rawDescData)
	})
	return file_mediawatch_sessions_v2_session_proto_rawDescData
}

var file_mediawatch_sessions_v2_session_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_mediawatch_sessions_v2_session_proto_goTypes = []interface{}{
	(*Session)(nil),               // 0: mediawatch.sessions.v2.Session
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_mediawatch_sessions_v2_session_proto_depIdxs = []int32{
	1, // 0: mediawatch.sessions.v2.Session.created_at:type_name -> google.protobuf.Timestamp
	1, // 1: mediawatch.sessions.v2.Session.expires_at:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_mediawatch_sessions_v2_session_proto_init() }
func file_mediawatch_sessions_v2_session_proto_init() {
	if File_mediawatch_sessions_v2_session_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mediawatch_sessions_v2_session_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Session); i {
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
			RawDescriptor: file_mediawatch_sessions_v2_session_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_mediawatch_sessions_v2_session_proto_goTypes,
		DependencyIndexes: file_mediawatch_sessions_v2_session_proto_depIdxs,
		MessageInfos:      file_mediawatch_sessions_v2_session_proto_msgTypes,
	}.Build()
	File_mediawatch_sessions_v2_session_proto = out.File
	file_mediawatch_sessions_v2_session_proto_rawDesc = nil
	file_mediawatch_sessions_v2_session_proto_goTypes = nil
	file_mediawatch_sessions_v2_session_proto_depIdxs = nil
}
