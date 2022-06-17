package southboundinterface

import (
	"context"

	"main-service/pkg/logger"

	gclient "github.com/openconfig/gnmi/client/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
)

var log = logger.GetLogger()

func SendRequestToStorage(data []byte) {
	ctx := context.Background()

	c, err := createGnmiClient("storage-service", ctx)
	if err != nil {
		return
	}

	defer c.Close()

	request := createSetRequest(data)

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

func GetConfigFromStorage() []byte {
	ctx := context.Background()

	c, err := createGnmiClient("storage-service", ctx)
	if err != nil {
		return nil
	}

	defer c.Close()

	request := createGetRequest()

	var response *gnmi.GetResponse
	response, err = c.(*gclient.Client).Get(ctx, request)
	if err != nil {
		log.Errorf("Set request failed: %v", err)
		return nil
	}

	if len(response.Notification) > 1 {
		log.Error("More than one update from storage-service")
	}

	var data []byte

	for _, notification := range response.Notification {
		log.Info("Received main conf from storage!")
		data = notification.Update[0].Val.GetJsonVal()
	}

	return data
}
