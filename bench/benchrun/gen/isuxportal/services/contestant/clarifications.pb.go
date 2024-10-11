// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: isuxportal/services/contestant/clarifications.proto

package contestant

import (
	resources "github.com/isucon/isucon14/bench/benchrun/gen/isuxportal/resources"
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

type ListClarificationsQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListClarificationsQuery) Reset() {
	*x = ListClarificationsQuery{}
	mi := &file_isuxportal_services_contestant_clarifications_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListClarificationsQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListClarificationsQuery) ProtoMessage() {}

func (x *ListClarificationsQuery) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_contestant_clarifications_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListClarificationsQuery.ProtoReflect.Descriptor instead.
func (*ListClarificationsQuery) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_contestant_clarifications_proto_rawDescGZIP(), []int{0}
}

type ListClarificationsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Clarifications []*resources.Clarification `protobuf:"bytes,1,rep,name=clarifications,proto3" json:"clarifications,omitempty"`
}

func (x *ListClarificationsResponse) Reset() {
	*x = ListClarificationsResponse{}
	mi := &file_isuxportal_services_contestant_clarifications_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListClarificationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListClarificationsResponse) ProtoMessage() {}

func (x *ListClarificationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_contestant_clarifications_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListClarificationsResponse.ProtoReflect.Descriptor instead.
func (*ListClarificationsResponse) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_contestant_clarifications_proto_rawDescGZIP(), []int{1}
}

func (x *ListClarificationsResponse) GetClarifications() []*resources.Clarification {
	if x != nil {
		return x.Clarifications
	}
	return nil
}

type RequestClarificationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Question string `protobuf:"bytes,1,opt,name=question,proto3" json:"question,omitempty"`
}

func (x *RequestClarificationRequest) Reset() {
	*x = RequestClarificationRequest{}
	mi := &file_isuxportal_services_contestant_clarifications_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RequestClarificationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestClarificationRequest) ProtoMessage() {}

func (x *RequestClarificationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_contestant_clarifications_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestClarificationRequest.ProtoReflect.Descriptor instead.
func (*RequestClarificationRequest) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_contestant_clarifications_proto_rawDescGZIP(), []int{2}
}

func (x *RequestClarificationRequest) GetQuestion() string {
	if x != nil {
		return x.Question
	}
	return ""
}

type RequestClarificationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Clarification *resources.Clarification `protobuf:"bytes,1,opt,name=clarification,proto3" json:"clarification,omitempty"`
}

func (x *RequestClarificationResponse) Reset() {
	*x = RequestClarificationResponse{}
	mi := &file_isuxportal_services_contestant_clarifications_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RequestClarificationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestClarificationResponse) ProtoMessage() {}

func (x *RequestClarificationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_contestant_clarifications_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestClarificationResponse.ProtoReflect.Descriptor instead.
func (*RequestClarificationResponse) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_contestant_clarifications_proto_rawDescGZIP(), []int{3}
}

func (x *RequestClarificationResponse) GetClarification() *resources.Clarification {
	if x != nil {
		return x.Clarification
	}
	return nil
}

var File_isuxportal_services_contestant_clarifications_proto protoreflect.FileDescriptor

var file_isuxportal_services_contestant_clarifications_proto_rawDesc = []byte{
	0x0a, 0x33, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x2f, 0x63, 0x6c, 0x61, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x24, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61,
	0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x1a, 0x28, 0x69, 0x73, 0x75,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2f, 0x63, 0x6c, 0x61, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x19, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x61,
	0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x22, 0x6f, 0x0a, 0x1a, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x61, 0x72, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51,
	0x0a, 0x0e, 0x63, 0x6c, 0x61, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72,
	0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x43, 0x6c, 0x61, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x0e, 0x63, 0x6c, 0x61, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x22, 0x39, 0x0a, 0x1b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x6c, 0x61, 0x72,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x6f, 0x0a, 0x1c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x6c, 0x61, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4f, 0x0a, 0x0d,
	0x63, 0x6c, 0x61, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x43, 0x6c, 0x61, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d,
	0x63, 0x6c, 0x61, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0xc1, 0x02,
	0x0a, 0x28, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x42, 0x13, 0x43, 0x6c, 0x61, 0x72,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x4c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x73,
	0x75, 0x63, 0x6f, 0x6e, 0x2f, 0x69, 0x73, 0x75, 0x63, 0x6f, 0x6e, 0x31, 0x34, 0x2f, 0x62, 0x65,
	0x6e, 0x63, 0x68, 0x2f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x72, 0x75, 0x6e, 0x2f, 0x67, 0x65, 0x6e,
	0x2f, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0xa2,
	0x02, 0x04, 0x49, 0x50, 0x53, 0x43, 0xaa, 0x02, 0x24, 0x49, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72,
	0x74, 0x61, 0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0xca, 0x02, 0x24,
	0x49, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x5c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x5c, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73,
	0x74, 0x61, 0x6e, 0x74, 0xe2, 0x02, 0x30, 0x49, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61,
	0x6c, 0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x5c, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x27, 0x49, 0x73, 0x75, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x61, 0x6c, 0x3a, 0x3a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x3a, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x3a, 0x3a, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_isuxportal_services_contestant_clarifications_proto_rawDescOnce sync.Once
	file_isuxportal_services_contestant_clarifications_proto_rawDescData = file_isuxportal_services_contestant_clarifications_proto_rawDesc
)

func file_isuxportal_services_contestant_clarifications_proto_rawDescGZIP() []byte {
	file_isuxportal_services_contestant_clarifications_proto_rawDescOnce.Do(func() {
		file_isuxportal_services_contestant_clarifications_proto_rawDescData = protoimpl.X.CompressGZIP(file_isuxportal_services_contestant_clarifications_proto_rawDescData)
	})
	return file_isuxportal_services_contestant_clarifications_proto_rawDescData
}

var file_isuxportal_services_contestant_clarifications_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_isuxportal_services_contestant_clarifications_proto_goTypes = []any{
	(*ListClarificationsQuery)(nil),      // 0: isuxportal.proto.services.contestant.ListClarificationsQuery
	(*ListClarificationsResponse)(nil),   // 1: isuxportal.proto.services.contestant.ListClarificationsResponse
	(*RequestClarificationRequest)(nil),  // 2: isuxportal.proto.services.contestant.RequestClarificationRequest
	(*RequestClarificationResponse)(nil), // 3: isuxportal.proto.services.contestant.RequestClarificationResponse
	(*resources.Clarification)(nil),      // 4: isuxportal.proto.resources.Clarification
}
var file_isuxportal_services_contestant_clarifications_proto_depIdxs = []int32{
	4, // 0: isuxportal.proto.services.contestant.ListClarificationsResponse.clarifications:type_name -> isuxportal.proto.resources.Clarification
	4, // 1: isuxportal.proto.services.contestant.RequestClarificationResponse.clarification:type_name -> isuxportal.proto.resources.Clarification
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_isuxportal_services_contestant_clarifications_proto_init() }
func file_isuxportal_services_contestant_clarifications_proto_init() {
	if File_isuxportal_services_contestant_clarifications_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_isuxportal_services_contestant_clarifications_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_isuxportal_services_contestant_clarifications_proto_goTypes,
		DependencyIndexes: file_isuxportal_services_contestant_clarifications_proto_depIdxs,
		MessageInfos:      file_isuxportal_services_contestant_clarifications_proto_msgTypes,
	}.Build()
	File_isuxportal_services_contestant_clarifications_proto = out.File
	file_isuxportal_services_contestant_clarifications_proto_rawDesc = nil
	file_isuxportal_services_contestant_clarifications_proto_goTypes = nil
	file_isuxportal_services_contestant_clarifications_proto_depIdxs = nil
}
