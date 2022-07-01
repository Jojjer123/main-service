package storewrapper

import (
	"context"
	"main-service/pkg/logger"
	"main-service/pkg/structures"
	"strings"

	"github.com/atomix/atomix-go-client/pkg/atomix"
	"github.com/golang/protobuf/proto"
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

func LogEvent(event *structures.Event) error {

	return nil
}

// Take in a config request from UNI and store it in the k/v
// store "streams" with a specific path for each request.
func StoreUniConfRequest(req *structures.Request) error {
	// Serialize request
	obj, err := proto.Marshal(req)
	if err != nil {
		log.Errorf("Failed to marshal request: %v", err)
		return err
	}

	// Create a URN where the serialized request will be stored
	urn := "streams.requests."

	// TODO: Generate or use some ID to keep track of the specific stream request
	urn += "someId"

	// Send serialized request to it's specific path in a store
	err = sendToStore(obj, urn)
	if err != nil {
		return err
	}

	return nil
}

// Takes in an object as a byte slice, a URN in the
// format of "storeName.Resource", and stores the
// structure at the URN.
func sendToStore(obj []byte, urn string) error {
	ctx := context.Background()

	// Create a slice of URN elements
	urnElems := strings.Split(urn, ".")

	// Get the store
	store, err := atomix.GetMap(ctx, urnElems[0])
	if err != nil {
		log.Errorf("Failed getting store \"%s\": %v", urnElems[0], err)
		return err
	}

	// TODO: Check if the URN contains more complex path and do something special then

	// Store the object
	_, err = store.Put(ctx, urnElems[1], obj)
	if err != nil {
		log.Errorf("Failed storing resource \"%s\": %v", urnElems[1], err)
		return err
	}

	return nil
}

// Takes in a URN in the format "storeName.Resource.Resource" and
// returns the structure for the requested resource.
func Get(urn string) (interface{}, error) {
	// Create a slice of urn elements
	urnElems := strings.Split(urn, ".")

	// Request object from store
	obj, err := getObjectFromStore(urnElems[0], urnElems[1])
	if err != nil {
		log.Errorf("Failed getting object from store: %v", err)
		return nil, err
	}

	// Get requested resource from the object

	return obj, nil
}

// Takes in name of store and resource and returns structure
// containing the resource.
func getObjectFromStore(storeName string, resourceName string) (interface{}, error) {
	obj, err := getResource(storeName, resourceName)

	return obj, err
}
