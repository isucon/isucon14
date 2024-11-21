// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: isuxportal/services/registration/activate_coupon.proto

package registration

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

type ActivateCouponRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TeamId int64 `protobuf:"varint,1,opt,name=team_id,json=teamId,proto3" json:"team_id,omitempty"`
}

func (x *ActivateCouponRequest) Reset() {
	*x = ActivateCouponRequest{}
	mi := &file_isuxportal_services_registration_activate_coupon_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ActivateCouponRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivateCouponRequest) ProtoMessage() {}

func (x *ActivateCouponRequest) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_registration_activate_coupon_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivateCouponRequest.ProtoReflect.Descriptor instead.
func (*ActivateCouponRequest) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_registration_activate_coupon_proto_rawDescGZIP(), []int{0}
}

func (x *ActivateCouponRequest) GetTeamId() int64 {
	if x != nil {
		return x.TeamId
	}
	return 0
}

type ActivateCouponResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ActivateCouponResponse) Reset() {
	*x = ActivateCouponResponse{}
	mi := &file_isuxportal_services_registration_activate_coupon_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ActivateCouponResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivateCouponResponse) ProtoMessage() {}

func (x *ActivateCouponResponse) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_registration_activate_coupon_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivateCouponResponse.ProtoReflect.Descriptor instead.
func (*ActivateCouponResponse) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_registration_activate_coupon_proto_rawDescGZIP(), []int{1}
}

var File_isuxportal_services_registration_activate_coupon_proto protoreflect.FileDescriptor

var file_isuxportal_services_registration_activate_coupon_proto_rawDesc = []byte{
	0x0a, 0x36, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x70,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x26, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x30, 0x0a, 0x15, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x70,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x65, 0x61,
	0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x74, 0x65, 0x61, 0x6d,
	0x49, 0x64, 0x22, 0x18, 0x0a, 0x16, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x75, 0x70, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0xcd, 0x02, 0x0a,
	0x2a, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x13, 0x41, 0x63, 0x74,
	0x69, 0x76, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x70, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x4e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69,
	0x73, 0x75, 0x63, 0x6f, 0x6e, 0x2f, 0x69, 0x73, 0x75, 0x63, 0x6f, 0x6e, 0x31, 0x34, 0x2f, 0x62,
	0x65, 0x6e, 0x63, 0x68, 0x2f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x72, 0x75, 0x6e, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0xa2, 0x02, 0x04, 0x49, 0x50, 0x53, 0x52, 0xaa, 0x02, 0x26, 0x49, 0x73, 0x75, 0x78,
	0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0xca, 0x02, 0x26, 0x49, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5c,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x5c, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0xe2, 0x02, 0x32, 0x49, 0x73,
	0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x5c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x29, 0x49, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x3a, 0x3a, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x3a, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x3a, 0x3a,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_isuxportal_services_registration_activate_coupon_proto_rawDescOnce sync.Once
	file_isuxportal_services_registration_activate_coupon_proto_rawDescData = file_isuxportal_services_registration_activate_coupon_proto_rawDesc
)

func file_isuxportal_services_registration_activate_coupon_proto_rawDescGZIP() []byte {
	file_isuxportal_services_registration_activate_coupon_proto_rawDescOnce.Do(func() {
		file_isuxportal_services_registration_activate_coupon_proto_rawDescData = protoimpl.X.CompressGZIP(file_isuxportal_services_registration_activate_coupon_proto_rawDescData)
	})
	return file_isuxportal_services_registration_activate_coupon_proto_rawDescData
}

var file_isuxportal_services_registration_activate_coupon_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_isuxportal_services_registration_activate_coupon_proto_goTypes = []any{
	(*ActivateCouponRequest)(nil),  // 0: isuxportal.proto.services.registration.ActivateCouponRequest
	(*ActivateCouponResponse)(nil), // 1: isuxportal.proto.services.registration.ActivateCouponResponse
}
var file_isuxportal_services_registration_activate_coupon_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_isuxportal_services_registration_activate_coupon_proto_init() }
func file_isuxportal_services_registration_activate_coupon_proto_init() {
	if File_isuxportal_services_registration_activate_coupon_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_isuxportal_services_registration_activate_coupon_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_isuxportal_services_registration_activate_coupon_proto_goTypes,
		DependencyIndexes: file_isuxportal_services_registration_activate_coupon_proto_depIdxs,
		MessageInfos:      file_isuxportal_services_registration_activate_coupon_proto_msgTypes,
	}.Build()
	File_isuxportal_services_registration_activate_coupon_proto = out.File
	file_isuxportal_services_registration_activate_coupon_proto_rawDesc = nil
	file_isuxportal_services_registration_activate_coupon_proto_goTypes = nil
	file_isuxportal_services_registration_activate_coupon_proto_depIdxs = nil
}
