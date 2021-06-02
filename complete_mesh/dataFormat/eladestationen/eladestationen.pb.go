// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: dataFormat/eladestationen/eladestationen.proto

package eladestationen

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type          string      `protobuf:"bytes,1,opt,name=Type,proto3" json:"Type,omitempty"`
	Totalfeatures int64       `protobuf:"varint,2,opt,name=Totalfeatures,proto3" json:"Totalfeatures,omitempty"`
	Features      []*Features `protobuf:"bytes,3,rep,name=features,proto3" json:"features,omitempty"`
	Crs           *Crs        `protobuf:"bytes,4,opt,name=crs,proto3" json:"crs,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dataFormat_eladestationen_eladestationen_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_dataFormat_eladestationen_eladestationen_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_dataFormat_eladestationen_eladestationen_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Message) GetTotalfeatures() int64 {
	if x != nil {
		return x.Totalfeatures
	}
	return 0
}

func (x *Message) GetFeatures() []*Features {
	if x != nil {
		return x.Features
	}
	return nil
}

func (x *Message) GetCrs() *Crs {
	if x != nil {
		return x.Crs
	}
	return nil
}

type Features struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type         string        `protobuf:"bytes,1,opt,name=Type,proto3" json:"Type,omitempty"`
	ID           string        `protobuf:"bytes,2,opt,name=ID,proto3" json:"ID,omitempty"`
	Geometry     []*Geometry   `protobuf:"bytes,3,rep,name=geometry,proto3" json:"geometry,omitempty"`
	GeometryName string        `protobuf:"bytes,4,opt,name=GeometryName,proto3" json:"GeometryName,omitempty"`
	Properties   []*Properties `protobuf:"bytes,5,rep,name=properties,proto3" json:"properties,omitempty"`
}

func (x *Features) Reset() {
	*x = Features{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dataFormat_eladestationen_eladestationen_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Features) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Features) ProtoMessage() {}

func (x *Features) ProtoReflect() protoreflect.Message {
	mi := &file_dataFormat_eladestationen_eladestationen_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Features.ProtoReflect.Descriptor instead.
func (*Features) Descriptor() ([]byte, []int) {
	return file_dataFormat_eladestationen_eladestationen_proto_rawDescGZIP(), []int{1}
}

func (x *Features) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Features) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Features) GetGeometry() []*Geometry {
	if x != nil {
		return x.Geometry
	}
	return nil
}

func (x *Features) GetGeometryName() string {
	if x != nil {
		return x.GeometryName
	}
	return ""
}

func (x *Features) GetProperties() []*Properties {
	if x != nil {
		return x.Properties
	}
	return nil
}

type Geometry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type        string    `protobuf:"bytes,1,opt,name=Type,proto3" json:"Type,omitempty"`
	Coordinates []float32 `protobuf:"fixed32,2,rep,packed,name=Coordinates,proto3" json:"Coordinates,omitempty"`
}

func (x *Geometry) Reset() {
	*x = Geometry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dataFormat_eladestationen_eladestationen_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Geometry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Geometry) ProtoMessage() {}

func (x *Geometry) ProtoReflect() protoreflect.Message {
	mi := &file_dataFormat_eladestationen_eladestationen_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Geometry.ProtoReflect.Descriptor instead.
func (*Geometry) Descriptor() ([]byte, []int) {
	return file_dataFormat_eladestationen_eladestationen_proto_rawDescGZIP(), []int{2}
}

func (x *Geometry) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Geometry) GetCoordinates() []float32 {
	if x != nil {
		return x.Coordinates
	}
	return nil
}

type Properties struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Objectid          int64      `protobuf:"varint,1,opt,name=Objectid,proto3" json:"Objectid,omitempty"`
	Address           string     `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
	Bezirk            int32      `protobuf:"varint,3,opt,name=Bezirk,proto3" json:"Bezirk,omitempty"`
	City              string     `protobuf:"bytes,5,opt,name=City,proto3" json:"City,omitempty"`
	Countrycode       string     `protobuf:"bytes,6,opt,name=Countrycode,proto3" json:"Countrycode,omitempty"`
	Designation       string     `protobuf:"bytes,7,opt,name=Designation,proto3" json:"Designation,omitempty"`
	Directpayment     int32      `protobuf:"varint,8,opt,name=Directpayment,proto3" json:"Directpayment,omitempty"`
	Evseid            string     `protobuf:"bytes,9,opt,name=Evseid,proto3" json:"Evseid,omitempty"`
	HubjectCompatible int64      `protobuf:"varint,10,opt,name=HubjectCompatible,proto3" json:"HubjectCompatible,omitempty"`
	Operatorname      string     `protobuf:"bytes,11,opt,name=Operatorname,proto3" json:"Operatorname,omitempty"`
	Source            string     `protobuf:"bytes,12,opt,name=Source,proto3" json:"Source,omitempty"`
	SeAnnoCadData     *anypb.Any `protobuf:"bytes,13,opt,name=SeAnnoCadData,proto3" json:"SeAnnoCadData,omitempty"`
}

func (x *Properties) Reset() {
	*x = Properties{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dataFormat_eladestationen_eladestationen_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Properties) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Properties) ProtoMessage() {}

func (x *Properties) ProtoReflect() protoreflect.Message {
	mi := &file_dataFormat_eladestationen_eladestationen_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Properties.ProtoReflect.Descriptor instead.
func (*Properties) Descriptor() ([]byte, []int) {
	return file_dataFormat_eladestationen_eladestationen_proto_rawDescGZIP(), []int{3}
}

func (x *Properties) GetObjectid() int64 {
	if x != nil {
		return x.Objectid
	}
	return 0
}

func (x *Properties) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Properties) GetBezirk() int32 {
	if x != nil {
		return x.Bezirk
	}
	return 0
}

func (x *Properties) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

func (x *Properties) GetCountrycode() string {
	if x != nil {
		return x.Countrycode
	}
	return ""
}

func (x *Properties) GetDesignation() string {
	if x != nil {
		return x.Designation
	}
	return ""
}

func (x *Properties) GetDirectpayment() int32 {
	if x != nil {
		return x.Directpayment
	}
	return 0
}

func (x *Properties) GetEvseid() string {
	if x != nil {
		return x.Evseid
	}
	return ""
}

func (x *Properties) GetHubjectCompatible() int64 {
	if x != nil {
		return x.HubjectCompatible
	}
	return 0
}

func (x *Properties) GetOperatorname() string {
	if x != nil {
		return x.Operatorname
	}
	return ""
}

func (x *Properties) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Properties) GetSeAnnoCadData() *anypb.Any {
	if x != nil {
		return x.SeAnnoCadData
	}
	return nil
}

type Crs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type          string         `protobuf:"bytes,1,opt,name=Type,proto3" json:"Type,omitempty"`
	Crsproperties *CrsProperties `protobuf:"bytes,2,opt,name=crsproperties,proto3" json:"crsproperties,omitempty"`
}

func (x *Crs) Reset() {
	*x = Crs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dataFormat_eladestationen_eladestationen_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Crs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Crs) ProtoMessage() {}

func (x *Crs) ProtoReflect() protoreflect.Message {
	mi := &file_dataFormat_eladestationen_eladestationen_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Crs.ProtoReflect.Descriptor instead.
func (*Crs) Descriptor() ([]byte, []int) {
	return file_dataFormat_eladestationen_eladestationen_proto_rawDescGZIP(), []int{4}
}

func (x *Crs) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Crs) GetCrsproperties() *CrsProperties {
	if x != nil {
		return x.Crsproperties
	}
	return nil
}

type CrsProperties struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
}

func (x *CrsProperties) Reset() {
	*x = CrsProperties{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dataFormat_eladestationen_eladestationen_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CrsProperties) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CrsProperties) ProtoMessage() {}

func (x *CrsProperties) ProtoReflect() protoreflect.Message {
	mi := &file_dataFormat_eladestationen_eladestationen_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CrsProperties.ProtoReflect.Descriptor instead.
func (*CrsProperties) Descriptor() ([]byte, []int) {
	return file_dataFormat_eladestationen_eladestationen_proto_rawDescGZIP(), []int{5}
}

func (x *CrsProperties) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_dataFormat_eladestationen_eladestationen_proto protoreflect.FileDescriptor

var file_dataFormat_eladestationen_eladestationen_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x2f, 0x65, 0x6c, 0x61,
	0x64, 0x65, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x65, 0x6e, 0x2f, 0x65, 0x6c, 0x61, 0x64,
	0x65, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x1a, 0x19, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x98, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x54, 0x6f, 0x74, 0x61, 0x6c,
	0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d,
	0x54, 0x6f, 0x74, 0x61, 0x6c, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12, 0x30, 0x0a,
	0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x2e, 0x46, 0x65, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x73, 0x52, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12,
	0x21, 0x0a, 0x03, 0x63, 0x72, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64,
	0x61, 0x74, 0x61, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x2e, 0x43, 0x72, 0x73, 0x52, 0x03, 0x63,
	0x72, 0x73, 0x22, 0xbc, 0x01, 0x0a, 0x08, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12,
	0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x49, 0x44, 0x12, 0x30, 0x0a, 0x08, 0x67, 0x65, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x46, 0x6f, 0x72, 0x6d,
	0x61, 0x74, 0x2e, 0x47, 0x65, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x52, 0x08, 0x67, 0x65, 0x6f,
	0x6d, 0x65, 0x74, 0x72, 0x79, 0x12, 0x22, 0x0a, 0x0c, 0x47, 0x65, 0x6f, 0x6d, 0x65, 0x74, 0x72,
	0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x47, 0x65, 0x6f,
	0x6d, 0x65, 0x74, 0x72, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x36, 0x0a, 0x0a, 0x70, 0x72, 0x6f,
	0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x65,
	0x72, 0x74, 0x69, 0x65, 0x73, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65,
	0x73, 0x22, 0x40, 0x0a, 0x08, 0x47, 0x65, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x12, 0x12, 0x0a,
	0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x02, 0x52, 0x0b, 0x43, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61,
	0x74, 0x65, 0x73, 0x22, 0x96, 0x03, 0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69,
	0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x42, 0x65, 0x7a, 0x69,
	0x72, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x42, 0x65, 0x7a, 0x69, 0x72, 0x6b,
	0x12, 0x12, 0x0a, 0x04, 0x43, 0x69, 0x74, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x43, 0x69, 0x74, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x72, 0x79, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x65, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x44, 0x65, 0x73,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x44, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0d, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x45, 0x76, 0x73, 0x65, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x45, 0x76, 0x73, 0x65, 0x69, 0x64, 0x12, 0x2c, 0x0a, 0x11, 0x48, 0x75, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x74, 0x69, 0x62, 0x6c, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x11, 0x48, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x74,
	0x69, 0x62, 0x6c, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x4f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x12, 0x3a, 0x0a, 0x0d, 0x53, 0x65, 0x41, 0x6e, 0x6e, 0x6f, 0x43, 0x61, 0x64, 0x44, 0x61, 0x74,
	0x61, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x0d, 0x53,
	0x65, 0x41, 0x6e, 0x6e, 0x6f, 0x43, 0x61, 0x64, 0x44, 0x61, 0x74, 0x61, 0x22, 0x5a, 0x0a, 0x03,
	0x43, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x3f, 0x0a, 0x0d, 0x63, 0x72, 0x73, 0x70, 0x72,
	0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x2e, 0x43, 0x72, 0x73, 0x50,
	0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x52, 0x0d, 0x63, 0x72, 0x73, 0x70, 0x72,
	0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x22, 0x23, 0x0a, 0x0d, 0x43, 0x72, 0x73, 0x50,
	0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x42, 0x02, 0x5a,
	0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dataFormat_eladestationen_eladestationen_proto_rawDescOnce sync.Once
	file_dataFormat_eladestationen_eladestationen_proto_rawDescData = file_dataFormat_eladestationen_eladestationen_proto_rawDesc
)

func file_dataFormat_eladestationen_eladestationen_proto_rawDescGZIP() []byte {
	file_dataFormat_eladestationen_eladestationen_proto_rawDescOnce.Do(func() {
		file_dataFormat_eladestationen_eladestationen_proto_rawDescData = protoimpl.X.CompressGZIP(file_dataFormat_eladestationen_eladestationen_proto_rawDescData)
	})
	return file_dataFormat_eladestationen_eladestationen_proto_rawDescData
}

var file_dataFormat_eladestationen_eladestationen_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_dataFormat_eladestationen_eladestationen_proto_goTypes = []interface{}{
	(*Message)(nil),       // 0: dataFormat.Message
	(*Features)(nil),      // 1: dataFormat.Features
	(*Geometry)(nil),      // 2: dataFormat.Geometry
	(*Properties)(nil),    // 3: dataFormat.Properties
	(*Crs)(nil),           // 4: dataFormat.Crs
	(*CrsProperties)(nil), // 5: dataFormat.CrsProperties
	(*anypb.Any)(nil),     // 6: google.protobuf.Any
}
var file_dataFormat_eladestationen_eladestationen_proto_depIdxs = []int32{
	1, // 0: dataFormat.Message.features:type_name -> dataFormat.Features
	4, // 1: dataFormat.Message.crs:type_name -> dataFormat.Crs
	2, // 2: dataFormat.Features.geometry:type_name -> dataFormat.Geometry
	3, // 3: dataFormat.Features.properties:type_name -> dataFormat.Properties
	6, // 4: dataFormat.Properties.SeAnnoCadData:type_name -> google.protobuf.Any
	5, // 5: dataFormat.Crs.crsproperties:type_name -> dataFormat.CrsProperties
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_dataFormat_eladestationen_eladestationen_proto_init() }
func file_dataFormat_eladestationen_eladestationen_proto_init() {
	if File_dataFormat_eladestationen_eladestationen_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dataFormat_eladestationen_eladestationen_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_dataFormat_eladestationen_eladestationen_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Features); i {
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
		file_dataFormat_eladestationen_eladestationen_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Geometry); i {
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
		file_dataFormat_eladestationen_eladestationen_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Properties); i {
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
		file_dataFormat_eladestationen_eladestationen_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Crs); i {
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
		file_dataFormat_eladestationen_eladestationen_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CrsProperties); i {
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
			RawDescriptor: file_dataFormat_eladestationen_eladestationen_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_dataFormat_eladestationen_eladestationen_proto_goTypes,
		DependencyIndexes: file_dataFormat_eladestationen_eladestationen_proto_depIdxs,
		MessageInfos:      file_dataFormat_eladestationen_eladestationen_proto_msgTypes,
	}.Build()
	File_dataFormat_eladestationen_eladestationen_proto = out.File
	file_dataFormat_eladestationen_eladestationen_proto_rawDesc = nil
	file_dataFormat_eladestationen_eladestationen_proto_goTypes = nil
	file_dataFormat_eladestationen_eladestationen_proto_depIdxs = nil
}
