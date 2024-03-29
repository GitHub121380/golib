// Code generated by protoc-gen-go.
// source: service.proto
// DO NOT EDIT!

package bns

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type InstanceInfo struct {
	HostName         *string         `protobuf:"bytes,1,opt,name=host_name" json:"host_name,omitempty"`
	ServiceName      *string         `protobuf:"bytes,2,opt,name=service_name" json:"service_name,omitempty"`
	InstanceStatus   *InstanceStatus `protobuf:"bytes,3,opt,name=instance_status" json:"instance_status,omitempty"`
	InstanceLoad     *InstanceLoad   `protobuf:"bytes,4,opt,name=instance_load" json:"instance_load,omitempty"`
	HostIp           *uint32         `protobuf:"varint,5,opt,name=host_ip" json:"host_ip,omitempty"`
	SuspectStatus    *bool           `protobuf:"varint,6,opt,name=suspect_status,def=0" json:"suspect_status,omitempty"`
	Offset           *int32          `protobuf:"varint,7,opt,name=offset,def=0" json:"offset,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *InstanceInfo) Reset()         { *m = InstanceInfo{} }
func (m *InstanceInfo) String() string { return proto.CompactTextString(m) }
func (*InstanceInfo) ProtoMessage()    {}

const Default_InstanceInfo_SuspectStatus bool = false
const Default_InstanceInfo_Offset int32 = 0

func (m *InstanceInfo) GetHostName() string {
	if m != nil && m.HostName != nil {
		return *m.HostName
	}
	return ""
}

func (m *InstanceInfo) GetServiceName() string {
	if m != nil && m.ServiceName != nil {
		return *m.ServiceName
	}
	return ""
}

func (m *InstanceInfo) GetInstanceStatus() *InstanceStatus {
	if m != nil {
		return m.InstanceStatus
	}
	return nil
}

func (m *InstanceInfo) GetInstanceLoad() *InstanceLoad {
	if m != nil {
		return m.InstanceLoad
	}
	return nil
}

func (m *InstanceInfo) GetHostIp() uint32 {
	if m != nil && m.HostIp != nil {
		return *m.HostIp
	}
	return 0
}

func (m *InstanceInfo) GetSuspectStatus() bool {
	if m != nil && m.SuspectStatus != nil {
		return *m.SuspectStatus
	}
	return Default_InstanceInfo_SuspectStatus
}

func (m *InstanceInfo) GetOffset() int32 {
	if m != nil && m.Offset != nil {
		return *m.Offset
	}
	return Default_InstanceInfo_Offset
}

type InstanceLoad struct {
	Load             *int32 `protobuf:"varint,1,opt,name=load" json:"load,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InstanceLoad) Reset()         { *m = InstanceLoad{} }
func (m *InstanceLoad) String() string { return proto.CompactTextString(m) }
func (*InstanceLoad) ProtoMessage()    {}

func (m *InstanceLoad) GetLoad() int32 {
	if m != nil && m.Load != nil {
		return *m.Load
	}
	return 0
}

type InstanceStatus struct {
	Port                 *int32  `protobuf:"varint,1,opt,name=port" json:"port,omitempty"`
	Status               *int32  `protobuf:"varint,2,opt,name=status" json:"status,omitempty"`
	Tags                 *string `protobuf:"bytes,3,opt,name=tags" json:"tags,omitempty"`
	Extra                *string `protobuf:"bytes,4,opt,name=extra" json:"extra,omitempty"`
	InterventionalStatus *int32  `protobuf:"varint,5,opt,name=interventional_status,def=0" json:"interventional_status,omitempty"`
	MultiPort            *string `protobuf:"bytes,6,opt,name=multi_port" json:"multi_port,omitempty"`
	ContainerId          *string `protobuf:"bytes,7,opt,name=container_id" json:"container_id,omitempty"`
	DeployPath           *string `protobuf:"bytes,8,opt,name=deploy_path" json:"deploy_path,omitempty"`
	XXX_unrecognized     []byte  `json:"-"`
}

func (m *InstanceStatus) Reset()         { *m = InstanceStatus{} }
func (m *InstanceStatus) String() string { return proto.CompactTextString(m) }
func (*InstanceStatus) ProtoMessage()    {}

const Default_InstanceStatus_InterventionalStatus int32 = 0

func (m *InstanceStatus) GetPort() int32 {
	if m != nil && m.Port != nil {
		return *m.Port
	}
	return 0
}

func (m *InstanceStatus) GetStatus() int32 {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return 0
}

func (m *InstanceStatus) GetTags() string {
	if m != nil && m.Tags != nil {
		return *m.Tags
	}
	return ""
}

func (m *InstanceStatus) GetExtra() string {
	if m != nil && m.Extra != nil {
		return *m.Extra
	}
	return ""
}

func (m *InstanceStatus) GetInterventionalStatus() int32 {
	if m != nil && m.InterventionalStatus != nil {
		return *m.InterventionalStatus
	}
	return Default_InstanceStatus_InterventionalStatus
}

func (m *InstanceStatus) GetMultiPort() string {
	if m != nil && m.MultiPort != nil {
		return *m.MultiPort
	}
	return ""
}

func (m *InstanceStatus) GetContainerId() string {
	if m != nil && m.ContainerId != nil {
		return *m.ContainerId
	}
	return ""
}

func (m *InstanceStatus) GetDeployPath() string {
	if m != nil && m.DeployPath != nil {
		return *m.DeployPath
	}
	return ""
}

type ServiceHostList struct {
	ServiceName      *string       `protobuf:"bytes,1,opt,name=service_name" json:"service_name,omitempty"`
	HostName         []string      `protobuf:"bytes,2,rep,name=host_name" json:"host_name,omitempty"`
	HostIpPair       []*HostIpPair `protobuf:"bytes,3,rep,name=host_ip_pair" json:"host_ip_pair,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *ServiceHostList) Reset()         { *m = ServiceHostList{} }
func (m *ServiceHostList) String() string { return proto.CompactTextString(m) }
func (*ServiceHostList) ProtoMessage()    {}

func (m *ServiceHostList) GetServiceName() string {
	if m != nil && m.ServiceName != nil {
		return *m.ServiceName
	}
	return ""
}

func (m *ServiceHostList) GetHostName() []string {
	if m != nil {
		return m.HostName
	}
	return nil
}

func (m *ServiceHostList) GetHostIpPair() []*HostIpPair {
	if m != nil {
		return m.HostIpPair
	}
	return nil
}

type ServiceAuthList struct {
	ServiceName      *string  `protobuf:"bytes,1,opt,name=service_name" json:"service_name,omitempty"`
	AuthServiceName  []string `protobuf:"bytes,2,rep,name=auth_service_name" json:"auth_service_name,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *ServiceAuthList) Reset()         { *m = ServiceAuthList{} }
func (m *ServiceAuthList) String() string { return proto.CompactTextString(m) }
func (*ServiceAuthList) ProtoMessage()    {}

func (m *ServiceAuthList) GetServiceName() string {
	if m != nil && m.ServiceName != nil {
		return *m.ServiceName
	}
	return ""
}

func (m *ServiceAuthList) GetAuthServiceName() []string {
	if m != nil {
		return m.AuthServiceName
	}
	return nil
}

type HostIpPair struct {
	HostName         *string `protobuf:"bytes,1,opt,name=host_name" json:"host_name,omitempty"`
	HostIp           *uint32 `protobuf:"varint,2,opt,name=host_ip" json:"host_ip,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *HostIpPair) Reset()         { *m = HostIpPair{} }
func (m *HostIpPair) String() string { return proto.CompactTextString(m) }
func (*HostIpPair) ProtoMessage()    {}

func (m *HostIpPair) GetHostName() string {
	if m != nil && m.HostName != nil {
		return *m.HostName
	}
	return ""
}

func (m *HostIpPair) GetHostIp() uint32 {
	if m != nil && m.HostIp != nil {
		return *m.HostIp
	}
	return 0
}

type ServiceInfo struct {
	ServiceName       *string `protobuf:"bytes,1,opt,name=service_name" json:"service_name,omitempty"`
	ServiceConf       *string `protobuf:"bytes,2,opt,name=service_conf" json:"service_conf,omitempty"`
	Threshold         *int32  `protobuf:"varint,3,opt,name=threshold" json:"threshold,omitempty"`
	CustomDefine      *string `protobuf:"bytes,4,opt,name=custom_define" json:"custom_define,omitempty"`
	OpenDeadhostCheck *bool   `protobuf:"varint,5,opt,name=open_deadhost_check,def=0" json:"open_deadhost_check,omitempty"`
	ThresholdPercent  *int32  `protobuf:"varint,6,opt,name=threshold_percent" json:"threshold_percent,omitempty"`
	OpenSmartBns      *bool   `protobuf:"varint,7,opt,name=open_smart_bns,def=0" json:"open_smart_bns,omitempty"`
	GroupNames        *string `protobuf:"bytes,8,opt,name=group_names" json:"group_names,omitempty"`
	QosInfo           *string `protobuf:"bytes,9,opt,name=qos_info" json:"qos_info,omitempty"`
	QosOpen           *bool   `protobuf:"varint,10,opt,name=qos_open,def=0" json:"qos_open,omitempty"`
	GianoInfo         *string `protobuf:"bytes,11,opt,name=giano_info" json:"giano_info,omitempty"`
	XXX_unrecognized  []byte  `json:"-"`
}

func (m *ServiceInfo) Reset()         { *m = ServiceInfo{} }
func (m *ServiceInfo) String() string { return proto.CompactTextString(m) }
func (*ServiceInfo) ProtoMessage()    {}

const Default_ServiceInfo_OpenDeadhostCheck bool = false
const Default_ServiceInfo_OpenSmartBns bool = false
const Default_ServiceInfo_QosOpen bool = false

func (m *ServiceInfo) GetServiceName() string {
	if m != nil && m.ServiceName != nil {
		return *m.ServiceName
	}
	return ""
}

func (m *ServiceInfo) GetServiceConf() string {
	if m != nil && m.ServiceConf != nil {
		return *m.ServiceConf
	}
	return ""
}

func (m *ServiceInfo) GetThreshold() int32 {
	if m != nil && m.Threshold != nil {
		return *m.Threshold
	}
	return 0
}

func (m *ServiceInfo) GetCustomDefine() string {
	if m != nil && m.CustomDefine != nil {
		return *m.CustomDefine
	}
	return ""
}

func (m *ServiceInfo) GetOpenDeadhostCheck() bool {
	if m != nil && m.OpenDeadhostCheck != nil {
		return *m.OpenDeadhostCheck
	}
	return Default_ServiceInfo_OpenDeadhostCheck
}

func (m *ServiceInfo) GetThresholdPercent() int32 {
	if m != nil && m.ThresholdPercent != nil {
		return *m.ThresholdPercent
	}
	return 0
}

func (m *ServiceInfo) GetOpenSmartBns() bool {
	if m != nil && m.OpenSmartBns != nil {
		return *m.OpenSmartBns
	}
	return Default_ServiceInfo_OpenSmartBns
}

func (m *ServiceInfo) GetGroupNames() string {
	if m != nil && m.GroupNames != nil {
		return *m.GroupNames
	}
	return ""
}

func (m *ServiceInfo) GetQosInfo() string {
	if m != nil && m.QosInfo != nil {
		return *m.QosInfo
	}
	return ""
}

func (m *ServiceInfo) GetQosOpen() bool {
	if m != nil && m.QosOpen != nil {
		return *m.QosOpen
	}
	return Default_ServiceInfo_QosOpen
}

func (m *ServiceInfo) GetGianoInfo() string {
	if m != nil && m.GianoInfo != nil {
		return *m.GianoInfo
	}
	return ""
}

func init() {
}
