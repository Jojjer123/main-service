package eventhandler

import (
	"context"

	"main-service/pkg/structures/notification"

	"google.golang.org/grpc"
)

// Notifies the TSN service through gRPC that it should start calculating
// a new configuration.
func notifyTsnService(reqId *notification.UUID) (*notification.UUID, error) {
	// Create gRPC client and connect to TSN service
	// (consider having a constant connection to TSN service)
	conn, err := grpc.Dial("tsn-service:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed dialing tsn-service: %v", err)
		return nil, err
	}

	defer conn.Close()

	client := notification.NewNotificationClient(conn)

	confId, err := client.CalcConfig(context.Background(), reqId)
	if err != nil {
		log.Fatalf("Calculating configuration failed: %v", err)
		return nil, err
	}

	return confId, nil
}
