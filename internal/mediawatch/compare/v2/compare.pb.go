// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: mediawatch/compare/v2/compare.proto

package comparev2

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

// SingleRequest
type SingleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// document id
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *SingleRequest) Reset() {
	*x = SingleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SingleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SingleRequest) ProtoMessage() {}

func (x *SingleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SingleRequest.ProtoReflect.Descriptor instead.
func (*SingleRequest) Descriptor() ([]byte, []int) {
	return file_mediawatch_compare_v2_compare_proto_rawDescGZIP(), []int{0}
}

func (x *SingleRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// OneToManyRequest
type OneToManyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// document id
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// list of target document ids
	Targets []string `protobuf:"bytes,2,rep,name=targets,proto3" json:"targets,omitempty"`
}

func (x *OneToManyRequest) Reset() {
	*x = OneToManyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OneToManyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OneToManyRequest) ProtoMessage() {}

func (x *OneToManyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OneToManyRequest.ProtoReflect.Descriptor instead.
func (*OneToManyRequest) Descriptor() ([]byte, []int) {
	return file_mediawatch_compare_v2_compare_proto_rawDescGZIP(), []int{1}
}

func (x *OneToManyRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *OneToManyRequest) GetTargets() []string {
	if x != nil {
		return x.Targets
	}
	return nil
}

// ManyToManyRequest
type ManyToManyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// list of source document ids
	Sources []string `protobuf:"bytes,1,rep,name=sources,proto3" json:"sources,omitempty"`
	// list of target document ids
	Targets []string `protobuf:"bytes,2,rep,name=targets,proto3" json:"targets,omitempty"`
}

func (x *ManyToManyRequest) Reset() {
	*x = ManyToManyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManyToManyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManyToManyRequest) ProtoMessage() {}

func (x *ManyToManyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManyToManyRequest.ProtoReflect.Descriptor instead.
func (*ManyToManyRequest) Descriptor() ([]byte, []int) {
	return file_mediawatch_compare_v2_compare_proto_rawDescGZIP(), []int{2}
}

func (x *ManyToManyRequest) GetSources() []string {
	if x != nil {
		return x.Sources
	}
	return nil
}

func (x *ManyToManyRequest) GetTargets() []string {
	if x != nil {
		return x.Targets
	}
	return nil
}

// SingleResponse
type SingleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// list of results
	Result []*Result `protobuf:"bytes,1,rep,name=result,proto3" json:"result,omitempty"`
}

func (x *SingleResponse) Reset() {
	*x = SingleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SingleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SingleResponse) ProtoMessage() {}

func (x *SingleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SingleResponse.ProtoReflect.Descriptor instead.
func (*SingleResponse) Descriptor() ([]byte, []int) {
	return file_mediawatch_compare_v2_compare_proto_rawDescGZIP(), []int{3}
}

func (x *SingleResponse) GetResult() []*Result {
	if x != nil {
		return x.Result
	}
	return nil
}

// OneToManyResponse
type OneToManyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// list of results
	Result []*Result `protobuf:"bytes,1,rep,name=result,proto3" json:"result,omitempty"`
}

func (x *OneToManyResponse) Reset() {
	*x = OneToManyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OneToManyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OneToManyResponse) ProtoMessage() {}

func (x *OneToManyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OneToManyResponse.ProtoReflect.Descriptor instead.
func (*OneToManyResponse) Descriptor() ([]byte, []int) {
	return file_mediawatch_compare_v2_compare_proto_rawDescGZIP(), []int{4}
}

func (x *OneToManyResponse) GetResult() []*Result {
	if x != nil {
		return x.Result
	}
	return nil
}

// ManyToManyResponse
type ManyToManyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// list of results
	Result []*Result `protobuf:"bytes,1,rep,name=result,proto3" json:"result,omitempty"`
}

func (x *ManyToManyResponse) Reset() {
	*x = ManyToManyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManyToManyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManyToManyResponse) ProtoMessage() {}

func (x *ManyToManyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManyToManyResponse.ProtoReflect.Descriptor instead.
func (*ManyToManyResponse) Descriptor() ([]byte, []int) {
	return file_mediawatch_compare_v2_compare_proto_rawDescGZIP(), []int{5}
}

func (x *ManyToManyResponse) GetResult() []*Result {
	if x != nil {
		return x.Result
	}
	return nil
}

// Result
type Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// success, error
	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	// 200, 500
	Code int32 `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	// message
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	// source document id
	Source string `protobuf:"bytes,4,opt,name=source,proto3" json:"source,omitempty"`
	// target document id
	Target string `protobuf:"bytes,5,opt,name=target,proto3" json:"target,omitempty"`
	// plagiarism score
	Score float64 `protobuf:"fixed64,6,opt,name=score,proto3" json:"score,omitempty"`
}

func (x *Result) Reset() {
	*x = Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_compare_v2_compare_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_mediawatch_compare_v2_compare_proto_rawDescGZIP(), []int{6}
}

func (x *Result) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Result) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Result) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Result) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Result) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *Result) GetScore() float64 {
	if x != nil {
		return x.Score
	}
	return 0
}

var File_mediawatch_compare_v2_compare_proto protoreflect.FileDescriptor

var file_mediawatch_compare_v2_compare_proto_rawDesc = []byte{
	0x0a, 0x23, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x63, 0x6f, 0x6d,
	0x70, 0x61, 0x72, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63,
	0x68, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x32, 0x22, 0x1f, 0x0a, 0x0d,
	0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x3c, 0x0a,
	0x10, 0x4f, 0x6e, 0x65, 0x54, 0x6f, 0x4d, 0x61, 0x6e, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x07, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x22, 0x47, 0x0a, 0x11, 0x4d,
	0x61, 0x6e, 0x79, 0x54, 0x6f, 0x4d, 0x61, 0x6e, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x07, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x73, 0x22, 0x47, 0x0a, 0x0e, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61,
	0x74, 0x63, 0x68, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x32, 0x2e, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x4a, 0x0a,
	0x11, 0x4f, 0x6e, 0x65, 0x54, 0x6f, 0x4d, 0x61, 0x6e, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x35, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e,
	0x63, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x32, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x4b, 0x0a, 0x12, 0x4d, 0x61, 0x6e,
	0x79, 0x54, 0x6f, 0x4d, 0x61, 0x6e, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x35, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x63, 0x6f, 0x6d,
	0x70, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x32, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x94, 0x01, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x32, 0xb0, 0x02,
	0x0a, 0x0e, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x57, 0x0a, 0x06, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x12, 0x24, 0x2e, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x2e,
	0x76, 0x32, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x25, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x63, 0x6f,
	0x6d, 0x70, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x32, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x60, 0x0a, 0x09, 0x4f, 0x6e, 0x65,
	0x54, 0x6f, 0x4d, 0x61, 0x6e, 0x79, 0x12, 0x27, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61,
	0x74, 0x63, 0x68, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x32, 0x2e, 0x4f,
	0x6e, 0x65, 0x54, 0x6f, 0x4d, 0x61, 0x6e, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x28, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x63, 0x6f, 0x6d,
	0x70, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x32, 0x2e, 0x4f, 0x6e, 0x65, 0x54, 0x6f, 0x4d, 0x61, 0x6e,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x63, 0x0a, 0x0a, 0x4d,
	0x61, 0x6e, 0x79, 0x54, 0x6f, 0x4d, 0x61, 0x6e, 0x79, 0x12, 0x28, 0x2e, 0x6d, 0x65, 0x64, 0x69,
	0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x2e, 0x76,
	0x32, 0x2e, 0x4d, 0x61, 0x6e, 0x79, 0x54, 0x6f, 0x4d, 0x61, 0x6e, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68,
	0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x32, 0x2e, 0x4d, 0x61, 0x6e, 0x79,
	0x54, 0x6f, 0x4d, 0x61, 0x6e, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0xe5, 0x01, 0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61,
	0x74, 0x63, 0x68, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x32, 0x42, 0x0c,
	0x43, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x44,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x76, 0x63, 0x69, 0x6f,
	0x2f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2f,
	0x63, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x2f, 0x76, 0x32, 0x3b, 0x63, 0x6f, 0x6d, 0x70, 0x61,
	0x72, 0x65, 0x76, 0x32, 0xa2, 0x02, 0x03, 0x4d, 0x43, 0x58, 0xaa, 0x02, 0x15, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x2e,
	0x56, 0x32, 0xca, 0x02, 0x15, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x5c,
	0x43, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x5c, 0x56, 0x32, 0xe2, 0x02, 0x21, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x5c, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x65, 0x5c,
	0x56, 0x32, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x17, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x3a, 0x3a, 0x43, 0x6f, 0x6d,
	0x70, 0x61, 0x72, 0x65, 0x3a, 0x3a, 0x56, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mediawatch_compare_v2_compare_proto_rawDescOnce sync.Once
	file_mediawatch_compare_v2_compare_proto_rawDescData = file_mediawatch_compare_v2_compare_proto_rawDesc
)

func file_mediawatch_compare_v2_compare_proto_rawDescGZIP() []byte {
	file_mediawatch_compare_v2_compare_proto_rawDescOnce.Do(func() {
		file_mediawatch_compare_v2_compare_proto_rawDescData = protoimpl.X.CompressGZIP(file_mediawatch_compare_v2_compare_proto_rawDescData)
	})
	return file_mediawatch_compare_v2_compare_proto_rawDescData
}

var file_mediawatch_compare_v2_compare_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_mediawatch_compare_v2_compare_proto_goTypes = []interface{}{
	(*SingleRequest)(nil),      // 0: mediawatch.compare.v2.SingleRequest
	(*OneToManyRequest)(nil),   // 1: mediawatch.compare.v2.OneToManyRequest
	(*ManyToManyRequest)(nil),  // 2: mediawatch.compare.v2.ManyToManyRequest
	(*SingleResponse)(nil),     // 3: mediawatch.compare.v2.SingleResponse
	(*OneToManyResponse)(nil),  // 4: mediawatch.compare.v2.OneToManyResponse
	(*ManyToManyResponse)(nil), // 5: mediawatch.compare.v2.ManyToManyResponse
	(*Result)(nil),             // 6: mediawatch.compare.v2.Result
}
var file_mediawatch_compare_v2_compare_proto_depIdxs = []int32{
	6, // 0: mediawatch.compare.v2.SingleResponse.result:type_name -> mediawatch.compare.v2.Result
	6, // 1: mediawatch.compare.v2.OneToManyResponse.result:type_name -> mediawatch.compare.v2.Result
	6, // 2: mediawatch.compare.v2.ManyToManyResponse.result:type_name -> mediawatch.compare.v2.Result
	0, // 3: mediawatch.compare.v2.CompareService.Single:input_type -> mediawatch.compare.v2.SingleRequest
	1, // 4: mediawatch.compare.v2.CompareService.OneToMany:input_type -> mediawatch.compare.v2.OneToManyRequest
	2, // 5: mediawatch.compare.v2.CompareService.ManyToMany:input_type -> mediawatch.compare.v2.ManyToManyRequest
	3, // 6: mediawatch.compare.v2.CompareService.Single:output_type -> mediawatch.compare.v2.SingleResponse
	4, // 7: mediawatch.compare.v2.CompareService.OneToMany:output_type -> mediawatch.compare.v2.OneToManyResponse
	5, // 8: mediawatch.compare.v2.CompareService.ManyToMany:output_type -> mediawatch.compare.v2.ManyToManyResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_mediawatch_compare_v2_compare_proto_init() }
func file_mediawatch_compare_v2_compare_proto_init() {
	if File_mediawatch_compare_v2_compare_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mediawatch_compare_v2_compare_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SingleRequest); i {
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
		file_mediawatch_compare_v2_compare_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OneToManyRequest); i {
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
		file_mediawatch_compare_v2_compare_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ManyToManyRequest); i {
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
		file_mediawatch_compare_v2_compare_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SingleResponse); i {
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
		file_mediawatch_compare_v2_compare_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OneToManyResponse); i {
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
		file_mediawatch_compare_v2_compare_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ManyToManyResponse); i {
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
		file_mediawatch_compare_v2_compare_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Result); i {
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
			RawDescriptor: file_mediawatch_compare_v2_compare_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mediawatch_compare_v2_compare_proto_goTypes,
		DependencyIndexes: file_mediawatch_compare_v2_compare_proto_depIdxs,
		MessageInfos:      file_mediawatch_compare_v2_compare_proto_msgTypes,
	}.Build()
	File_mediawatch_compare_v2_compare_proto = out.File
	file_mediawatch_compare_v2_compare_proto_rawDesc = nil
	file_mediawatch_compare_v2_compare_proto_goTypes = nil
	file_mediawatch_compare_v2_compare_proto_depIdxs = nil
}
