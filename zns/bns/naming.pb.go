// Code generated by protoc-gen-go.
// source: naming.proto
// DO NOT EDIT!

/*
Package bns is a generated protocol buffer package.

It is generated from these files:

	naming.proto
	naminglib.proto
	service.proto

It has these top-level messages:

	BnsKvPair
	BnsInput
	BnsOutput
	BnsInstance
*/
package bns

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

// *
//
//	fill this struct to ask BNS for a instance list
type BnsKvPair struct {
	Key              *string `protobuf:"bytes,1,req,name=key" json:"key,omitempty"`
	Value            *string `protobuf:"bytes,2,req,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *BnsKvPair) Reset()         { *m = BnsKvPair{} }
func (m *BnsKvPair) String() string { return proto.CompactTextString(m) }
func (*BnsKvPair) ProtoMessage()    {}

func (m *BnsKvPair) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *BnsKvPair) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

type BnsInput struct {
	ServiceName *string `protobuf:"bytes,1,req,name=service_name" json:"service_name,omitempty"`
	TimeoutMs   *uint32 `protobuf:"varint,2,opt,name=timeout_ms,def=1500" json:"timeout_ms,omitempty"`
	Type        *uint32 `protobuf:"varint,3,opt,name=type,def=0" json:"type,omitempty"`
	// type=1 use for ip white list
	TagConstrain     []*BnsKvPair `protobuf:"bytes,4,rep,name=tag_constrain" json:"tag_constrain,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *BnsInput) Reset()         { *m = BnsInput{} }
func (m *BnsInput) String() string { return proto.CompactTextString(m) }
func (*BnsInput) ProtoMessage()    {}

const Default_BnsInput_TimeoutMs uint32 = 1500
const Default_BnsInput_Type uint32 = 0

func (m *BnsInput) GetServiceName() string {
	if m != nil && m.ServiceName != nil {
		return *m.ServiceName
	}
	return ""
}

func (m *BnsInput) GetTimeoutMs() uint32 {
	if m != nil && m.TimeoutMs != nil {
		return *m.TimeoutMs
	}
	return Default_BnsInput_TimeoutMs
}

func (m *BnsInput) GetType() uint32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Default_BnsInput_Type
}

func (m *BnsInput) GetTagConstrain() []*BnsKvPair {
	if m != nil {
		return m.TagConstrain
	}
	return nil
}

type BnsOutput struct {
	ServiceName      *string        `protobuf:"bytes,1,req,name=service_name" json:"service_name,omitempty"`
	Instance         []*BnsInstance `protobuf:"bytes,2,rep,name=instance" json:"instance,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *BnsOutput) Reset()         { *m = BnsOutput{} }
func (m *BnsOutput) String() string { return proto.CompactTextString(m) }
func (*BnsOutput) ProtoMessage()    {}

func (m *BnsOutput) GetServiceName() string {
	if m != nil && m.ServiceName != nil {
		return *m.ServiceName
	}
	return ""
}

func (m *BnsOutput) GetInstance() []*BnsInstance {
	if m != nil {
		return m.Instance
	}
	return nil
}

type BnsInstance struct {
	ServiceName      *string      `protobuf:"bytes,1,req,name=service_name" json:"service_name,omitempty"`
	HostName         *string      `protobuf:"bytes,2,req,name=host_name" json:"host_name,omitempty"`
	HostIp           *string      `protobuf:"bytes,3,req,name=host_ip" json:"host_ip,omitempty"`
	HostIpUint       *uint32      `protobuf:"varint,4,req,name=host_ip_uint" json:"host_ip_uint,omitempty"`
	Status           *int32       `protobuf:"varint,5,opt,name=status,def=0" json:"status,omitempty"`
	Port             *int32       `protobuf:"varint,6,opt,name=port,def=0" json:"port,omitempty"`
	Tag              *string      `protobuf:"bytes,7,opt,name=tag,def=" json:"tag,omitempty"`
	Load             *int32       `protobuf:"varint,8,opt,name=load,def=-1" json:"load,omitempty"`
	Offset           *int32       `protobuf:"varint,9,opt,name=offset,def=0" json:"offset,omitempty"`
	Extra            *string      `protobuf:"bytes,10,opt,name=extra,def=" json:"extra,omitempty"`
	MultiPort        *string      `protobuf:"bytes,11,opt,name=multi_port,def=" json:"multi_port,omitempty"`
	TagKvPairFormat  []*BnsKvPair `protobuf:"bytes,12,rep,name=tag_kv_pair_format" json:"tag_kv_pair_format,omitempty"`
	ContainerId      *string      `protobuf:"bytes,13,opt,name=container_id" json:"container_id,omitempty"`
	DeployPath       *string      `protobuf:"bytes,14,opt,name=deploy_path" json:"deploy_path,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *BnsInstance) Reset()         { *m = BnsInstance{} }
func (m *BnsInstance) String() string { return proto.CompactTextString(m) }
func (*BnsInstance) ProtoMessage()    {}

const Default_BnsInstance_Status int32 = 0
const Default_BnsInstance_Port int32 = 0
const Default_BnsInstance_Load int32 = -1
const Default_BnsInstance_Offset int32 = 0

func (m *BnsInstance) GetServiceName() string {
	if m != nil && m.ServiceName != nil {
		return *m.ServiceName
	}
	return ""
}

func (m *BnsInstance) GetHostName() string {
	if m != nil && m.HostName != nil {
		return *m.HostName
	}
	return ""
}

func (m *BnsInstance) GetHostIp() string {
	if m != nil && m.HostIp != nil {
		return *m.HostIp
	}
	return ""
}

func (m *BnsInstance) GetHostIpUint() uint32 {
	if m != nil && m.HostIpUint != nil {
		return *m.HostIpUint
	}
	return 0
}

func (m *BnsInstance) GetStatus() int32 {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return Default_BnsInstance_Status
}

func (m *BnsInstance) GetPort() int32 {
	if m != nil && m.Port != nil {
		return *m.Port
	}
	return Default_BnsInstance_Port
}

func (m *BnsInstance) GetTag() string {
	if m != nil && m.Tag != nil {
		return *m.Tag
	}
	return ""
}

func (m *BnsInstance) GetLoad() int32 {
	if m != nil && m.Load != nil {
		return *m.Load
	}
	return Default_BnsInstance_Load
}

func (m *BnsInstance) GetOffset() int32 {
	if m != nil && m.Offset != nil {
		return *m.Offset
	}
	return Default_BnsInstance_Offset
}

func (m *BnsInstance) GetExtra() string {
	if m != nil && m.Extra != nil {
		return *m.Extra
	}
	return ""
}

func (m *BnsInstance) GetMultiPort() string {
	if m != nil && m.MultiPort != nil {
		return *m.MultiPort
	}
	return ""
}

func (m *BnsInstance) GetTagKvPairFormat() []*BnsKvPair {
	if m != nil {
		return m.TagKvPairFormat
	}
	return nil
}

func (m *BnsInstance) GetContainerId() string {
	if m != nil && m.ContainerId != nil {
		return *m.ContainerId
	}
	return ""
}

func (m *BnsInstance) GetDeployPath() string {
	if m != nil && m.DeployPath != nil {
		return *m.DeployPath
	}
	return ""
}

func init() {
}
