package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"main-service/pkg/logger"
	northboundinterface "main-service/pkg/northbound-interface"
	store "main-service/pkg/store-wrapper"
	monitor "main-service/pkg/structures/temp-monitor-conf"

	"github.com/atomix/atomix-go-client/pkg/atomix"
	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/certs"

	// "github.com/onosproject/onos-lib-go/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log = logger.GetLogger()

func main() {
	time.Sleep(1 * time.Minute)
	// Temporarily add switches to onos-topo
	addSwitches()
	// Temporarily add monitor configs & adapter to atomix
	addMonitorConf()

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
	if err := createDevice("switch-0", "gnmi-netconf-adapter:11161", "netconf-device", "tsn-model", "1.0.2"); err != nil {
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

	obj := topo.Object{
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

	var configurable = topo.Configurable{
		Type:    model,
		Address: addr,
		Target:  "192.168.0.1",
		Version: modelVersion,
		Timeout: uint64(10 * time.Second),
	}

	m := jsonpb.Marshaler{}

	configData, err := m.MarshalToString(&configurable)
	if err != nil {
		log.Errorf("Failed marshaling configurable: %v", err)
		return err
	}

	obj.SetAspectBytes("onos.topo.Configurable", []byte(configData))
	obj.SetAspectBytes("onos.topo.TLSOptions", []byte(`{"insecure": true, "plain": true}`))
	obj.SetAspectBytes("onos.topo.Asset", []byte(fmt.Sprintf(`{"name": "%v"}`, name)))
	obj.SetAspectBytes("onos.topo.MastershipState", []byte(`{}`))

	req := &topo.CreateRequest{
		Object: &obj,
	}

	resp, err := client.Create(context.Background(), req)
	if err != nil {
		log.Errorf("Failed creating device: %v", err)
		return err
	}

	log.Infof("Created device %v", resp)

	return nil
}

func addMonitorConf() {
	fileContent, err := ioutil.ReadFile("monitor-conf-example.yaml")
	if err != nil {
		log.Errorf("Failed reading file: %v", err)
		return
	}

	// log.Infof("Read file: %v", fileContent)

	jsonBytes, err := yaml.YAMLToJSON(fileContent)
	if err != nil {
		log.Errorf("Failed converting file content from yaml to json: %v", err)
		return
	}

	var conf = &monitor.Config{}
	if err = jsonpb.Unmarshal(bytes.NewReader(jsonBytes), conf); err != nil {
		log.Errorf("Failed unmarshaling json to protobuf: %v", err)
		return
	}

	// log.Infof("Umarshaled into: %v", conf)

	rawConf, err := proto.Marshal(conf)
	if err != nil {
		log.Errorf("Failed marshaling config to byte slice: %v", err)
		return
	}

	// log.Infof("Marshaled into: %v", rawConf)

	sendToStore(rawConf, "configurations.monitor-config.192.168.0.1")

	data, err := proto.Marshal(&monitor.Adapter{
		Protocol: "NETCONF",
		Address:  "gnmi-netconf-adapter",
	})
	if err != nil {
		log.Errorf("Failed marshaling adapter: %v", err)
		return
	}

	sendToStore(data, "configurations.adapter.NETCONF")
}

func sendToStore(obj []byte, urn string) error {
	ctx := context.Background()

	// Create a slice of URN elements
	urnElems := strings.SplitN(urn, ".", 2)

	// Get the store
	store, err := atomix.GetMap(ctx, urnElems[0])
	if err != nil {
		log.Errorf("Failed getting store \"%s\": %v", urnElems[0], err)
		return err
	}

	// TODO: Check if the URN contains more complex path and do something special then

	// Store the object
	_, err = store.Put(ctx, urnElems[1], obj)
	if err != nil {
		log.Errorf("Failed storing resource \"%s\": %v", urnElems[1], err)
		return err
	}

	return nil
}
