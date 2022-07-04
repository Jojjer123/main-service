package southboundinterface

import (
	// "context"

	"main-service/pkg/logger"
	// store "main-service/pkg/store-wrapper"
	"main-service/pkg/structures/configuration"
	// "github.com/gogo/protobuf/proto"
	// gclient "github.com/openconfig/gnmi/client/gnmi"
	// "github.com/openconfig/gnmi/proto/gnmi"
)

var log = logger.GetLogger()

func StoreRequestInStorage(configRequest *configuration.ConfigRequest) {
	// data, err := proto.Marshal(configRequest)
	// if err != nil {
	// 	log.Errorf("Failed marshaling: %v", err)
	// 	return
	// }

	// err = store.Set("configurations.testing", data)
	// if err != nil {
	// 	log.Errorf("Error storing configuration request: %v", err)
	// 	return
	// }

	// log.Info("Stored configuration request!")

	// ctx := context.Background()

	// c, err := createGnmiClient("storage-service", ctx)
	// if err != nil {
	// 	return
	// }

	// defer c.Close()

	// request := createSetRequest(data)

	// var response *gnmi.SetResponse
	// response, err = c.(*gclient.Client).Set(ctx, request)
	// if err != nil {
	// 	log.Errorf("Set request failed: %v", err)
	// 	return
	// }

	// if len(response.Response) > 1 {
	// 	log.Error("More than one result from storage-service")
	// }

	// for _, result := range response.Response {
	// 	if result.Path.Elem[0].Name != "ActionResult" {
	// 		log.Error("Missing action result from storage-service")
	// 	} else {
	// 		if result.Path.Elem[0].Key["ActionResult"] == "Failed" {
	// 			log.Error("Storing request failed!")
	// 		} else {
	// 			log.Infof("Stored request successfully!")
	// 		}
	// 	}
	// }
}

func GetConfigFromStorage() []byte {
	// resource, err := store.Get("configurations.testing")
	// if err != nil {
	// 	log.Errorf("Error storing configuration request: %v", err)
	// 	return nil
	// }

	// test, ok := resource.(structures.ConfigRequest)
	// if !ok {
	// 	log.Errorf("Failed type assertion on: %v", test)
	// 	return nil
	// }

	// config, err := proto.Marshal(&test)
	// if err != nil {
	// 	log.Errorf("Failed marshaling config: %v", err)
	// 	return nil
	// }

	// return config
	return []byte{}
}
