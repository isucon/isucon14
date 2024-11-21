// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: isuxportal/services/admin/teams.proto

package admin

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

type ListTeamsQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListTeamsQuery) Reset() {
	*x = ListTeamsQuery{}
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListTeamsQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTeamsQuery) ProtoMessage() {}

func (x *ListTeamsQuery) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTeamsQuery.ProtoReflect.Descriptor instead.
func (*ListTeamsQuery) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_admin_teams_proto_rawDescGZIP(), []int{0}
}

type ListTeamsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Teams []*ListTeamsResponse_TeamListItem `protobuf:"bytes,1,rep,name=teams,proto3" json:"teams,omitempty"`
}

func (x *ListTeamsResponse) Reset() {
	*x = ListTeamsResponse{}
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListTeamsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTeamsResponse) ProtoMessage() {}

func (x *ListTeamsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTeamsResponse.ProtoReflect.Descriptor instead.
func (*ListTeamsResponse) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_admin_teams_proto_rawDescGZIP(), []int{1}
}

func (x *ListTeamsResponse) GetTeams() []*ListTeamsResponse_TeamListItem {
	if x != nil {
		return x.Teams
	}
	return nil
}

type GetTeamQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetTeamQuery) Reset() {
	*x = GetTeamQuery{}
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTeamQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTeamQuery) ProtoMessage() {}

func (x *GetTeamQuery) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTeamQuery.ProtoReflect.Descriptor instead.
func (*GetTeamQuery) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_admin_teams_proto_rawDescGZIP(), []int{2}
}

func (x *GetTeamQuery) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetTeamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Team *resources.Team `protobuf:"bytes,1,opt,name=team,proto3" json:"team,omitempty"`
}

func (x *GetTeamResponse) Reset() {
	*x = GetTeamResponse{}
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTeamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTeamResponse) ProtoMessage() {}

func (x *GetTeamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTeamResponse.ProtoReflect.Descriptor instead.
func (*GetTeamResponse) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_admin_teams_proto_rawDescGZIP(), []int{3}
}

func (x *GetTeamResponse) GetTeam() *resources.Team {
	if x != nil {
		return x.Team
	}
	return nil
}

type UpdateTeamQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *UpdateTeamQuery) Reset() {
	*x = UpdateTeamQuery{}
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateTeamQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTeamQuery) ProtoMessage() {}

func (x *UpdateTeamQuery) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTeamQuery.ProtoReflect.Descriptor instead.
func (*UpdateTeamQuery) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_admin_teams_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateTeamQuery) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type UpdateTeamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Team *resources.Team `protobuf:"bytes,1,opt,name=team,proto3" json:"team,omitempty"`
	// Update only specified contestants
	Contestants []*resources.Contestant `protobuf:"bytes,2,rep,name=contestants,proto3" json:"contestants,omitempty"`
}

func (x *UpdateTeamRequest) Reset() {
	*x = UpdateTeamRequest{}
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateTeamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTeamRequest) ProtoMessage() {}

func (x *UpdateTeamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTeamRequest.ProtoReflect.Descriptor instead.
func (*UpdateTeamRequest) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_admin_teams_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateTeamRequest) GetTeam() *resources.Team {
	if x != nil {
		return x.Team
	}
	return nil
}

func (x *UpdateTeamRequest) GetContestants() []*resources.Contestant {
	if x != nil {
		return x.Contestants
	}
	return nil
}

type UpdateTeamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateTeamResponse) Reset() {
	*x = UpdateTeamResponse{}
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateTeamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTeamResponse) ProtoMessage() {}

func (x *UpdateTeamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTeamResponse.ProtoReflect.Descriptor instead.
func (*UpdateTeamResponse) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_admin_teams_proto_rawDescGZIP(), []int{6}
}

type ListTeamsResponse_TeamListItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TeamId             int64    `protobuf:"varint,1,opt,name=team_id,json=teamId,proto3" json:"team_id,omitempty"`
	Name               string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	MemberNames        []string `protobuf:"bytes,3,rep,name=member_names,json=memberNames,proto3" json:"member_names,omitempty"`
	FinalParticipation bool     `protobuf:"varint,4,opt,name=final_participation,json=finalParticipation,proto3" json:"final_participation,omitempty"`
	IsStudent          bool     `protobuf:"varint,5,opt,name=is_student,json=isStudent,proto3" json:"is_student,omitempty"`
	Withdrawn          bool     `protobuf:"varint,6,opt,name=withdrawn,proto3" json:"withdrawn,omitempty"`
	Disqualified       bool     `protobuf:"varint,7,opt,name=disqualified,proto3" json:"disqualified,omitempty"`
	Hidden             bool     `protobuf:"varint,8,opt,name=hidden,proto3" json:"hidden,omitempty"`
}

func (x *ListTeamsResponse_TeamListItem) Reset() {
	*x = ListTeamsResponse_TeamListItem{}
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListTeamsResponse_TeamListItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTeamsResponse_TeamListItem) ProtoMessage() {}

func (x *ListTeamsResponse_TeamListItem) ProtoReflect() protoreflect.Message {
	mi := &file_isuxportal_services_admin_teams_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTeamsResponse_TeamListItem.ProtoReflect.Descriptor instead.
func (*ListTeamsResponse_TeamListItem) Descriptor() ([]byte, []int) {
	return file_isuxportal_services_admin_teams_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ListTeamsResponse_TeamListItem) GetTeamId() int64 {
	if x != nil {
		return x.TeamId
	}
	return 0
}

func (x *ListTeamsResponse_TeamListItem) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListTeamsResponse_TeamListItem) GetMemberNames() []string {
	if x != nil {
		return x.MemberNames
	}
	return nil
}

func (x *ListTeamsResponse_TeamListItem) GetFinalParticipation() bool {
	if x != nil {
		return x.FinalParticipation
	}
	return false
}

func (x *ListTeamsResponse_TeamListItem) GetIsStudent() bool {
	if x != nil {
		return x.IsStudent
	}
	return false
}

func (x *ListTeamsResponse_TeamListItem) GetWithdrawn() bool {
	if x != nil {
		return x.Withdrawn
	}
	return false
}

func (x *ListTeamsResponse_TeamListItem) GetDisqualified() bool {
	if x != nil {
		return x.Disqualified
	}
	return false
}

func (x *ListTeamsResponse_TeamListItem) GetHidden() bool {
	if x != nil {
		return x.Hidden
	}
	return false
}

var File_isuxportal_services_admin_teams_proto protoreflect.FileDescriptor

var file_isuxportal_services_admin_teams_proto_rawDesc = []byte{
	0x0a, 0x25, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x74, 0x65, 0x61, 0x6d,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1f, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72,
	0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x1a, 0x1f, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x61, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74,
	0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25, 0x69, 0x73, 0x75, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x10, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x73, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x22, 0xf5, 0x02, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55, 0x0a, 0x05, 0x74, 0x65, 0x61, 0x6d,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3f, 0x2e, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x65,
	0x61, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x54, 0x65, 0x61, 0x6d,
	0x4c, 0x69, 0x73, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x1a,
	0x88, 0x02, 0x0a, 0x0c, 0x54, 0x65, 0x61, 0x6d, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x74, 0x65, 0x6d,
	0x12, 0x17, 0x0a, 0x07, 0x74, 0x65, 0x61, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x74, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a,
	0x0c, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x73,
	0x12, 0x2f, 0x0a, 0x13, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63,
	0x69, 0x70, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x12, 0x66,
	0x69, 0x6e, 0x61, 0x6c, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74,
	0x12, 0x1c, 0x0a, 0x09, 0x77, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x6e, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x09, 0x77, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x6e, 0x12, 0x22,
	0x0a, 0x0c, 0x64, 0x69, 0x73, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x64, 0x69, 0x73, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x66, 0x69,
	0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x22, 0x1e, 0x0a, 0x0c, 0x47, 0x65,
	0x74, 0x54, 0x65, 0x61, 0x6d, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x47, 0x0a, 0x0f, 0x47, 0x65,
	0x74, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a,
	0x04, 0x74, 0x65, 0x61, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x69, 0x73,
	0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x04, 0x74,
	0x65, 0x61, 0x6d, 0x22, 0x21, 0x0a, 0x0f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x61,
	0x6d, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x93, 0x01, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x04,
	0x74, 0x65, 0x61, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x69, 0x73, 0x75,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x04, 0x74, 0x65,
	0x61, 0x6d, 0x12, 0x48, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x52,
	0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x22, 0x14, 0x0a, 0x12,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x9a, 0x02, 0x0a, 0x23, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x73, 0x75, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x42, 0x0a, 0x54, 0x65, 0x61, 0x6d,
	0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x73, 0x75, 0x63, 0x6f, 0x6e, 0x2f, 0x69, 0x73, 0x75, 0x63,
	0x6f, 0x6e, 0x31, 0x34, 0x2f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x2f, 0x62, 0x65, 0x6e, 0x63, 0x68,
	0x72, 0x75, 0x6e, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x69, 0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74,
	0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0xa2, 0x02, 0x04, 0x49, 0x50, 0x53, 0x41, 0xaa, 0x02, 0x1f, 0x49, 0x73, 0x75, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0xca, 0x02, 0x1f, 0x49, 0x73, 0x75,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x5c, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0xe2, 0x02, 0x2b, 0x49,
	0x73, 0x75, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x5c, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x5c, 0x47,
	0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x22, 0x49, 0x73, 0x75,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x3a, 0x3a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x3a,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x3a, 0x3a, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_isuxportal_services_admin_teams_proto_rawDescOnce sync.Once
	file_isuxportal_services_admin_teams_proto_rawDescData = file_isuxportal_services_admin_teams_proto_rawDesc
)

func file_isuxportal_services_admin_teams_proto_rawDescGZIP() []byte {
	file_isuxportal_services_admin_teams_proto_rawDescOnce.Do(func() {
		file_isuxportal_services_admin_teams_proto_rawDescData = protoimpl.X.CompressGZIP(file_isuxportal_services_admin_teams_proto_rawDescData)
	})
	return file_isuxportal_services_admin_teams_proto_rawDescData
}

var file_isuxportal_services_admin_teams_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_isuxportal_services_admin_teams_proto_goTypes = []any{
	(*ListTeamsQuery)(nil),                 // 0: isuxportal.proto.services.admin.ListTeamsQuery
	(*ListTeamsResponse)(nil),              // 1: isuxportal.proto.services.admin.ListTeamsResponse
	(*GetTeamQuery)(nil),                   // 2: isuxportal.proto.services.admin.GetTeamQuery
	(*GetTeamResponse)(nil),                // 3: isuxportal.proto.services.admin.GetTeamResponse
	(*UpdateTeamQuery)(nil),                // 4: isuxportal.proto.services.admin.UpdateTeamQuery
	(*UpdateTeamRequest)(nil),              // 5: isuxportal.proto.services.admin.UpdateTeamRequest
	(*UpdateTeamResponse)(nil),             // 6: isuxportal.proto.services.admin.UpdateTeamResponse
	(*ListTeamsResponse_TeamListItem)(nil), // 7: isuxportal.proto.services.admin.ListTeamsResponse.TeamListItem
	(*resources.Team)(nil),                 // 8: isuxportal.proto.resources.Team
	(*resources.Contestant)(nil),           // 9: isuxportal.proto.resources.Contestant
}
var file_isuxportal_services_admin_teams_proto_depIdxs = []int32{
	7, // 0: isuxportal.proto.services.admin.ListTeamsResponse.teams:type_name -> isuxportal.proto.services.admin.ListTeamsResponse.TeamListItem
	8, // 1: isuxportal.proto.services.admin.GetTeamResponse.team:type_name -> isuxportal.proto.resources.Team
	8, // 2: isuxportal.proto.services.admin.UpdateTeamRequest.team:type_name -> isuxportal.proto.resources.Team
	9, // 3: isuxportal.proto.services.admin.UpdateTeamRequest.contestants:type_name -> isuxportal.proto.resources.Contestant
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_isuxportal_services_admin_teams_proto_init() }
func file_isuxportal_services_admin_teams_proto_init() {
	if File_isuxportal_services_admin_teams_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_isuxportal_services_admin_teams_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_isuxportal_services_admin_teams_proto_goTypes,
		DependencyIndexes: file_isuxportal_services_admin_teams_proto_depIdxs,
		MessageInfos:      file_isuxportal_services_admin_teams_proto_msgTypes,
	}.Build()
	File_isuxportal_services_admin_teams_proto = out.File
	file_isuxportal_services_admin_teams_proto_rawDesc = nil
	file_isuxportal_services_admin_teams_proto_goTypes = nil
	file_isuxportal_services_admin_teams_proto_depIdxs = nil
}
