package storewrapper

import (
	"context"
	"errors"
	"main-service/pkg/structures"

	"github.com/atomix/atomix-go-client/pkg/atomix"
	"github.com/gogo/protobuf/proto"
)

// Get the resource specified from the given store.
func getResource(storeName string, resourceName string) (interface{}, error) {
	ctx := context.Background()

	store, err := atomix.GetMap(ctx, storeName)
	if err != nil {
		log.Errorf("Error getting store: %v", err)
		return nil, err
	}

	resourceEntry, err := store.Get(ctx, resourceName)
	// _, err = store.Get(ctx, resourceName)
	if err != nil {
		log.Errorf("Error getting resource: %v", err)
		return nil, err
	}

	// resourceStruct, err := proto.Unmarshal(resourceEntry.Value, myProtoMessage)
	// if err != nil {
	// 	return nil, err
	// }

	resourceStruct, err := unmarshalMessage(resourceEntry.Value, resourceName)
	if err != nil {
		log.Errorf("Error unmarshaling resource: %v", err)
		return nil, err
	}

	return resourceStruct, nil
}

func unmarshalMessage(val []byte, resource string) (interface{}, error) {
	var err error

	switch resource {
	case "testing":
		obj := structures.ConfigRequest{}
		if err = proto.Unmarshal(val, &obj); err != nil {
			log.Errorf("Failed unmarshaling %s: %v", resource, err)
			return nil, err
		}

		return obj, nil
	// case "resources":
	// 	obj, err = getResource(resourceName)
	// case "streams":
	// 	obj, err = getStream(resourceName)
	// case "topology":
	// 	obj, err = getTopology(resourceName)
	// case "metrics":
	// 	obj, err = getMetric(resourceName)
	// case "events":
	// 	obj, err = getEvent(resourceName)
	default:
		return nil, errors.New("Store not found!")
	}

	if err != nil {
		log.Errorf("Could not find resource \"%s\": %v", resource, err)
		return nil, err
	}

	return "", nil
}

// var obj interface{}
// var err error

// switch storeName {
// case "configurations":
// 	obj, err = getConfiguration(resourceName)
// case "resources":
// 	obj, err = getResource(resourceName)
// case "streams":
// 	obj, err = getStream(resourceName)
// case "topology":
// 	obj, err = getTopology(resourceName)
// case "metrics":
// 	obj, err = getMetric(resourceName)
// case "events":
// 	obj, err = getEvent(resourceName)
// default:
// 	return nil, errors.New("Store not found!")
// }
