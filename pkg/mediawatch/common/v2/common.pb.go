// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: mediawatch/common/v2/common.proto

package commonv2

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

// Status Enumeration
type Status int32

const (
	Status_STATUS_UNSPECIFIED Status = 0
	Status_STATUS_PENDING     Status = 1
	Status_STATUS_ACTIVE      Status = 2
	Status_STATUS_SUSPENDED   Status = 3
	Status_STATUS_CLOSED      Status = 4
	Status_STATUS_DELETED     Status = 5
	Status_STATUS_OFFLINE     Status = 6
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0: "STATUS_UNSPECIFIED",
		1: "STATUS_PENDING",
		2: "STATUS_ACTIVE",
		3: "STATUS_SUSPENDED",
		4: "STATUS_CLOSED",
		5: "STATUS_DELETED",
		6: "STATUS_OFFLINE",
	}
	Status_value = map[string]int32{
		"STATUS_UNSPECIFIED": 0,
		"STATUS_PENDING":     1,
		"STATUS_ACTIVE":      2,
		"STATUS_SUSPENDED":   3,
		"STATUS_CLOSED":      4,
		"STATUS_DELETED":     5,
		"STATUS_OFFLINE":     6,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_mediawatch_common_v2_common_proto_enumTypes[0].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_mediawatch_common_v2_common_proto_enumTypes[0]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_mediawatch_common_v2_common_proto_rawDescGZIP(), []int{0}
}

// StreamType Enumeration
type StreamType int32

const (
	StreamType_STREAM_TYPE_UNSPECIFIED StreamType = 0
	StreamType_STREAM_TYPE_OTHER       StreamType = 1
	StreamType_STREAM_TYPE_TWITTER     StreamType = 2
	StreamType_STREAM_TYPE_RSS         StreamType = 3
)

// Enum value maps for StreamType.
var (
	StreamType_name = map[int32]string{
		0: "STREAM_TYPE_UNSPECIFIED",
		1: "STREAM_TYPE_OTHER",
		2: "STREAM_TYPE_TWITTER",
		3: "STREAM_TYPE_RSS",
	}
	StreamType_value = map[string]int32{
		"STREAM_TYPE_UNSPECIFIED": 0,
		"STREAM_TYPE_OTHER":       1,
		"STREAM_TYPE_TWITTER":     2,
		"STREAM_TYPE_RSS":         3,
	}
)

func (x StreamType) Enum() *StreamType {
	p := new(StreamType)
	*p = x
	return p
}

func (x StreamType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StreamType) Descriptor() protoreflect.EnumDescriptor {
	return file_mediawatch_common_v2_common_proto_enumTypes[1].Descriptor()
}

func (StreamType) Type() protoreflect.EnumType {
	return &file_mediawatch_common_v2_common_proto_enumTypes[1]
}

func (x StreamType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StreamType.Descriptor instead.
func (StreamType) EnumDescriptor() ([]byte, []int) {
	return file_mediawatch_common_v2_common_proto_rawDescGZIP(), []int{1}
}

// Locality Enumeration
type Locality int32

const (
	Locality_LOCALITY_UNSPECIFIED   Locality = 0
	Locality_LOCALITY_OTHER         Locality = 1
	Locality_LOCALITY_LOCAL         Locality = 2
	Locality_LOCALITY_NATIONAL      Locality = 3
	Locality_LOCALITY_INTERNATIONAL Locality = 4
	Locality_LOCALITY_MIXED         Locality = 5
)

// Enum value maps for Locality.
var (
	Locality_name = map[int32]string{
		0: "LOCALITY_UNSPECIFIED",
		1: "LOCALITY_OTHER",
		2: "LOCALITY_LOCAL",
		3: "LOCALITY_NATIONAL",
		4: "LOCALITY_INTERNATIONAL",
		5: "LOCALITY_MIXED",
	}
	Locality_value = map[string]int32{
		"LOCALITY_UNSPECIFIED":   0,
		"LOCALITY_OTHER":         1,
		"LOCALITY_LOCAL":         2,
		"LOCALITY_NATIONAL":      3,
		"LOCALITY_INTERNATIONAL": 4,
		"LOCALITY_MIXED":         5,
	}
)

func (x Locality) Enum() *Locality {
	p := new(Locality)
	*p = x
	return p
}

func (x Locality) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Locality) Descriptor() protoreflect.EnumDescriptor {
	return file_mediawatch_common_v2_common_proto_enumTypes[2].Descriptor()
}

func (Locality) Type() protoreflect.EnumType {
	return &file_mediawatch_common_v2_common_proto_enumTypes[2]
}

func (x Locality) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Locality.Descriptor instead.
func (Locality) EnumDescriptor() ([]byte, []int) {
	return file_mediawatch_common_v2_common_proto_rawDescGZIP(), []int{2}
}

// BusinessType Enumeration
type BusinessType int32

const (
	BusinessType_BUSINESS_TYPE_UNSPECIFIED  BusinessType = 0
	BusinessType_BUSINESS_TYPE_OTHER        BusinessType = 1
	BusinessType_BUSINESS_TYPE_AGENCY       BusinessType = 2
	BusinessType_BUSINESS_TYPE_ORGANIZATION BusinessType = 3
	BusinessType_BUSINESS_TYPE_BLOG         BusinessType = 4
	BusinessType_BUSINESS_TYPE_PORTAL       BusinessType = 5
)

// Enum value maps for BusinessType.
var (
	BusinessType_name = map[int32]string{
		0: "BUSINESS_TYPE_UNSPECIFIED",
		1: "BUSINESS_TYPE_OTHER",
		2: "BUSINESS_TYPE_AGENCY",
		3: "BUSINESS_TYPE_ORGANIZATION",
		4: "BUSINESS_TYPE_BLOG",
		5: "BUSINESS_TYPE_PORTAL",
	}
	BusinessType_value = map[string]int32{
		"BUSINESS_TYPE_UNSPECIFIED":  0,
		"BUSINESS_TYPE_OTHER":        1,
		"BUSINESS_TYPE_AGENCY":       2,
		"BUSINESS_TYPE_ORGANIZATION": 3,
		"BUSINESS_TYPE_BLOG":         4,
		"BUSINESS_TYPE_PORTAL":       5,
	}
)

func (x BusinessType) Enum() *BusinessType {
	p := new(BusinessType)
	*p = x
	return p
}

func (x BusinessType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BusinessType) Descriptor() protoreflect.EnumDescriptor {
	return file_mediawatch_common_v2_common_proto_enumTypes[3].Descriptor()
}

func (BusinessType) Type() protoreflect.EnumType {
	return &file_mediawatch_common_v2_common_proto_enumTypes[3]
}

func (x BusinessType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BusinessType.Descriptor instead.
func (BusinessType) EnumDescriptor() ([]byte, []int) {
	return file_mediawatch_common_v2_common_proto_rawDescGZIP(), []int{3}
}

// ContentType Enumeration
type ContentType int32

const (
	ContentType_CONTENT_TYPE_UNSPECIFIED         ContentType = 0
	ContentType_CONTENT_TYPE_OTHER               ContentType = 1
	ContentType_CONTENT_TYPE_NEWS                ContentType = 2
	ContentType_CONTENT_TYPE_MARKET_BUSINESS     ContentType = 3
	ContentType_CONTENT_TYPE_DEFENCE_ARMY_POLICE ContentType = 4
	ContentType_CONTENT_TYPE_ENTERTAINMENT       ContentType = 5
	ContentType_CONTENT_TYPE_HEALTH_BEAUTY       ContentType = 6
	ContentType_CONTENT_TYPE_SPORTS              ContentType = 7
	ContentType_CONTENT_TYPE_RELIGION            ContentType = 8
	ContentType_CONTENT_TYPE_OPINION             ContentType = 9
	ContentType_CONTENT_TYPE_AGRICULTURE         ContentType = 10
	ContentType_CONTENT_TYPE_SCIENCE             ContentType = 11
	ContentType_CONTENT_TYPE_EDUCATION           ContentType = 12
	ContentType_CONTENT_TYPE_JUSTICE             ContentType = 13
)

// Enum value maps for ContentType.
var (
	ContentType_name = map[int32]string{
		0:  "CONTENT_TYPE_UNSPECIFIED",
		1:  "CONTENT_TYPE_OTHER",
		2:  "CONTENT_TYPE_NEWS",
		3:  "CONTENT_TYPE_MARKET_BUSINESS",
		4:  "CONTENT_TYPE_DEFENCE_ARMY_POLICE",
		5:  "CONTENT_TYPE_ENTERTAINMENT",
		6:  "CONTENT_TYPE_HEALTH_BEAUTY",
		7:  "CONTENT_TYPE_SPORTS",
		8:  "CONTENT_TYPE_RELIGION",
		9:  "CONTENT_TYPE_OPINION",
		10: "CONTENT_TYPE_AGRICULTURE",
		11: "CONTENT_TYPE_SCIENCE",
		12: "CONTENT_TYPE_EDUCATION",
		13: "CONTENT_TYPE_JUSTICE",
	}
	ContentType_value = map[string]int32{
		"CONTENT_TYPE_UNSPECIFIED":         0,
		"CONTENT_TYPE_OTHER":               1,
		"CONTENT_TYPE_NEWS":                2,
		"CONTENT_TYPE_MARKET_BUSINESS":     3,
		"CONTENT_TYPE_DEFENCE_ARMY_POLICE": 4,
		"CONTENT_TYPE_ENTERTAINMENT":       5,
		"CONTENT_TYPE_HEALTH_BEAUTY":       6,
		"CONTENT_TYPE_SPORTS":              7,
		"CONTENT_TYPE_RELIGION":            8,
		"CONTENT_TYPE_OPINION":             9,
		"CONTENT_TYPE_AGRICULTURE":         10,
		"CONTENT_TYPE_SCIENCE":             11,
		"CONTENT_TYPE_EDUCATION":           12,
		"CONTENT_TYPE_JUSTICE":             13,
	}
)

func (x ContentType) Enum() *ContentType {
	p := new(ContentType)
	*p = x
	return p
}

func (x ContentType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ContentType) Descriptor() protoreflect.EnumDescriptor {
	return file_mediawatch_common_v2_common_proto_enumTypes[4].Descriptor()
}

func (ContentType) Type() protoreflect.EnumType {
	return &file_mediawatch_common_v2_common_proto_enumTypes[4]
}

func (x ContentType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ContentType.Descriptor instead.
func (ContentType) EnumDescriptor() ([]byte, []int) {
	return file_mediawatch_common_v2_common_proto_rawDescGZIP(), []int{4}
}

// PoliticalOrientation Enumeration
type PoliticalOrientation int32

const (
	PoliticalOrientation_POLITICAL_ORIENTATION_UNSPECIFIED  PoliticalOrientation = 0
	PoliticalOrientation_POLITICAL_ORIENTATION_OTHER        PoliticalOrientation = 1
	PoliticalOrientation_POLITICAL_ORIENTATION_LEFT         PoliticalOrientation = 2
	PoliticalOrientation_POLITICAL_ORIENTATION_CENTER_LEFT  PoliticalOrientation = 3
	PoliticalOrientation_POLITICAL_ORIENTATION_CENTER       PoliticalOrientation = 4
	PoliticalOrientation_POLITICAL_ORIENTATION_CENTER_RIGHT PoliticalOrientation = 5
	PoliticalOrientation_POLITICAL_ORIENTATION_RIGHT        PoliticalOrientation = 6
	PoliticalOrientation_POLITICAL_ORIENTATION_FAR_RIGHT    PoliticalOrientation = 7
)

// Enum value maps for PoliticalOrientation.
var (
	PoliticalOrientation_name = map[int32]string{
		0: "POLITICAL_ORIENTATION_UNSPECIFIED",
		1: "POLITICAL_ORIENTATION_OTHER",
		2: "POLITICAL_ORIENTATION_LEFT",
		3: "POLITICAL_ORIENTATION_CENTER_LEFT",
		4: "POLITICAL_ORIENTATION_CENTER",
		5: "POLITICAL_ORIENTATION_CENTER_RIGHT",
		6: "POLITICAL_ORIENTATION_RIGHT",
		7: "POLITICAL_ORIENTATION_FAR_RIGHT",
	}
	PoliticalOrientation_value = map[string]int32{
		"POLITICAL_ORIENTATION_UNSPECIFIED":  0,
		"POLITICAL_ORIENTATION_OTHER":        1,
		"POLITICAL_ORIENTATION_LEFT":         2,
		"POLITICAL_ORIENTATION_CENTER_LEFT":  3,
		"POLITICAL_ORIENTATION_CENTER":       4,
		"POLITICAL_ORIENTATION_CENTER_RIGHT": 5,
		"POLITICAL_ORIENTATION_RIGHT":        6,
		"POLITICAL_ORIENTATION_FAR_RIGHT":    7,
	}
)

func (x PoliticalOrientation) Enum() *PoliticalOrientation {
	p := new(PoliticalOrientation)
	*p = x
	return p
}

func (x PoliticalOrientation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PoliticalOrientation) Descriptor() protoreflect.EnumDescriptor {
	return file_mediawatch_common_v2_common_proto_enumTypes[5].Descriptor()
}

func (PoliticalOrientation) Type() protoreflect.EnumType {
	return &file_mediawatch_common_v2_common_proto_enumTypes[5]
}

func (x PoliticalOrientation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PoliticalOrientation.Descriptor instead.
func (PoliticalOrientation) EnumDescriptor() ([]byte, []int) {
	return file_mediawatch_common_v2_common_proto_rawDescGZIP(), []int{5}
}

// Tier Enumeration
type Tier int32

const (
	Tier_TIER_UNSPECIFIED        Tier = 0
	Tier_TIER_OTHER              Tier = 1
	Tier_TIER_TRADITIONAL        Tier = 2
	Tier_TIER_DIGITAL            Tier = 3
	Tier_TIER_BROADCASTING_TV    Tier = 4
	Tier_TIER_BROADCASTING_RADIO Tier = 5
	Tier_TIER_MIXED              Tier = 6
)

// Enum value maps for Tier.
var (
	Tier_name = map[int32]string{
		0: "TIER_UNSPECIFIED",
		1: "TIER_OTHER",
		2: "TIER_TRADITIONAL",
		3: "TIER_DIGITAL",
		4: "TIER_BROADCASTING_TV",
		5: "TIER_BROADCASTING_RADIO",
		6: "TIER_MIXED",
	}
	Tier_value = map[string]int32{
		"TIER_UNSPECIFIED":        0,
		"TIER_OTHER":              1,
		"TIER_TRADITIONAL":        2,
		"TIER_DIGITAL":            3,
		"TIER_BROADCASTING_TV":    4,
		"TIER_BROADCASTING_RADIO": 5,
		"TIER_MIXED":              6,
	}
)

func (x Tier) Enum() *Tier {
	p := new(Tier)
	*p = x
	return p
}

func (x Tier) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Tier) Descriptor() protoreflect.EnumDescriptor {
	return file_mediawatch_common_v2_common_proto_enumTypes[6].Descriptor()
}

func (Tier) Type() protoreflect.EnumType {
	return &file_mediawatch_common_v2_common_proto_enumTypes[6]
}

func (x Tier) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Tier.Descriptor instead.
func (Tier) EnumDescriptor() ([]byte, []int) {
	return file_mediawatch_common_v2_common_proto_rawDescGZIP(), []int{6}
}

// Pagination
type Pagination struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total int64 `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Pages int64 `protobuf:"varint,2,opt,name=pages,proto3" json:"pages,omitempty"`
}

func (x *Pagination) Reset() {
	*x = Pagination{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_common_v2_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pagination) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pagination) ProtoMessage() {}

func (x *Pagination) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_common_v2_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pagination.ProtoReflect.Descriptor instead.
func (*Pagination) Descriptor() ([]byte, []int) {
	return file_mediawatch_common_v2_common_proto_rawDescGZIP(), []int{0}
}

func (x *Pagination) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *Pagination) GetPages() int64 {
	if x != nil {
		return x.Pages
	}
	return 0
}

// ResponseWithMessage
type ResponseWithMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ResponseWithMessage) Reset() {
	*x = ResponseWithMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_common_v2_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseWithMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseWithMessage) ProtoMessage() {}

func (x *ResponseWithMessage) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_common_v2_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseWithMessage.ProtoReflect.Descriptor instead.
func (*ResponseWithMessage) Descriptor() ([]byte, []int) {
	return file_mediawatch_common_v2_common_proto_rawDescGZIP(), []int{1}
}

func (x *ResponseWithMessage) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ResponseWithMessage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type SortBy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	By  string `protobuf:"bytes,1,opt,name=by,proto3" json:"by,omitempty"`
	Asc bool   `protobuf:"varint,2,opt,name=asc,proto3" json:"asc,omitempty"`
}

func (x *SortBy) Reset() {
	*x = SortBy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_common_v2_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SortBy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SortBy) ProtoMessage() {}

func (x *SortBy) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_common_v2_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SortBy.ProtoReflect.Descriptor instead.
func (*SortBy) Descriptor() ([]byte, []int) {
	return file_mediawatch_common_v2_common_proto_rawDescGZIP(), []int{2}
}

func (x *SortBy) GetBy() string {
	if x != nil {
		return x.By
	}
	return ""
}

func (x *SortBy) GetAsc() bool {
	if x != nil {
		return x.Asc
	}
	return false
}

type RangeBy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	By   string `protobuf:"bytes,1,opt,name=by,proto3" json:"by,omitempty"`
	From string `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To   string `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
}

func (x *RangeBy) Reset() {
	*x = RangeBy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mediawatch_common_v2_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RangeBy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RangeBy) ProtoMessage() {}

func (x *RangeBy) ProtoReflect() protoreflect.Message {
	mi := &file_mediawatch_common_v2_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RangeBy.ProtoReflect.Descriptor instead.
func (*RangeBy) Descriptor() ([]byte, []int) {
	return file_mediawatch_common_v2_common_proto_rawDescGZIP(), []int{3}
}

func (x *RangeBy) GetBy() string {
	if x != nil {
		return x.By
	}
	return ""
}

func (x *RangeBy) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *RangeBy) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

var File_mediawatch_common_v2_common_proto protoreflect.FileDescriptor

var file_mediawatch_common_v2_common_proto_rawDesc = []byte{
	0x0a, 0x21, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x14, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x22, 0x38, 0x0a, 0x0a, 0x50, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x14, 0x0a,
	0x05, 0x70, 0x61, 0x67, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x70, 0x61,
	0x67, 0x65, 0x73, 0x22, 0x47, 0x0a, 0x13, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x57,
	0x69, 0x74, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x2a, 0x0a, 0x06,
	0x53, 0x6f, 0x72, 0x74, 0x42, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x62, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x62, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x73, 0x63, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x03, 0x61, 0x73, 0x63, 0x22, 0x3d, 0x0a, 0x07, 0x52, 0x61, 0x6e, 0x67,
	0x65, 0x42, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x62, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x62, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x2a, 0x98, 0x01, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x16, 0x0a, 0x12, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x11,
	0x0a, 0x0d, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10,
	0x02, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x53, 0x55, 0x53, 0x50,
	0x45, 0x4e, 0x44, 0x45, 0x44, 0x10, 0x03, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x54, 0x41, 0x54, 0x55,
	0x53, 0x5f, 0x43, 0x4c, 0x4f, 0x53, 0x45, 0x44, 0x10, 0x04, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x05, 0x12, 0x12,
	0x0a, 0x0e, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x4f, 0x46, 0x46, 0x4c, 0x49, 0x4e, 0x45,
	0x10, 0x06, 0x2a, 0x6e, 0x0a, 0x0a, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x1b, 0x0a, 0x17, 0x53, 0x54, 0x52, 0x45, 0x41, 0x4d, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a,
	0x11, 0x53, 0x54, 0x52, 0x45, 0x41, 0x4d, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4f, 0x54, 0x48,
	0x45, 0x52, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x53, 0x54, 0x52, 0x45, 0x41, 0x4d, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x54, 0x57, 0x49, 0x54, 0x54, 0x45, 0x52, 0x10, 0x02, 0x12, 0x13, 0x0a,
	0x0f, 0x53, 0x54, 0x52, 0x45, 0x41, 0x4d, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x53, 0x53,
	0x10, 0x03, 0x2a, 0x93, 0x01, 0x0a, 0x08, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x12,
	0x18, 0x0a, 0x14, 0x4c, 0x4f, 0x43, 0x41, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x4c, 0x4f, 0x43,
	0x41, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x4f, 0x54, 0x48, 0x45, 0x52, 0x10, 0x01, 0x12, 0x12, 0x0a,
	0x0e, 0x4c, 0x4f, 0x43, 0x41, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x4c, 0x4f, 0x43, 0x41, 0x4c, 0x10,
	0x02, 0x12, 0x15, 0x0a, 0x11, 0x4c, 0x4f, 0x43, 0x41, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x41,
	0x54, 0x49, 0x4f, 0x4e, 0x41, 0x4c, 0x10, 0x03, 0x12, 0x1a, 0x0a, 0x16, 0x4c, 0x4f, 0x43, 0x41,
	0x4c, 0x49, 0x54, 0x59, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x41, 0x54, 0x49, 0x4f, 0x4e,
	0x41, 0x4c, 0x10, 0x04, 0x12, 0x12, 0x0a, 0x0e, 0x4c, 0x4f, 0x43, 0x41, 0x4c, 0x49, 0x54, 0x59,
	0x5f, 0x4d, 0x49, 0x58, 0x45, 0x44, 0x10, 0x05, 0x2a, 0xb2, 0x01, 0x0a, 0x0c, 0x42, 0x75, 0x73,
	0x69, 0x6e, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x19, 0x42, 0x55, 0x53,
	0x49, 0x4e, 0x45, 0x53, 0x53, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x42, 0x55, 0x53, 0x49,
	0x4e, 0x45, 0x53, 0x53, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4f, 0x54, 0x48, 0x45, 0x52, 0x10,
	0x01, 0x12, 0x18, 0x0a, 0x14, 0x42, 0x55, 0x53, 0x49, 0x4e, 0x45, 0x53, 0x53, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x41, 0x47, 0x45, 0x4e, 0x43, 0x59, 0x10, 0x02, 0x12, 0x1e, 0x0a, 0x1a, 0x42,
	0x55, 0x53, 0x49, 0x4e, 0x45, 0x53, 0x53, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4f, 0x52, 0x47,
	0x41, 0x4e, 0x49, 0x5a, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x03, 0x12, 0x16, 0x0a, 0x12, 0x42,
	0x55, 0x53, 0x49, 0x4e, 0x45, 0x53, 0x53, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x42, 0x4c, 0x4f,
	0x47, 0x10, 0x04, 0x12, 0x18, 0x0a, 0x14, 0x42, 0x55, 0x53, 0x49, 0x4e, 0x45, 0x53, 0x53, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x50, 0x4f, 0x52, 0x54, 0x41, 0x4c, 0x10, 0x05, 0x2a, 0x9e, 0x03,
	0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a,
	0x18, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x43,
	0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4f, 0x54, 0x48, 0x45,
	0x52, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x4e, 0x45, 0x57, 0x53, 0x10, 0x02, 0x12, 0x20, 0x0a, 0x1c, 0x43, 0x4f,
	0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x41, 0x52, 0x4b, 0x45,
	0x54, 0x5f, 0x42, 0x55, 0x53, 0x49, 0x4e, 0x45, 0x53, 0x53, 0x10, 0x03, 0x12, 0x24, 0x0a, 0x20,
	0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x45, 0x46,
	0x45, 0x4e, 0x43, 0x45, 0x5f, 0x41, 0x52, 0x4d, 0x59, 0x5f, 0x50, 0x4f, 0x4c, 0x49, 0x43, 0x45,
	0x10, 0x04, 0x12, 0x1e, 0x0a, 0x1a, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x54, 0x41, 0x49, 0x4e, 0x4d, 0x45, 0x4e, 0x54,
	0x10, 0x05, 0x12, 0x1e, 0x0a, 0x1a, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x5f, 0x42, 0x45, 0x41, 0x55, 0x54, 0x59,
	0x10, 0x06, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x53, 0x50, 0x4f, 0x52, 0x54, 0x53, 0x10, 0x07, 0x12, 0x19, 0x0a, 0x15, 0x43,
	0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x45, 0x4c, 0x49,
	0x47, 0x49, 0x4f, 0x4e, 0x10, 0x08, 0x12, 0x18, 0x0a, 0x14, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e,
	0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4f, 0x50, 0x49, 0x4e, 0x49, 0x4f, 0x4e, 0x10, 0x09,
	0x12, 0x1c, 0x0a, 0x18, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x41, 0x47, 0x52, 0x49, 0x43, 0x55, 0x4c, 0x54, 0x55, 0x52, 0x45, 0x10, 0x0a, 0x12, 0x18,
	0x0a, 0x14, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53,
	0x43, 0x49, 0x45, 0x4e, 0x43, 0x45, 0x10, 0x0b, 0x12, 0x1a, 0x0a, 0x16, 0x43, 0x4f, 0x4e, 0x54,
	0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x45, 0x44, 0x55, 0x43, 0x41, 0x54, 0x49,
	0x4f, 0x4e, 0x10, 0x0c, 0x12, 0x18, 0x0a, 0x14, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x4a, 0x55, 0x53, 0x54, 0x49, 0x43, 0x45, 0x10, 0x0d, 0x2a, 0xb5,
	0x02, 0x0a, 0x14, 0x50, 0x6f, 0x6c, 0x69, 0x74, 0x69, 0x63, 0x61, 0x6c, 0x4f, 0x72, 0x69, 0x65,
	0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x21, 0x50, 0x4f, 0x4c, 0x49, 0x54,
	0x49, 0x43, 0x41, 0x4c, 0x5f, 0x4f, 0x52, 0x49, 0x45, 0x4e, 0x54, 0x41, 0x54, 0x49, 0x4f, 0x4e,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1f,
	0x0a, 0x1b, 0x50, 0x4f, 0x4c, 0x49, 0x54, 0x49, 0x43, 0x41, 0x4c, 0x5f, 0x4f, 0x52, 0x49, 0x45,
	0x4e, 0x54, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x4f, 0x54, 0x48, 0x45, 0x52, 0x10, 0x01, 0x12,
	0x1e, 0x0a, 0x1a, 0x50, 0x4f, 0x4c, 0x49, 0x54, 0x49, 0x43, 0x41, 0x4c, 0x5f, 0x4f, 0x52, 0x49,
	0x45, 0x4e, 0x54, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x4c, 0x45, 0x46, 0x54, 0x10, 0x02, 0x12,
	0x25, 0x0a, 0x21, 0x50, 0x4f, 0x4c, 0x49, 0x54, 0x49, 0x43, 0x41, 0x4c, 0x5f, 0x4f, 0x52, 0x49,
	0x45, 0x4e, 0x54, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x43, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x5f,
	0x4c, 0x45, 0x46, 0x54, 0x10, 0x03, 0x12, 0x20, 0x0a, 0x1c, 0x50, 0x4f, 0x4c, 0x49, 0x54, 0x49,
	0x43, 0x41, 0x4c, 0x5f, 0x4f, 0x52, 0x49, 0x45, 0x4e, 0x54, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f,
	0x43, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x10, 0x04, 0x12, 0x26, 0x0a, 0x22, 0x50, 0x4f, 0x4c, 0x49,
	0x54, 0x49, 0x43, 0x41, 0x4c, 0x5f, 0x4f, 0x52, 0x49, 0x45, 0x4e, 0x54, 0x41, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x43, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x52, 0x49, 0x47, 0x48, 0x54, 0x10, 0x05,
	0x12, 0x1f, 0x0a, 0x1b, 0x50, 0x4f, 0x4c, 0x49, 0x54, 0x49, 0x43, 0x41, 0x4c, 0x5f, 0x4f, 0x52,
	0x49, 0x45, 0x4e, 0x54, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x52, 0x49, 0x47, 0x48, 0x54, 0x10,
	0x06, 0x12, 0x23, 0x0a, 0x1f, 0x50, 0x4f, 0x4c, 0x49, 0x54, 0x49, 0x43, 0x41, 0x4c, 0x5f, 0x4f,
	0x52, 0x49, 0x45, 0x4e, 0x54, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x41, 0x52, 0x5f, 0x52,
	0x49, 0x47, 0x48, 0x54, 0x10, 0x07, 0x2a, 0x9b, 0x01, 0x0a, 0x04, 0x54, 0x69, 0x65, 0x72, 0x12,
	0x14, 0x0a, 0x10, 0x54, 0x49, 0x45, 0x52, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x54, 0x49, 0x45, 0x52, 0x5f, 0x4f, 0x54,
	0x48, 0x45, 0x52, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x49, 0x45, 0x52, 0x5f, 0x54, 0x52,
	0x41, 0x44, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x41, 0x4c, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x54,
	0x49, 0x45, 0x52, 0x5f, 0x44, 0x49, 0x47, 0x49, 0x54, 0x41, 0x4c, 0x10, 0x03, 0x12, 0x18, 0x0a,
	0x14, 0x54, 0x49, 0x45, 0x52, 0x5f, 0x42, 0x52, 0x4f, 0x41, 0x44, 0x43, 0x41, 0x53, 0x54, 0x49,
	0x4e, 0x47, 0x5f, 0x54, 0x56, 0x10, 0x04, 0x12, 0x1b, 0x0a, 0x17, 0x54, 0x49, 0x45, 0x52, 0x5f,
	0x42, 0x52, 0x4f, 0x41, 0x44, 0x43, 0x41, 0x53, 0x54, 0x49, 0x4e, 0x47, 0x5f, 0x52, 0x41, 0x44,
	0x49, 0x4f, 0x10, 0x05, 0x12, 0x0e, 0x0a, 0x0a, 0x54, 0x49, 0x45, 0x52, 0x5f, 0x4d, 0x49, 0x58,
	0x45, 0x44, 0x10, 0x06, 0x42, 0xd8, 0x01, 0x0a, 0x18, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76,
	0x32, 0x42, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x76, 0x63,
	0x69, 0x6f, 0x2f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x3b, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x76, 0x32, 0xa2,
	0x02, 0x03, 0x4d, 0x43, 0x58, 0xaa, 0x02, 0x14, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74,
	0x63, 0x68, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x56, 0x32, 0xca, 0x02, 0x14, 0x4d,
	0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68, 0x5c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x5c, 0x56, 0x32, 0xe2, 0x02, 0x20, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61, 0x74, 0x63, 0x68,
	0x5c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x56, 0x32, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x16, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x77, 0x61,
	0x74, 0x63, 0x68, 0x3a, 0x3a, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x3a, 0x3a, 0x56, 0x32, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mediawatch_common_v2_common_proto_rawDescOnce sync.Once
	file_mediawatch_common_v2_common_proto_rawDescData = file_mediawatch_common_v2_common_proto_rawDesc
)

func file_mediawatch_common_v2_common_proto_rawDescGZIP() []byte {
	file_mediawatch_common_v2_common_proto_rawDescOnce.Do(func() {
		file_mediawatch_common_v2_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_mediawatch_common_v2_common_proto_rawDescData)
	})
	return file_mediawatch_common_v2_common_proto_rawDescData
}

var file_mediawatch_common_v2_common_proto_enumTypes = make([]protoimpl.EnumInfo, 7)
var file_mediawatch_common_v2_common_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_mediawatch_common_v2_common_proto_goTypes = []interface{}{
	(Status)(0),                 // 0: mediawatch.common.v2.Status
	(StreamType)(0),             // 1: mediawatch.common.v2.StreamType
	(Locality)(0),               // 2: mediawatch.common.v2.Locality
	(BusinessType)(0),           // 3: mediawatch.common.v2.BusinessType
	(ContentType)(0),            // 4: mediawatch.common.v2.ContentType
	(PoliticalOrientation)(0),   // 5: mediawatch.common.v2.PoliticalOrientation
	(Tier)(0),                   // 6: mediawatch.common.v2.Tier
	(*Pagination)(nil),          // 7: mediawatch.common.v2.Pagination
	(*ResponseWithMessage)(nil), // 8: mediawatch.common.v2.ResponseWithMessage
	(*SortBy)(nil),              // 9: mediawatch.common.v2.SortBy
	(*RangeBy)(nil),             // 10: mediawatch.common.v2.RangeBy
}
var file_mediawatch_common_v2_common_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_mediawatch_common_v2_common_proto_init() }
func file_mediawatch_common_v2_common_proto_init() {
	if File_mediawatch_common_v2_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mediawatch_common_v2_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pagination); i {
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
		file_mediawatch_common_v2_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseWithMessage); i {
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
		file_mediawatch_common_v2_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SortBy); i {
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
		file_mediawatch_common_v2_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RangeBy); i {
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
			RawDescriptor: file_mediawatch_common_v2_common_proto_rawDesc,
			NumEnums:      7,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_mediawatch_common_v2_common_proto_goTypes,
		DependencyIndexes: file_mediawatch_common_v2_common_proto_depIdxs,
		EnumInfos:         file_mediawatch_common_v2_common_proto_enumTypes,
		MessageInfos:      file_mediawatch_common_v2_common_proto_msgTypes,
	}.Build()
	File_mediawatch_common_v2_common_proto = out.File
	file_mediawatch_common_v2_common_proto_rawDesc = nil
	file_mediawatch_common_v2_common_proto_goTypes = nil
	file_mediawatch_common_v2_common_proto_depIdxs = nil
}