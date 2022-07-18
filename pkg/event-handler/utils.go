package eventhandler

import (
	"context"
	"crypto/tls"
	"strings"

	// "crypto/tls"
	"main-service/pkg/structures/notification"
	"time"

	// "github.com/onosproject/onos-lib-go/pkg/certs"
	// "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi_ext"
	"google.golang.org/grpc"
)

// Notifies the TSN service through gRPC that it should start calculating
// a new configuration.
func notifyTsnService(reqIds *notification.IdList) (*notification.UUID, error) {
	// Create gRPC client and connect to TSN service
	// (consider having a constant connection to TSN service)
	conn, err := grpc.Dial("tsn-service:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed dialing tsn-service: %v", err)
		return nil, err
	}

	defer conn.Close()

	client := notification.NewNotificationClient(conn)

	confId, err := client.CalcConfig(context.Background(), reqIds)
	if err != nil {
		log.Errorf("Calculating configuration failed: %v", err)
		return nil, err
	}

	return confId, nil
}

func applyConfiguration(id *notification.UUID) error {

	/************************ APPLY CONFIG ************************/
	client, err := connectToGnmiService("onos-config:5150")
	if err != nil {
		log.Errorf("Failed connecting to gNMI service: %v", err)
		return err
	}

	confReq := getSetRequestForConfig()

	response, err := client.(*gclient.Client).Set(context.Background(), confReq)
	if err != nil {
		log.Errorf("Target returned RPC error for Set: %v", err)
		return err
	}

	log.Infof("Response from device-monitor is: %v", response)

	return nil
}

// Creates a set request for applying a new configuration.
func getSetRequestForConfig() *pb.SetRequest {
	// TODO: Generate all pb.Update objects for all the values that should be changed... (Let Hamza know
	// that the current implementation of onos-config only takes in config and not a config ID).
	confSetRequest := pb.SetRequest{
		Update: []*pb.Update{ // List of updated values for the configuration
			{
				Path: &pb.Path{
					Target: "192.168.0.1",
					Elem:   []*pb.PathElem{}, // Path to an element that should be updated
				},
			},
		},
		Extension: []*gnmi_ext.Extension{
			{
				Ext: &gnmi_ext.Extension_RegisteredExt{
					RegisteredExt: &gnmi_ext.RegisteredExtension{
						Id:  gnmi_ext.ExtensionID(100),
						Msg: []byte("my_network_change"),
					},
				},
			},
			{
				Ext: &gnmi_ext.Extension_RegisteredExt{
					RegisteredExt: &gnmi_ext.RegisteredExtension{
						Id:  gnmi_ext.ExtensionID(101),
						Msg: []byte("1.0.2"),
					},
				},
			},
			{
				Ext: &gnmi_ext.Extension_RegisteredExt{
					RegisteredExt: &gnmi_ext.RegisteredExtension{
						Id:  gnmi_ext.ExtensionID(102),
						Msg: []byte("tsn-model"),
					},
				},
			},
		},
	}

	return &confSetRequest
}

// Takes in addr such as "onos-config:5150" and returns a gNMI-client.
func connectToGnmiService(addr string) (client.Impl, error) {
	cert, err := tls.X509KeyPair([]byte(certs.DefaultClientCrt), []byte(certs.DefaultClientKey))
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	client, err := gclient.New(context.Background(), client.Destination{
		Addrs:       []string{addr},
		Target:      strings.Split(addr, ":")[0],
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         tlsConfig,
	})

	if err != nil {
		log.Errorf("Failed creating gNMI client to onos-config: %v", err)
		return nil, err
	}

	return client, nil
}
