package eventhandler

import (
	"context"
	"fmt"
	"main-service/pkg/logger"
	store "main-service/pkg/store-wrapper"
	"main-service/pkg/structures"
	"main-service/pkg/structures/grpc/notification"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var log = logger.GetLogger()

// TODO: Change response structure to instead be protobuf, it is currently normal structure
func HandleEvent(event *structures.ConfigRequest) (*structures.Response, error) {
	// Store requests in storage and log the events
	if err := storeRequestsInStore(event.Requests); err != nil {
		log.Errorf("Failed storing and logging events: %v", err)
		return nil, err
	}

	log.Info("Configuration requests stored successfully!")

	// Notify TSN service that it should calculate a new configuration
	if err := notifyTsnService(); err != nil {
		log.Errorf("Failed to notify TSN service: %v", err)
		return nil, err
	}

	log.Info("Configuration calculated!")

	// TODO: Finalize configuration

	// TODO: Send something to config-service to use new configuration

	return nil, nil
}

// Notifies the TSN service through gRPC that it should start calculating
// a new configuration.
func notifyTsnService() error {
	// TODO: Create gRPC client and connect to TSN service
	// (consider having a constant connection to TSN service)

	conn, err := grpc.Dial("tsn-service:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed dialing tsn-service: %v", err)
		return err
	}

	defer conn.Close()

	client := notification.NewNotificationClient(conn)

	// Will have to be changed, CalcConfig should return something so that
	// we know what configuration it just calculated.
	_, err = client.CalcConfig(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Calculating configuration failed: %v", err)
		return err
	}

	return nil
}

// Takes in requests, stores them, and logs the events
func storeRequestsInStore(requestList []*structures.Request) error {
	var storingOk = true
	var err error

	// Store all requests in a k/v store
	for _, request := range requestList {
		err = store.StoreUniConfRequest(request)
		if err != nil {
			storingOk = false
		} else {
			if err = storeEvent(request); err != nil {
				return err
			}
		}
	}

	// Stop handling event if storing of configurations failed
	if !storingOk {
		log.Errorf("Storing configuration requests failed: %v", err)
		return err
	}

	return nil
}

// Create and store an event
func storeEvent(req *structures.Request) error {
	// Create an event from the request
	ev, err := createEvent(req)
	if err != nil {
		log.Errorf("Failed creating event from request: %v", err)
		return err
	}

	// Log the event
	if err = store.LogEvent(ev); err != nil {
		log.Errorf("Failed to log event: %v", err)
		return err
	}

	return nil
}

// Create an event from the request
func createEvent(req *structures.Request) (*structures.Event, error) {
	// TODO: Add correct data to event, still don't know:
	// 		* Event types that should exist
	// 		* What is a Handler
	// 		* Should EventGroupId just be a uuid?
	// 		* OccuranceTime comes from where?
	// 		* Duration is measured where?
	// 		* What is LogInfo?

	var event = structures.Event{
		EventId:       fmt.Sprintf("%v", uuid.New()),
		EventType:     structures.EventType_ADD_STREAM,
		Status:        structures.EventStatus_PASSED,
		Handlers:      []*structures.EventHandler{},
		EventGroupId:  fmt.Sprintf("%v", uuid.New()),
		OccuranceTime: 123,
		Duration:      123,
		LogInfo:       &structures.LogInfo{},
	}

	return &event, nil
}
