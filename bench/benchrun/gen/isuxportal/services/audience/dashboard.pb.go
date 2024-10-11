// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: isuxportal/services/audience/dashboard.proto

package audience

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

type DashboardQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DashboardQuery) Reset() {
	*x = DashboardQuery{}
	mi := &file_isuxportal_services_audience_dashboard_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DashboardQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DashboardQuery) ProtoMessage() {}

func (x *DashboardQuery) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_audience_dashboard_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DashboardQuery.ProtoReflect.Descriptor instead.
func (*DashboardQuery) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_audience_dashboard_proto_rawDescGZIP(), []int{0}
}

type DashboardResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Leaderboard *resources.Leaderboard `protobuf:"bytes,1,opt,name=leaderboard,proto3" json:"leaderboard,omitempty"`
}

func (x *DashboardResponse) Reset() {
	*x = DashboardResponse{}
	mi := &file_isuxportal_services_audience_dashboard_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DashboardResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DashboardResponse) ProtoMessage() {}

func (x *DashboardResponse) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_audience_dashboard_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DashboardResponse.ProtoReflect.Descriptor instead.
func (*DashboardResponse) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_audience_dashboard_proto_rawDescGZIP(), []int{1}
}

func (x *DashboardResponse) GetLeaderboard() *resources.Leaderboard {
	if x != nil {
		return x.Leaderboard
	}
	return nil
}

type SoloDashboardQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SoloDashboardQuery) Reset() {
	*x = SoloDashboardQuery{}
	mi := &file_isuxportal_services_audience_dashboard_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SoloDashboardQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SoloDashboardQuery) ProtoMessage() {}

func (x *SoloDashboardQuery) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_audience_dashboard_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SoloDashboardQuery.ProtoReflect.Descriptor instead.
func (*SoloDashboardQuery) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_audience_dashboard_proto_rawDescGZIP(), []int{2}
}

type SoloDashboardResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LeaderboardItem *resources.LeaderboardItem `protobuf:"bytes,1,opt,name=leaderboard_item,json=leaderboardItem,proto3" json:"leaderboard_item,omitempty"`
}

func (x *SoloDashboardResponse) Reset() {
	*x = SoloDashboardResponse{}
	mi := &file_isuxportal_services_audience_dashboard_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SoloDashboardResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SoloDashboardResponse) ProtoMessage() {}

func (x *SoloDashboardResponse) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_audience_dashboard_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SoloDashboardResponse.ProtoReflect.Descriptor instead.
func (*SoloDashboardResponse) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_audience_dashboard_proto_rawDescGZIP(), []int{3}
}

func (x *SoloDashboardResponse) GetLeaderboardItem() *resources.LeaderboardItem {
	if x != nil {
		return x.LeaderboardItem
	}
	return nil
}

var File_isuxportal_services_audience_dashboard_proto protoreflect.FileDescriptor

var file_isuxportal_services_audience_dashboard_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x64,
	0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x22,
	0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x65, 0x6e,
	0x63, 0x65, 0x1a, 0x26, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62,
	0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x10, 0x0a, 0x0e, 0x44, 0x61,
	0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x51, 0x75, 0x65, 0x72, 0x79, 0x22, 0x5e, 0x0a, 0x11,
	0x44, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x49, 0x0a, 0x0b, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72,
	0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x52,
	0x0b, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x22, 0x14, 0x0a, 0x12,
	0x53, 0x6f, 0x6c, 0x6f, 0x44, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x22, 0x6f, 0x0a, 0x15, 0x53, 0x6f, 0x6c, 0x6f, 0x44, 0x61, 0x73, 0x68, 0x62, 0x6f,
	0x61, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x56, 0x0a, 0x10, 0x6c,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74,
	0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x49, 0x74,
	0x65, 0x6d, 0x52, 0x0f, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x49,
	0x74, 0x65, 0x6d, 0x42, 0xb0, 0x02, 0x0a, 0x26, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x73, 0x75, 0x78,
	0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x42, 0x0e,
	0x44, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x4a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x73, 0x75,
	0x63, 0x6f, 0x6e, 0x2f, 0x69, 0x73, 0x75, 0x63, 0x6f, 0x6e, 0x31, 0x34, 0x2f, 0x62, 0x65, 0x6e,
	0x63, 0x68, 0x2f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x72, 0x75, 0x6e, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0xa2, 0x02, 0x04, 0x49,
	0x50, 0x53, 0x41, 0xaa, 0x02, 0x22, 0x49, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c,
	0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e,
	0x41, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0xca, 0x02, 0x22, 0x49, 0x73, 0x75, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x5c, 0x41, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0xe2, 0x02, 0x2e,
	0x49, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x5c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x5c, 0x41, 0x75, 0x64, 0x69, 0x65, 0x6e,
	0x63, 0x65, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x25, 0x49, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x3a, 0x3a, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x3a, 0x3a, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x3a, 0x3a, 0x41, 0x75,
	0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_isuxportal_services_audience_dashboard_proto_rawDescOnce sync.Once
	file_isuxportal_services_audience_dashboard_proto_rawDescData = file_isuxportal_services_audience_dashboard_proto_rawDesc
)

func file_isuxportal_services_audience_dashboard_proto_rawDescGZIP() []byte {
	file_isuxportal_services_audience_dashboard_proto_rawDescOnce.Do(func() {
		file_isuxportal_services_audience_dashboard_proto_rawDescData = protoimpl.X.CompressGZIP(file_isuxportal_services_audience_dashboard_proto_rawDescData)
	})
	return file_isuxportal_services_audience_dashboard_proto_rawDescData
}

var file_isuxportal_services_audience_dashboard_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_isuxportal_services_audience_dashboard_proto_goTypes = []any{
	(*DashboardQuery)(nil),            // 0: isuxportal.proto.services.audience.DashboardQuery
	(*DashboardResponse)(nil),         // 1: isuxportal.proto.services.audience.DashboardResponse
	(*SoloDashboardQuery)(nil),        // 2: isuxportal.proto.services.audience.SoloDashboardQuery
	(*SoloDashboardResponse)(nil),     // 3: isuxportal.proto.services.audience.SoloDashboardResponse
	(*resources.Leaderboard)(nil),     // 4: isuxportal.proto.resources.Leaderboard
	(*resources.LeaderboardItem)(nil), // 5: isuxportal.proto.resources.LeaderboardItem
}
var file_isuxportal_services_audience_dashboard_proto_depIdxs = []int32{
	4, // 0: isuxportal.proto.services.audience.DashboardResponse.leaderboard:type_name -> isuxportal.proto.resources.Leaderboard
	5, // 1: isuxportal.proto.services.audience.SoloDashboardResponse.leaderboard_item:type_name -> isuxportal.proto.resources.LeaderboardItem
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_isuxportal_services_audience_dashboard_proto_init() }
func file_isuxportal_services_audience_dashboard_proto_init() {
	if File_isuxportal_services_audience_dashboard_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_isuxportal_services_audience_dashboard_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_isuxportal_services_audience_dashboard_proto_goTypes,
		DependencyIndexes: file_isuxportal_services_audience_dashboard_proto_depIdxs,
		MessageInfos:      file_isuxportal_services_audience_dashboard_proto_msgTypes,
	}.Build()
	File_isuxportal_services_audience_dashboard_proto = out.File
	file_isuxportal_services_audience_dashboard_proto_rawDesc = nil
	file_isuxportal_services_audience_dashboard_proto_goTypes = nil
	file_isuxportal_services_audience_dashboard_proto_depIdxs = nil
}