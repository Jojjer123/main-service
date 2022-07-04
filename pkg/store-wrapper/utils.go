package storewrapper

import (
	"context"
	"strings"

	"github.com/atomix/atomix-go-client/pkg/atomix"
)

// Takes in an object as a byte slice, a URN in the
// format of "storeName.Resource", and stores the
// structure at the URN.
func sendToStore(obj []byte, urn string) error {
	ctx := context.Background()

	// Create a slice of URN elements
	urnElems := strings.SplitN(urn, ".", 2)

	// log.Infof("Getting store \"%s\"...", urnElems[0])

	// Get the store
	store, err := atomix.GetMap(ctx, urnElems[0])
	if err != nil {
		log.Errorf("Failed getting store \"%s\": %v", urnElems[0], err)
		return err
	}

	// TODO: Check if the URN contains more complex path and do something special then

	// log.Infof("Storing object at \"%s\"...", urnElems[1])

	// Store the object
	_, err = store.Put(ctx, urnElems[1], obj)
	if err != nil {
		log.Errorf("Failed storing resource \"%s\": %v", urnElems[1], err)
		return err
	}

	// log.Infof("Stored object at \"%s\"", urn)

	return nil
}
