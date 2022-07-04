package eventhandler

import (
	"main-service/pkg/logger"
	"main-service/pkg/structures/configuration"
)

var log = logger.GetLogger()

// TODO: Change response structure to instead be protobuf, it is currently normal structure
func HandleEvent(event *configuration.ConfigRequest) (*configuration.Response, error) {
	// Store requests in storage and log the events
	requestIds, err := storeRequestsInStore(event.Requests)
	if err != nil {
		log.Errorf("Failed storing and logging events: %v", err)
		return nil, err
	}

	log.Info("Configuration requests stored successfully!")

	// Notify TSN service that it should calculate a new configurations
	for _, req := range requestIds {
		// TODO: Make the code in this loop into a function and run in parallel?
		configId, err := notifyTsnService(req)
		if err != nil {
			log.Errorf("Failed to notify TSN service: %v", err)
			return nil, err
		}

		log.Infof("Configuration calculated with ID: %s", configId.GetValue())
	}

	// TODO: Finalize configuration

	// TODO: Send something to config-service to use new configuration

	return nil, nil
}
