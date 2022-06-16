package southboundinterface

import (
	"context"
	"fmt"
	"time"

	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
)

func createSetRequest(data []byte) *gnmi.SetRequest {
	request := &gnmi.SetRequest{
		Update: []*gnmi.Update{
			{
				Path: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{
							Name: "Action",
							Key: map[string]string{
								"Action": "StoreGetReq",
							},
						},
					},
				},
				Val: &gnmi.TypedValue{
					Value: &gnmi.TypedValue_JsonVal{
						JsonVal: data,
					},
				},
			},
		},
	}

	return request
}

func createGetRequest() *gnmi.GetRequest {
	request := &gnmi.GetRequest{
		Path: []*gnmi.Path{
			{
				Elem: []*gnmi.PathElem{
					{
						Name: "Action",
						Key: map[string]string{
							"Action": "GetMainConf",
						},
					},
				},
			},
		},
	}

	return request
}

func createGnmiClient(addr string, ctx context.Context) (client.Impl, error) {
	// Use secure communication (port 10161)
	c, err := gclient.New(ctx, client.Destination{
		Addrs:       []string{fmt.Sprintf("%s:11161", addr)},
		Target:      addr,
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		log.Errorf("Could not create a gNMI client: %+v", err)

		return nil, err
	}

	return c, nil
}
