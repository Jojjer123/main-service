package northboundinterface

import (
	"main-service/pkg/structures/configuration"
	"main-service/pkg/structures/notification"

	"google.golang.org/protobuf/encoding/protojson"
)

func createResponse(confId *notification.UUID) ([]byte, error) {
	var baseResp = &configuration.ConfigResponse{
		Version: 123, // random value for now
		Responses: []*configuration.Response{
			{
				StatusGroup: &configuration.StatusGroup{
					StrId: &configuration.StreamId{
						MacAddress: "11:22:33:44:55", // random value for now
						UniqueId:   "123",            // random value for now
					},
					StatusInfo: &configuration.StatusInfo{
						TalkerStatus:   1,   // random value for now
						ListenerStatus: 1,   // random value for now
						FailureCode:    123, // random value for now
					},
					FailedInterfaces: []*configuration.InterfaceId{
						{
							MacAddress:    "11:22:33:44:55", // random value for now
							InterfaceName: "I-do-not-know",  // random value for now
						},
					},
					StatusTalkerListener: []*configuration.TalkerListenerStatus{
						{
							AccumulatedLatency: &configuration.AccumulatedLatency{
								AccumulatedLatency: 123, // random value for now
							},
							InterfaceConfiguration: []*configuration.InterfaceConfiguration{
								{
									InterfaceId: &configuration.InterfaceId{
										MacAddress:    "11:22:33:44:55", // random value for now
										InterfaceName: "test",           // random value for now
									},
									Type: 123, // random value for now
									MacAddr: &configuration.IeeeMacAddress{
										DestinationMac: "destMac", // random value for now
										SourceMac:      "srcMac",  // random value for now
									},
									VlanTag: &configuration.IeeeVlanTag{
										PriorityCodePoint: 1, // random value for now
										VlanId:            1, // random value for now
									},
									Ipv4Tup: &configuration.Ipv4Tuple{},
									Ipv6Tup: &configuration.Ipv6Tuple{},
									TimeAwareOffset: &configuration.TimeAwareOffset{
										Offset: 123, // random value for now
									},
								},
							},
						},
					},
					EndStationInterfaces: []*configuration.Interface{
						{
							Index: 0, // random value for now
							InterfaceId: &configuration.InterfaceId{
								MacAddress:    "11:22:33:44:55", // random value for now
								InterfaceName: "testInterface",  // random value for now
							},
						},
					},
				},
			},
		},
	}

	rawData, err := protojson.Marshal(baseResp)
	if err != nil {
		log.Errorf("Failed to marshal UNI response: %v", err)
		return nil, err
	}

	return rawData, nil
}
