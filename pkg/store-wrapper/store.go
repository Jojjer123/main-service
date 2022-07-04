package storewrapper

import (
	"context"
	"fmt"
	"main-service/pkg/logger"
	"main-service/pkg/structures/configuration"
	"main-service/pkg/structures/event"
	"main-service/pkg/structures/notification"

	"github.com/atomix/atomix-go-client/pkg/atomix"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

var stores = []string{
	"configurations",
	"resources",
	"streams",
	"topology",
	"metrics",
	"events",
}

var log = logger.GetLogger()

// Generates all stores defined in the global variable "stores"
// (they are all generated as the primitive Map).
func CreateStores() {
	ctx := context.Background()

	var successful = true

	// For each store-name generate a new k/v store
	for _, name := range stores {
		_, err := atomix.GetMap(ctx, name)
		if err != nil {
			log.Errorf("Failed creating store \"%s\": %v", name, err)
			successful = false
		}
	}

	if successful {
		log.Info("All stores created!")
	}
}

// Log an event to k/v store
func LogEvent(event *event.Event) error {
	// Serialize event
	obj, err := proto.Marshal(event)
	if err != nil {
		log.Errorf("Failed to marshal request: %v", err)
		return err
	}

	// Create a URN where the serialized request will be stored
	urn := "events.addStream."

	// TODO: Generate or use some ID to keep track of the specific event
	// urn += fmt.Sprintf("%v", uuid.New())
	urn += fmt.Sprintf("%v", event.EventId)

	// Send serialized event to it's specific path in a store
	err = sendToStore(obj, urn)
	if err != nil {
		return err
	}

	return nil
}

// Take in a config request from UNI and store it in the k/v
// store "streams" with a specific path for each request.
func StoreUniConfRequest(req *configuration.Request) (*notification.UUID, error) {
	// Serialize request
	obj, err := proto.Marshal(req)
	if err != nil {
		log.Errorf("Failed to marshal request: %v", err)
		return nil, err
	}

	// Create a URN where the serialized request will be stored
	urn := "streams.requests."

	// TODO: Generate or use some ID to keep track of the specific stream request
	var requestId = notification.UUID{
		Value: fmt.Sprintf("%v", uuid.New()),
	}
	urn += fmt.Sprintf("%v", requestId.Value)

	// log.Infof("URN now looks like: %s", urn)

	// Send serialized request to it's specific path in a store
	err = sendToStore(obj, urn)
	if err != nil {
		return nil, err
	}

	return &requestId, nil
}

//////////////////////////////////////////////////
/*                   TEMPLATE                   */
//////////////////////////////////////////////////
/*

func PublicFunctionName(req structureType) error {
	// Serialize request
	obj, err := proto.Marshal(req)
	if err != nil {
		log.Errorf("Failed to marshal request: %v", err)
		return err
	}

	// Create a URN where the serialized request will be stored
	urn := "store.type."

	// TODO: Generate or use some ID to keep track of the specific stream request
	urn += fmt.Sprintf("%v", uuid.New())

	// Send serialized request to it's specific path in a store
	err = sendToStore(obj, urn)
	if err != nil {
		return err
	}

	return nil
}

*/
