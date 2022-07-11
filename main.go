package main

import (
	"context"
	"crypto/tls"
	"main-service/pkg/logger"
	northboundinterface "main-service/pkg/northbound-interface"
	store "main-service/pkg/store-wrapper"

	"github.com/google/uuid"
	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/errors"
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
	ctx := context.Background()

	cert, err := tls.X509KeyPair([]byte(certs.DefaultClientCrt), []byte(certs.DefaultClientKey))
	if err != nil {
		log.Errorf("Failed generating tls certs: %v", err)
		return
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial("onos-topo:5150", opts...)
	if err != nil {
		log.Fatalf("Failed dialing onos-topo: %v", err)
		return
	}

	defer conn.Close()

	client := topo.CreateTopoClient(conn)

	obj := &topo.Object{
		UUID:     topo.UUID(uuid.New().String()),
		ID:       topo.ID("0"),
		Revision: topo.Revision(123),
		Type:     topo.Object_ENTITY,
		Obj: &topo.Object_Entity{
			Entity: &topo.Entity{},
		},
	}

	obj.SetAspectBytes("onos.topo.Configurable", []byte(`{"address": "192.168.0.1", "version": "1.0.0", "type": "my-model"}`))
	obj.SetAspectBytes("onos.topo.TLSOptions", []byte(`{"insecure": true, "plain": true}`))
	obj.SetAspectBytes("onos.topo.Asset", []byte(`{"name": "switch-0"}`))
	obj.SetAspectBytes("onos.topo.MastershipState", []byte(`{}`))

	resp, err := client.Create(ctx, &topo.CreateRequest{Object: obj})
	if err != nil {
		log.Fatalf("Failed creating topo object: %v", errors.FromGRPC(err))
		return
	}

	log.Infof("onos-topo create response: %v", resp)
}
