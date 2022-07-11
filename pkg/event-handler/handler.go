package eventhandler

import (
	"main-service/pkg/logger"
	// store "main-service/pkg/store-wrapper"
	"main-service/pkg/structures/configuration"
	"main-service/pkg/structures/notification"

	"time"
)

var log = logger.GetLogger()

// Take in a configuratin request, process it and once a configuration
// has been calculated, return ID of the new configuration.
func HandleAddStreamEvent(event *configuration.ConfigRequest, timeOfReq time.Time) (*notification.UUID, error) {
	// start := time.Now().UnixMilli()

	// Store requests in storage and log the events
	requestIds, err := storeRequestsInStore(event.Requests, timeOfReq)
	if err != nil {
		log.Errorf("Failed storing and logging events: %v", err)
		return nil, err
	}

	log.Info("Configuration requests stored successfully!")

	// Notify TSN service that it should calculate a new configuration
	configId, err := notifyTsnService(requestIds)
	if err != nil {
		log.Errorf("Failed to notify TSN service: %v", err)
		return nil, err
	}

	log.Infof("Configuration calculated with ID: %s", configId.GetValue())

	// end := time.Now().UnixMilli()
	// log.Infof("Time to complete: %v ms", end-start)

	// TODO: Finalize configuration

	// TODO: Send something to config-service to use new configuration
	if err = applyConfiguration(configId); err != nil {
		log.Errorf("Failed notifying onos-config of new configuration: %v", err)
		return nil, err
	}

	return configId, nil
}
