// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.6.1
// source: event.proto

package event

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

type EventType int32

const (
	EventType_ADD_STREAM    EventType = 0
	EventType_REMOVE_STREAM EventType = 1
)

// Enum value maps for EventType.
var (
	EventType_name = map[int32]string{
		0: "ADD_STREAM",
		1: "REMOVE_STREAM",
	}
	EventType_value = map[string]int32{
		"ADD_STREAM":    0,
		"REMOVE_STREAM": 1,
	}
)

func (x EventType) Enum() *EventType {
	p := new(EventType)
	*p = x
	return p
}

func (x EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_event_proto_enumTypes[0].Descriptor()
}

func (EventType) Type() protoreflect.EnumType {
	return &file_event_proto_enumTypes[0]
}

func (x EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventType.Descriptor instead.
func (EventType) EnumDescriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{0}
}

type EventStatus int32

const (
	EventStatus_PLANNED EventStatus = 0
	EventStatus_PASSED  EventStatus = 1
)

// Enum value maps for EventStatus.
var (
	EventStatus_name = map[int32]string{
		0: "PLANNED",
		1: "PASSED",
	}
	EventStatus_value = map[string]int32{
		"PLANNED": 0,
		"PASSED":  1,
	}
)

func (x EventStatus) Enum() *EventStatus {
	p := new(EventStatus)
	*p = x
	return p
}

func (x EventStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_event_proto_enumTypes[1].Descriptor()
}

func (EventStatus) Type() protoreflect.EnumType {
	return &file_event_proto_enumTypes[1]
}

func (x EventStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventStatus.Descriptor instead.
func (EventStatus) EnumDescriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{1}
}

type EventHandler struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Handler string `protobuf:"bytes,1,opt,name=Handler,proto3" json:"Handler,omitempty"`
}

func (x *EventHandler) Reset() {
	*x = EventHandler{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventHandler) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventHandler) ProtoMessage() {}

func (x *EventHandler) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventHandler.ProtoReflect.Descriptor instead.
func (*EventHandler) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{0}
}

func (x *EventHandler) GetHandler() string {
	if x != nil {
		return x.Handler
	}
	return ""
}

type LogInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info string `protobuf:"bytes,1,opt,name=Info,proto3" json:"Info,omitempty"`
}

func (x *LogInfo) Reset() {
	*x = LogInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogInfo) ProtoMessage() {}

func (x *LogInfo) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogInfo.ProtoReflect.Descriptor instead.
func (*LogInfo) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{1}
}

func (x *LogInfo) GetInfo() string {
	if x != nil {
		return x.Info
	}
	return ""
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId       string          `protobuf:"bytes,1,opt,name=EventId,proto3" json:"EventId,omitempty"`
	EventType     EventType       `protobuf:"varint,2,opt,name=EventType,proto3,enum=event.EventType" json:"EventType,omitempty"`
	Status        EventStatus     `protobuf:"varint,3,opt,name=Status,proto3,enum=event.EventStatus" json:"Status,omitempty"`
	Handlers      []*EventHandler `protobuf:"bytes,4,rep,name=Handlers,proto3" json:"Handlers,omitempty"`
	EventGroupId  string          `protobuf:"bytes,5,opt,name=EventGroupId,proto3" json:"EventGroupId,omitempty"`
	ExecutionTime uint32          `protobuf:"varint,6,opt,name=ExecutionTime,proto3" json:"ExecutionTime,omitempty"`
	OccuranceTime uint32          `protobuf:"varint,7,opt,name=OccuranceTime,proto3" json:"OccuranceTime,omitempty"`
	Duration      uint32          `protobuf:"varint,8,opt,name=Duration,proto3" json:"Duration,omitempty"`
	LogInfo       *LogInfo        `protobuf:"bytes,9,opt,name=LogInfo,proto3" json:"LogInfo,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{2}
}

func (x *Event) GetEventId() string {
	if x != nil {
		return x.EventId
	}
	return ""
}

func (x *Event) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_ADD_STREAM
}

func (x *Event) GetStatus() EventStatus {
	if x != nil {
		return x.Status
	}
	return EventStatus_PLANNED
}

func (x *Event) GetHandlers() []*EventHandler {
	if x != nil {
		return x.Handlers
	}
	return nil
}

func (x *Event) GetEventGroupId() string {
	if x != nil {
		return x.EventGroupId
	}
	return ""
}

func (x *Event) GetExecutionTime() uint32 {
	if x != nil {
		return x.ExecutionTime
	}
	return 0
}

func (x *Event) GetOccuranceTime() uint32 {
	if x != nil {
		return x.OccuranceTime
	}
	return 0
}

func (x *Event) GetDuration() uint32 {
	if x != nil {
		return x.Duration
	}
	return 0
}

func (x *Event) GetLogInfo() *LogInfo {
	if x != nil {
		return x.LogInfo
	}
	return nil
}

var File_event_proto protoreflect.FileDescriptor

var file_event_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x22, 0x28, 0x0a, 0x0c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x6e,
	0x64, 0x6c, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x22, 0x1d,
	0x0a, 0x07, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x49, 0x6e, 0x66,
	0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0xe4, 0x02,
	0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x2e, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x2a, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x12, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x0a,
	0x08, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x6e,
	0x64, 0x6c, 0x65, 0x72, 0x52, 0x08, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x73, 0x12, 0x22,
	0x0a, 0x0c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x49, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x54,
	0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x45, 0x78, 0x65, 0x63, 0x75,
	0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x4f, 0x63, 0x63, 0x75,
	0x72, 0x61, 0x6e, 0x63, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0d, 0x4f, 0x63, 0x63, 0x75, 0x72, 0x61, 0x6e, 0x63, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x08, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x07, 0x4c, 0x6f,
	0x67, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x2e, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x4c, 0x6f, 0x67,
	0x49, 0x6e, 0x66, 0x6f, 0x2a, 0x2e, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x41, 0x44, 0x44, 0x5f, 0x53, 0x54, 0x52, 0x45, 0x41, 0x4d, 0x10,
	0x00, 0x12, 0x11, 0x0a, 0x0d, 0x52, 0x45, 0x4d, 0x4f, 0x56, 0x45, 0x5f, 0x53, 0x54, 0x52, 0x45,
	0x41, 0x4d, 0x10, 0x01, 0x2a, 0x26, 0x0a, 0x0b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x4c, 0x41, 0x4e, 0x4e, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x0a, 0x0a, 0x06, 0x50, 0x41, 0x53, 0x53, 0x45, 0x44, 0x10, 0x01, 0x42, 0x08, 0x5a, 0x06,
	0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_event_proto_rawDescOnce sync.Once
	file_event_proto_rawDescData = file_event_proto_rawDesc
)

func file_event_proto_rawDescGZIP() []byte {
	file_event_proto_rawDescOnce.Do(func() {
		file_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_event_proto_rawDescData)
	})
	return file_event_proto_rawDescData
}

var file_event_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_event_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_event_proto_goTypes = []interface{}{
	(EventType)(0),       // 0: event.EventType
	(EventStatus)(0),     // 1: event.EventStatus
	(*EventHandler)(nil), // 2: event.EventHandler
	(*LogInfo)(nil),      // 3: event.LogInfo
	(*Event)(nil),        // 4: event.Event
}
var file_event_proto_depIdxs = []int32{
	0, // 0: event.Event.EventType:type_name -> event.EventType
	1, // 1: event.Event.Status:type_name -> event.EventStatus
	2, // 2: event.Event.Handlers:type_name -> event.EventHandler
	3, // 3: event.Event.LogInfo:type_name -> event.LogInfo
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_event_proto_init() }
func file_event_proto_init() {
	if File_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventHandler); i {
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
		file_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogInfo); i {
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
		file_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
			RawDescriptor: file_event_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_event_proto_goTypes,
		DependencyIndexes: file_event_proto_depIdxs,
		EnumInfos:         file_event_proto_enumTypes,
		MessageInfos:      file_event_proto_msgTypes,
	}.Build()
	File_event_proto = out.File
	file_event_proto_rawDesc = nil
	file_event_proto_goTypes = nil
	file_event_proto_depIdxs = nil
}
