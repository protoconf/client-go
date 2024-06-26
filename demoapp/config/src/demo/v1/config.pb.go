// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: demo/v1/config.proto

package democonfig

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

type DemoConfig_LogLevel int32

const (
	DemoConfig_LOG_LEVEL_UNSPECIFIED DemoConfig_LogLevel = 0
	DemoConfig_LOG_LEVEL_DEBUG       DemoConfig_LogLevel = 1
	DemoConfig_LOG_LEVEL_INFO        DemoConfig_LogLevel = 2
	DemoConfig_LOG_LEVEL_WARN        DemoConfig_LogLevel = 3
	DemoConfig_LOG_LEVEL_ERROR       DemoConfig_LogLevel = 4
)

// Enum value maps for DemoConfig_LogLevel.
var (
	DemoConfig_LogLevel_name = map[int32]string{
		0: "LOG_LEVEL_UNSPECIFIED",
		1: "LOG_LEVEL_DEBUG",
		2: "LOG_LEVEL_INFO",
		3: "LOG_LEVEL_WARN",
		4: "LOG_LEVEL_ERROR",
	}
	DemoConfig_LogLevel_value = map[string]int32{
		"LOG_LEVEL_UNSPECIFIED": 0,
		"LOG_LEVEL_DEBUG":       1,
		"LOG_LEVEL_INFO":        2,
		"LOG_LEVEL_WARN":        3,
		"LOG_LEVEL_ERROR":       4,
	}
)

func (x DemoConfig_LogLevel) Enum() *DemoConfig_LogLevel {
	p := new(DemoConfig_LogLevel)
	*p = x
	return p
}

func (x DemoConfig_LogLevel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DemoConfig_LogLevel) Descriptor() protoreflect.EnumDescriptor {
	return file_demo_v1_config_proto_enumTypes[0].Descriptor()
}

func (DemoConfig_LogLevel) Type() protoreflect.EnumType {
	return &file_demo_v1_config_proto_enumTypes[0]
}

func (x DemoConfig_LogLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DemoConfig_LogLevel.Descriptor instead.
func (DemoConfig_LogLevel) EnumDescriptor() ([]byte, []int) {
	return file_demo_v1_config_proto_rawDescGZIP(), []int{0, 0}
}

type DemoConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PodName       string                 `protobuf:"bytes,1,opt,name=pod_name,json=podName,proto3" json:"pod_name,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Version       string                 `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	LogLevel      DemoConfig_LogLevel    `protobuf:"varint,4,opt,name=log_level,json=logLevel,proto3,enum=demoapp.v1.DemoConfig_LogLevel" json:"log_level,omitempty"`
	Timeout       *durationpb.Duration   `protobuf:"bytes,5,opt,name=timeout,proto3" json:"timeout,omitempty"`
	LastUpdate    *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=last_update,json=lastUpdate,proto3" json:"last_update,omitempty"`
	TotalRequests uint64                 `protobuf:"varint,7,opt,name=total_requests,json=totalRequests,proto3" json:"total_requests,omitempty"`
}

func (x *DemoConfig) Reset() {
	*x = DemoConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_demo_v1_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DemoConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DemoConfig) ProtoMessage() {}

func (x *DemoConfig) ProtoReflect() protoreflect.Message {
	mi := &file_demo_v1_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DemoConfig.ProtoReflect.Descriptor instead.
func (*DemoConfig) Descriptor() ([]byte, []int) {
	return file_demo_v1_config_proto_rawDescGZIP(), []int{0}
}

func (x *DemoConfig) GetPodName() string {
	if x != nil {
		return x.PodName
	}
	return ""
}

func (x *DemoConfig) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *DemoConfig) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *DemoConfig) GetLogLevel() DemoConfig_LogLevel {
	if x != nil {
		return x.LogLevel
	}
	return DemoConfig_LOG_LEVEL_UNSPECIFIED
}

func (x *DemoConfig) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *DemoConfig) GetLastUpdate() *timestamppb.Timestamp {
	if x != nil {
		return x.LastUpdate
	}
	return nil
}

func (x *DemoConfig) GetTotalRequests() uint64 {
	if x != nil {
		return x.TotalRequests
	}
	return 0
}

var File_demo_v1_config_proto protoreflect.FileDescriptor

var file_demo_v1_config_proto_rawDesc = []byte{
	0x0a, 0x14, 0x64, 0x65, 0x6d, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x64, 0x65, 0x6d, 0x6f, 0x61, 0x70, 0x70, 0x2e,
	0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xa7, 0x03, 0x0a, 0x0a, 0x44, 0x65, 0x6d, 0x6f, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x6f, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3c, 0x0a,
	0x09, 0x6c, 0x6f, 0x67, 0x5f, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x1f, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65,
	0x6d, 0x6f, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65,
	0x6c, 0x52, 0x08, 0x6c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x33, 0x0a, 0x07, 0x74,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x12, 0x3b, 0x0a, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x25, 0x0a,
	0x0e, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x73, 0x22, 0x77, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c,
	0x12, 0x19, 0x0a, 0x15, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x4c,
	0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x44, 0x45, 0x42, 0x55, 0x47, 0x10, 0x01,
	0x12, 0x12, 0x0a, 0x0e, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x49, 0x4e,
	0x46, 0x4f, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45,
	0x4c, 0x5f, 0x57, 0x41, 0x52, 0x4e, 0x10, 0x03, 0x12, 0x13, 0x0a, 0x0f, 0x4c, 0x4f, 0x47, 0x5f,
	0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x04, 0x42, 0x46, 0x5a,
	0x44, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x67, 0x6f, 0x2f,
	0x64, 0x65, 0x6d, 0x6f, 0x61, 0x70, 0x70, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x73,
	0x72, 0x63, 0x2f, 0x64, 0x65, 0x6d, 0x6f, 0x2f, 0x76, 0x31, 0x3b, 0x64, 0x65, 0x6d, 0x6f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_demo_v1_config_proto_rawDescOnce sync.Once
	file_demo_v1_config_proto_rawDescData = file_demo_v1_config_proto_rawDesc
)

func file_demo_v1_config_proto_rawDescGZIP() []byte {
	file_demo_v1_config_proto_rawDescOnce.Do(func() {
		file_demo_v1_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_demo_v1_config_proto_rawDescData)
	})
	return file_demo_v1_config_proto_rawDescData
}

var file_demo_v1_config_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_demo_v1_config_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_demo_v1_config_proto_goTypes = []any{
	(DemoConfig_LogLevel)(0),      // 0: demoapp.v1.DemoConfig.LogLevel
	(*DemoConfig)(nil),            // 1: demoapp.v1.DemoConfig
	(*durationpb.Duration)(nil),   // 2: google.protobuf.Duration
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_demo_v1_config_proto_depIdxs = []int32{
	0, // 0: demoapp.v1.DemoConfig.log_level:type_name -> demoapp.v1.DemoConfig.LogLevel
	2, // 1: demoapp.v1.DemoConfig.timeout:type_name -> google.protobuf.Duration
	3, // 2: demoapp.v1.DemoConfig.last_update:type_name -> google.protobuf.Timestamp
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_demo_v1_config_proto_init() }
func file_demo_v1_config_proto_init() {
	if File_demo_v1_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_demo_v1_config_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*DemoConfig); i {
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
			RawDescriptor: file_demo_v1_config_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_demo_v1_config_proto_goTypes,
		DependencyIndexes: file_demo_v1_config_proto_depIdxs,
		EnumInfos:         file_demo_v1_config_proto_enumTypes,
		MessageInfos:      file_demo_v1_config_proto_msgTypes,
	}.Build()
	File_demo_v1_config_proto = out.File
	file_demo_v1_config_proto_rawDesc = nil
	file_demo_v1_config_proto_goTypes = nil
	file_demo_v1_config_proto_depIdxs = nil
}
