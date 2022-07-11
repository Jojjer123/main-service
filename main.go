package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"main-service/pkg/logger"
	northboundinterface "main-service/pkg/northbound-interface"
	store "main-service/pkg/store-wrapper"

	"github.com/google/uuid"
	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/certs"

	// "github.com/onosproject/onos-lib-go/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log = logger.GetLogger()

func main() {
	// Temporarily add switches to onos-topo
	addSwitches()

	// Create TSN stores
	store.CreateStores()

	// Start UNI server
	northboundinterface.StartServer()
}

func addSwitches() {
	/************************ CREATE KIND ************************/
	if err := createKind("netconf-device"); err != nil {
		log.Errorf("Failed creating kind: %v", err)
		return
	}

	/************************ CREATE DEVICE ************************/
	if err := createDevice("switch-0", "192.168.0.1", "netconf-device", "Devicesim", "1.0.0"); err != nil {
		log.Errorf("Failed creating device: %v", err)
		return
	}
}

func connectToGrpcService(addr string) (*grpc.ClientConn, error) {
	cert, err := tls.X509KeyPair([]byte(certs.DefaultClientCrt), []byte(certs.DefaultClientKey))
	if err != nil {
		log.Errorf("Failed generating tls certs: %v", err)
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Errorf("Failed dialing %s: %v", addr, err)
		return nil, err
	}

	return conn, nil
}

func createKind(name string) error {
	conn, err := connectToGrpcService("onos-topo:5150")
	if err != nil {
		log.Errorf("Failed connecting to service: %v", err)
		return err
	}

	defer conn.Close()

	client := topo.CreateTopoClient(conn)

	req := &topo.CreateRequest{
		Object: &topo.Object{
			UUID:     topo.UUID(uuid.NewString()),
			ID:       topo.ID(name),
			Revision: topo.Revision(5),
			Type:     topo.Object_KIND,
			Obj: &topo.Object_Kind{
				Kind: &topo.Kind{
					Name: name,
				},
			},
		},
	}

	resp, err := client.Create(context.Background(), req)
	if err != nil {
		log.Errorf("Failed creating kind: %v", err)
		return err
	}

	log.Infof("Created kind %v", resp)

	return nil
}

// TODO: Add dynamic src and target IDs whatever they are...
func createDevice(name string, addr string, kind string, model string, modelVersion string) error {
	conn, err := connectToGrpcService("onos-topo:5150")
	if err != nil {
		log.Errorf("Failed connecting to service: %v", err)
		return err
	}

	defer conn.Close()

	client := topo.CreateTopoClient(conn)

	obj := &topo.Object{
		UUID:     topo.UUID(uuid.NewString()),
		ID:       topo.ID(name),
		Revision: topo.Revision(2),
		Type:     topo.Object_ENTITY,
		Obj: &topo.Object_Entity{
			Entity: &topo.Entity{
				KindID: topo.ID(kind),
			},
		},
	}

	obj.SetAspectBytes("onos.topo.Configurable", []byte(fmt.Sprintf(`{"address": "%s", "version": "%s", "type": "%s"}`, addr, modelVersion, model)))
	obj.SetAspectBytes("onos.topo.TLSOptions", []byte(`{"insecure": true, "plain": true}`))
	obj.SetAspectBytes("onos.topo.Asset", []byte(fmt.Sprintf(`{"name": "%v"}`, name)))
	obj.SetAspectBytes("onos.topo.MastershipState", []byte(`{}`))

	req := &topo.CreateRequest{
		Object: obj,
	}

	resp, err := client.Create(context.Background(), req)
	if err != nil {
		log.Errorf("Failed creating device: %v", err)
		return err
	}

	log.Infof("Created device %v", resp)

	return nil
}
