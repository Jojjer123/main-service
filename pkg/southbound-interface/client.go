package southboundinterface

import (
	"context"
	"time"

	"github.com/main-service/pkg/logger"
	"github.com/openconfig/gnmi/client"
	"github.com/openconfig/gnmi/proto/gnmi"

	gclient "github.com/openconfig/gnmi/client/gnmi"
)

var log = logger.GetLogger()

func SendRequestToStorage(data []byte) {
	ctx := context.Background()

	c, err := createGnmiClient("storage-service", ctx)
	if err != nil {
		return
	}

	request := createRequest(data)

	var response *gnmi.SetResponse
	response, err = c.(*gclient.Client).Set(ctx, request)
	if err != nil {
		log.Errorf("Set request failed: %v", err)
		return
	}

	if len(response.Response) > 1 {
		log.Error("More than one result from storage-service")
	}

	for _, result := range response.Response {
		if result.Path.Elem[0].Name != "ActionResult" {
			log.Error("Missing action result from storage-service")
		} else {
			if result.Path.Elem[0].Key["ActionResult"] == "Failed" {
				log.Error("Storing request failed!")
			} else {
				log.Infof("Stored request successfully!")
			}
		}
	}
}

func createRequest(data []byte) *gnmi.SetRequest {
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

func createGnmiClient(addr string, ctx context.Context) (client.Impl, error) {
	// Use secure communication (port 10161)
	c, err := gclient.New(ctx, client.Destination{
		Addrs:       []string{addr + ":11161"},
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
