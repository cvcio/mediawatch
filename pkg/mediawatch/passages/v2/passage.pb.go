// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: mediawatch/passages/v2/passage.proto

package passagesv2

import (
	v2 "github.com/cvcio/mediawatch/pkg/mediawatch/common/v2"
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

type PassageType int32

const (
	PassageType_TYPE_UNSPECIFIED PassageType = 0
	PassageType_TRIM_LEFT        PassageType = 1
	PassageType_TRIM_RIGHT       PassageType = 2
	PassageType_SPLIT            PassageType = 3
)

// Enum value maps for PassageType.
var (
	PassageType_name = map[int32]string{
		0: "TYPE_UNSPECIFIED",
		1: "TRIM_LEFT",
		2: "TRIM_RIGHT",
		3: "SPLIT",
	}
	PassageType_value = map[string]int32{
		"TYPE_UNSPECIFIED": 0,
		"TRIM_LEFT":        1,
		"TRIM_RIGHT":       2,
		"SPLIT":            3,
	}
)

func (x PassageType) Enum() *PassageType {
	p := new(PassageType)
	*p = x
	return p
}

func (x PassageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PassageType) Descriptor() protoreflect.EnumDescriptor {
	return file_mediawatch_passages_v2_passage_proto_enumTypes[0].Descriptor()
}

func (PassageType) Type() protoreflect.EnumType {
	return &file_mediawatch_passages_v2_passage_proto_enumTypes[0]
}

func (x PassageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PassageType.Descriptor instead.
func (PassageType) EnumDescriptor() ([]byte, []int) {
	return file_mediawatch_passages_v2_passage_proto_rawDescGZIP(), []int{0}
}

type Passage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Text     string      `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	Language string      `protobuf:"bytes,3,opt,name=language,proto3" json:"language,omitempty"`
	Type     PassageType `protobuf:"varint,4,opt,name=type,proto3,enum=mediawatch.passages.v2.PassageType" json:"type,omitempty"`
}

func (x *Passage) Reset() {
	*x = Passage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_passages_v2_passage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Passage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Passage) ProtoMessage() {}

func (x *Passage) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_passages_v2_passage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Passage.ProtoReflect.Descriptor instead.
func (*Passage) Descriptor() ([]byte, []int) {
	return file_mediawatch_passages_v2_passage_proto_rawDescGZIP(), []int{0}
}

func (x *Passage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Passage) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Passage) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

func (x *Passage) GetType() PassageType {
	if x != nil {
		return x.Type
	}
	return PassageType_TYPE_UNSPECIFIED
}

type PassageList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data       []*Passage     `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	Pagination *v2.Pagination `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (x *PassageList) Reset() {
	*x = PassageList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_passages_v2_passage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PassageList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PassageList) ProtoMessage() {}

func (x *PassageList) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_passages_v2_passage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PassageList.ProtoReflect.Descriptor instead.
func (*PassageList) Descriptor() ([]byte, []int) {
	return file_mediawatch_passages_v2_passage_proto_rawDescGZIP(), []int{1}
}

func (x *PassageList) GetData() []*Passage {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *PassageList) GetPagination() *v2.Pagination {
	if x != nil {
		return x.Pagination
	}
	return nil
}

type QueryPassage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Language string `protobuf:"bytes,2,opt,name=language,proto3" json:"language,omitempty"`
}

func (x *QueryPassage) Reset() {
	*x = QueryPassage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_passages_v2_passage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryPassage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryPassage) ProtoMessage() {}

func (x *QueryPassage) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_passages_v2_passage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryPassage.ProtoReflect.Descriptor instead.
func (*QueryPassage) Descriptor() ([]byte, []int) {
	return file_mediawatch_passages_v2_passage_proto_rawDescGZIP(), []int{2}
}

func (x *QueryPassage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *QueryPassage) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

var File_mediawatch_passages_v2_passage_proto protoreflect.FileDescriptor

var file_mediawatch_passages_v2_passage_proto_rawDesc = []byte{
	0x0a, 0x24, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x70, 0x61, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x76, 0x32, 0x2f, 0x70, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74,
	0x63, 0x68, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x76, 0x32, 0x1a, 0x21,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x82, 0x01, 0x0a, 0x07, 0x50, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x37, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e, 0x6d, 0x65,
	0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x2e, 0x76, 0x32, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x84, 0x01, 0x0a, 0x0b, 0x50, 0x61, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63,
	0x68, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x76, 0x32, 0x2e, 0x50, 0x61,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x40, 0x0a, 0x0a, 0x70,
	0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x0a,
	0x0c, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a,
	0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2a, 0x4d, 0x0a, 0x0b, 0x50, 0x61, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0d,
	0x0a, 0x09, 0x54, 0x52, 0x49, 0x4d, 0x5f, 0x4c, 0x45, 0x46, 0x54, 0x10, 0x01, 0x12, 0x0e, 0x0a,
	0x0a, 0x54, 0x52, 0x49, 0x4d, 0x5f, 0x52, 0x49, 0x47, 0x48, 0x54, 0x10, 0x02, 0x12, 0x09, 0x0a,
	0x05, 0x53, 0x50, 0x4c, 0x49, 0x54, 0x10, 0x03, 0x32, 0xc1, 0x01, 0x0a, 0x0e, 0x50, 0x61, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x53, 0x0a, 0x0d, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x2e, 0x6d,
	0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x2e, 0x76, 0x32, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x1f, 0x2e,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x2e, 0x76, 0x32, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00,
	0x12, 0x5a, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x50, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12,
	0x24, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x70, 0x61, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x76, 0x32, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x61,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x23, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74,
	0x63, 0x68, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x76, 0x32, 0x2e, 0x50,
	0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x42, 0xe7, 0x01, 0x0a,
	0x1a, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e,
	0x70, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x76, 0x32, 0x42, 0x0c, 0x50, 0x61, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x41, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x76, 0x63, 0x69, 0x6f, 0x2f, 0x6d, 0x65,
	0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x70, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x2f, 0x76, 0x32, 0x3b, 0x70, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x76, 0x32, 0xa2, 0x02,
	0x03, 0x4d, 0x50, 0x58, 0xaa, 0x02, 0x16, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63,
	0x68, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x56, 0x32, 0xca, 0x02, 0x16,
	0x4d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x5c, 0x50, 0x61, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x5c, 0x56, 0x32, 0xe2, 0x02, 0x22, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61,
	0x74, 0x63, 0x68, 0x5c, 0x50, 0x61, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x5c, 0x56, 0x32, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x18, 0x4d, 0x65,
	0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x3a, 0x3a, 0x50, 0x61, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x3a, 0x3a, 0x56, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mediawatch_passages_v2_passage_proto_rawDescOnce sync.Once
	file_mediawatch_passages_v2_passage_proto_rawDescData = file_mediawatch_passages_v2_passage_proto_rawDesc
)

func file_mediawatch_passages_v2_passage_proto_rawDescGZIP() []byte {
	file_mediawatch_passages_v2_passage_proto_rawDescOnce.Do(func() {
		file_mediawatch_passages_v2_passage_proto_rawDescData = protoimpl.X.CompressGZIP(file_mediawatch_passages_v2_passage_proto_rawDescData)
	})
	return file_mediawatch_passages_v2_passage_proto_rawDescData
}

var file_mediawatch_passages_v2_passage_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_mediawatch_passages_v2_passage_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_mediawatch_passages_v2_passage_proto_goTypes = []interface{}{
	(PassageType)(0),      // 0: mediawatch.passages.v2.PassageType
	(*Passage)(nil),       // 1: mediawatch.passages.v2.Passage
	(*PassageList)(nil),   // 2: mediawatch.passages.v2.PassageList
	(*QueryPassage)(nil),  // 3: mediawatch.passages.v2.QueryPassage
	(*v2.Pagination)(nil), // 4: mediawatch.common.v2.Pagination
}
var file_mediawatch_passages_v2_passage_proto_depIdxs = []int32{
	0, // 0: mediawatch.passages.v2.Passage.type:type_name -> mediawatch.passages.v2.PassageType
	1, // 1: mediawatch.passages.v2.PassageList.data:type_name -> mediawatch.passages.v2.Passage
	4, // 2: mediawatch.passages.v2.PassageList.pagination:type_name -> mediawatch.common.v2.Pagination
	1, // 3: mediawatch.passages.v2.PassageService.CreatePassage:input_type -> mediawatch.passages.v2.Passage
	3, // 4: mediawatch.passages.v2.PassageService.GetPassages:input_type -> mediawatch.passages.v2.QueryPassage
	1, // 5: mediawatch.passages.v2.PassageService.CreatePassage:output_type -> mediawatch.passages.v2.Passage
	2, // 6: mediawatch.passages.v2.PassageService.GetPassages:output_type -> mediawatch.passages.v2.PassageList
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_mediawatch_passages_v2_passage_proto_init() }
func file_mediawatch_passages_v2_passage_proto_init() {
	if File_mediawatch_passages_v2_passage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mediawatch_passages_v2_passage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Passage); i {
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
		file_mediawatch_passages_v2_passage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PassageList); i {
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
		file_mediawatch_passages_v2_passage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryPassage); i {
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
			RawDescriptor: file_mediawatch_passages_v2_passage_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mediawatch_passages_v2_passage_proto_goTypes,
		DependencyIndexes: file_mediawatch_passages_v2_passage_proto_depIdxs,
		EnumInfos:         file_mediawatch_passages_v2_passage_proto_enumTypes,
		MessageInfos:      file_mediawatch_passages_v2_passage_proto_msgTypes,
	}.Build()
	File_mediawatch_passages_v2_passage_proto = out.File
	file_mediawatch_passages_v2_passage_proto_rawDesc = nil
	file_mediawatch_passages_v2_passage_proto_goTypes = nil
	file_mediawatch_passages_v2_passage_proto_depIdxs = nil
}
