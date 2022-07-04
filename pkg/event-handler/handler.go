package eventhandler

import (
	"main-service/pkg/logger"
	"main-service/pkg/structures/configuration"
	"time"
)

var log = logger.GetLogger()

// Take in a configuratin request, process it and once a configuration
// has been calculated, return a configuration response.
func HandleEvent(event *configuration.ConfigRequest) (*configuration.Response, error) {

	start := time.Now().UnixMilli()

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

	end := time.Now().UnixMilli()

	log.Infof("Time to complete: %v ms", end-start)

	// TODO: Finalize configuration

	// TODO: Send something to config-service to use new configuration

	return nil, nil
}
