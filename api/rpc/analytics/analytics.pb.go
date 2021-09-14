// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: analytics/analytics.proto

package analytics

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Log data request
type LogDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId    string `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
	Operation string `protobuf:"bytes,3,opt,name=operation,proto3" json:"operation,omitempty"`
	Data      string `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *LogDataRequest) Reset() {
	*x = LogDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_analytics_analytics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogDataRequest) ProtoMessage() {}

func (x *LogDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_analytics_analytics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogDataRequest.ProtoReflect.Descriptor instead.
func (*LogDataRequest) Descriptor() ([]byte, []int) {
	return file_analytics_analytics_proto_rawDescGZIP(), []int{0}
}

func (x *LogDataRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *LogDataRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *LogDataRequest) GetOperation() string {
	if x != nil {
		return x.Operation
	}
	return ""
}

func (x *LogDataRequest) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

// Get a log request
type GetLogRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset int32 `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *GetLogRequest) Reset() {
	*x = GetLogRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_analytics_analytics_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLogRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLogRequest) ProtoMessage() {}

func (x *GetLogRequest) ProtoReflect() protoreflect.Message {
	mi := &file_analytics_analytics_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLogRequest.ProtoReflect.Descriptor instead.
func (*GetLogRequest) Descriptor() ([]byte, []int) {
	return file_analytics_analytics_proto_rawDescGZIP(), []int{1}
}

func (x *GetLogRequest) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

// Analytics entry message
type AnalyticsEntryMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId    string `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
	Operation string `protobuf:"bytes,3,opt,name=operation,proto3" json:"operation,omitempty"`
	Data      string `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *AnalyticsEntryMessage) Reset() {
	*x = AnalyticsEntryMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_analytics_analytics_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyticsEntryMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyticsEntryMessage) ProtoMessage() {}

func (x *AnalyticsEntryMessage) ProtoReflect() protoreflect.Message {
	mi := &file_analytics_analytics_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyticsEntryMessage.ProtoReflect.Descriptor instead.
func (*AnalyticsEntryMessage) Descriptor() ([]byte, []int) {
	return file_analytics_analytics_proto_rawDescGZIP(), []int{2}
}

func (x *AnalyticsEntryMessage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AnalyticsEntryMessage) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *AnalyticsEntryMessage) GetOperation() string {
	if x != nil {
		return x.Operation
	}
	return ""
}

func (x *AnalyticsEntryMessage) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File_analytics_analytics_proto protoreflect.FileDescriptor

var file_analytics_analytics_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x74, 0x69, 0x63, 0x73, 0x2f, 0x61, 0x6e, 0x61, 0x6c,
	0x79, 0x74, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x61, 0x6e, 0x61,
	0x6c, 0x79, 0x74, 0x69, 0x63, 0x73, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x6a, 0x0a, 0x0e, 0x4c, 0x6f, 0x67, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x27, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x71, 0x0a, 0x15, 0x41, 0x6e, 0x61, 0x6c,
	0x79, 0x74, 0x69, 0x63, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x93, 0x01, 0x0a, 0x09,
	0x41, 0x6e, 0x61, 0x6c, 0x79, 0x74, 0x69, 0x63, 0x73, 0x12, 0x3e, 0x0a, 0x07, 0x4c, 0x6f, 0x67,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x19, 0x2e, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x74, 0x69, 0x63, 0x73,
	0x2e, 0x4c, 0x6f, 0x67, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x06, 0x47, 0x65, 0x74,
	0x4c, 0x6f, 0x67, 0x12, 0x18, 0x2e, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x74, 0x69, 0x63, 0x73, 0x2e,
	0x47, 0x65, 0x74, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e,
	0x61, 0x6e, 0x61, 0x6c, 0x79, 0x74, 0x69, 0x63, 0x73, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x74,
	0x69, 0x63, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x00, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6d, 0x6f, 0x72, 0x7a, 0x68, 0x61, 0x6e, 0x6f, 0x76, 0x2f, 0x67, 0x6f, 0x2d, 0x72, 0x65, 0x61,
	0x6c, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x61,
	0x6e, 0x61, 0x6c, 0x79, 0x74, 0x69, 0x63, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_analytics_analytics_proto_rawDescOnce sync.Once
	file_analytics_analytics_proto_rawDescData = file_analytics_analytics_proto_rawDesc
)

func file_analytics_analytics_proto_rawDescGZIP() []byte {
	file_analytics_analytics_proto_rawDescOnce.Do(func() {
		file_analytics_analytics_proto_rawDescData = protoimpl.X.CompressGZIP(file_analytics_analytics_proto_rawDescData)
	})
	return file_analytics_analytics_proto_rawDescData
}

var file_analytics_analytics_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_analytics_analytics_proto_goTypes = []interface{}{
	(*LogDataRequest)(nil),        // 0: analytics.LogDataRequest
	(*GetLogRequest)(nil),         // 1: analytics.GetLogRequest
	(*AnalyticsEntryMessage)(nil), // 2: analytics.AnalyticsEntryMessage
	(*emptypb.Empty)(nil),         // 3: google.protobuf.Empty
}
var file_analytics_analytics_proto_depIdxs = []int32{
	0, // 0: analytics.Analytics.LogData:input_type -> analytics.LogDataRequest
	1, // 1: analytics.Analytics.GetLog:input_type -> analytics.GetLogRequest
	3, // 2: analytics.Analytics.LogData:output_type -> google.protobuf.Empty
	2, // 3: analytics.Analytics.GetLog:output_type -> analytics.AnalyticsEntryMessage
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_analytics_analytics_proto_init() }
func file_analytics_analytics_proto_init() {
	if File_analytics_analytics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_analytics_analytics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogDataRequest); i {
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
		file_analytics_analytics_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLogRequest); i {
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
		file_analytics_analytics_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyticsEntryMessage); i {
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
			RawDescriptor: file_analytics_analytics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_analytics_analytics_proto_goTypes,
		DependencyIndexes: file_analytics_analytics_proto_depIdxs,
		MessageInfos:      file_analytics_analytics_proto_msgTypes,
	}.Build()
	File_analytics_analytics_proto = out.File
	file_analytics_analytics_proto_rawDesc = nil
	file_analytics_analytics_proto_goTypes = nil
	file_analytics_analytics_proto_depIdxs = nil
}
