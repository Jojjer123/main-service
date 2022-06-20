package structures

import "net"

//////////////////////////////////
//      Request structure       //
//////////////////////////////////

type UserToNetworkRequirements struct {
	NumSeamlessTrees uint8  `json:"num-seamless-trees"`
	MaxLatency       uint32 `json:"max-latency"`
}

type InterfaceCapabilities struct {
	VlanTagCapable bool `json:"vlan-tag-capable"`
}

type StreamId struct {
	// MacAddress net.HardwareAddr
	// UniqueId   uint16
	MacAddress string `json:"mac-address"`
	UniqueId   string `json:"unique-id"`
}

type InterfaceId struct {
	// MacAddress    net.HardwareAddr
	// InterfaceName string
	Index         int    `json:"index"`
	MacAddress    string `json:"mac-address"`
	InterfaceName string `json:"interface-name"`
}

type ListenerGroup struct {
	// StrId streamId
	Index                uint16                    `json:"index"`
	StrId                StreamId                  `json:"stream-id"`
	EndStationInterfaces []InterfaceId             `json:"end-station-interfaces"`
	UserToNetReq         UserToNetworkRequirements `json:"user-to-network-requirements"`
	InterfCap            InterfaceCapabilities     `json:"interface-capabilities"`
}

////////////////////////////

type StreamRank struct {
	Rank uint8 `json:"Rank"`
}

type IeeeMacAddress struct {
	// DestinationMac net.HardwareAddr
	// SourceMac      net.HardwareAddr
	DestinationMac string `json:"destination-mac-address"`
	SourceMac      string `json:"source-mac"`
}

type IeeeVlanTag struct {
	PriorityCodePoint uint8  `json:"priority-code-point"`
	VlanId            uint16 `json:"vlan-id"`
}

type Ipv4Tuple struct {
	SrcIpAddr  net.IPAddr
	DestIpAddr net.IPAddr
	Dscp       uint8
	Protocol   uint16
	SrcPort    uint16
	DestPort   uint16
}

type Ipv6Tuple struct {
	SrcIpAddr  net.IPAddr
	DestIpAddr net.IPAddr
	Dscp       uint8
	Protocol   uint16
	SrcPort    uint16
	DestPort   uint16
}

// type DataFrameSpecification struct {
// 	Typ     int // 1 = macaddress/vlan, 2 = ipv4, 3 = ipv6
// 	MacAddr *IeeeMacAddress
// 	VlanTag *IeeeVlanTag
// 	Ipv4Tup *Ipv4Tuple
// 	Ipv6Tup *Ipv6Tuple
// }

type DataFrameSpecification struct {
	Index   uint16
	MacAddr *IeeeMacAddress `json:"ieee802-mac-addresses"`
	VlanTag *IeeeVlanTag    `json:"ieee802-vlan-tag"`
	Ipv4Tup *Ipv4Tuple      //`json:""`
	Ipv6Tup *Ipv6Tuple      //`json:""`
}

type Interval struct {
	Numerator   uint32
	Denominator uint32
}

type TimeAware struct {
	EarliestTransmitOffset uint32 `json:"earliest-transmit-offset"`
	LatestTransmitOffset   uint32 `json:"latest-transmit-offset"`
	Jitter                 uint32 `json:"jitter"`
}

type TrafficSpecification struct {
	Interval              Interval   `json:"Interval"`
	MaxFramesPerInterval  uint16     `json:"max-frames-per-interval"`
	MaxFrameSize          uint16     `json:"max-frame-size"`
	TransmissionSelection uint8      `json:"transmission-selection"`
	TimeAware             *TimeAware `json:"time-aware"`
}

type TalkerGroup struct {
	StrId                  StreamId                  `json:"stream-id"`
	StrRank                StreamRank                `json:"stream-rank"`
	EndStationInterfaces   []InterfaceId             `json:"end-station-interfaces"`
	DataFrameSpecification []DataFrameSpecification  `json:"data-frame-specification"`
	TrafficSpecification   TrafficSpecification      `json:"traffic-specification"`
	UserToNetReq           UserToNetworkRequirements `json:"user-to-network-requirements"`
	InterfCap              InterfaceCapabilities     `json:"interface-capabilities"`
}

type Request struct {
	Talker       TalkerGroup     `json:"talker"`
	ListenerList []ListenerGroup `json:"listener-list"`
}

////////////////////////////

// type ConfigRequest struct {
// 	Talker       TalkerGroup
// 	ListenerList []ListenerGroup
// }
type ConfigRequest struct {
	Version  float32   `json:"version"`
	Requests []Request `json:"requests"`
}

//////////////////////////////////
//      Response structure      //
//////////////////////////////////

type StatusInfo struct {
	TalkerStatus   int //todo: update to enum: 0=None, 1=ready, 2=failed
	ListenerStatus int //todo: update to enum: 0=None, 1=ready, 2=partial-failed, 3=failed
	FailureCode    uint8
}

type AccumulatedLatency struct {
	AccumulatedLatency uint32
}

type TimeAwareOffset struct {
	Offset uint32 //ns, EarliestTransmitOffset < offset < LatestTransmitOffset
}

type InterfaceConfiguration struct {
	InterfaceId InterfaceId
	Typ         int // 1 = macaddress/vlan, 2 = ipv4, 3 = ipv6
	MacAddr     *IeeeMacAddress
	VlanTag     *IeeeVlanTag
	Ipv4Tup     *Ipv4Tuple
	Ipv6Tup     *Ipv6Tuple
	// provided iff TimeAware is present in talker group
	TimeAwareOffset *TimeAwareOffset
}

type TalkerListenerStatus struct {
	AccumulatedLatency     AccumulatedLatency
	InterfaceConfiguration []InterfaceConfiguration
}

type StatusGroup struct {
	//StrId StreamId
	StatusInfo           StatusInfo
	FailedInterfaces     *[]InterfaceId
	StatusTalkerListener []TalkerListenerStatus
	EndStationInterfaces *[]InterfaceId
}

type Response struct {
	StatusGroup StatusGroup
}
