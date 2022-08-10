package storewrapper

import (
	"context"
	"crypto/tls"
	"io"
	"main-service/pkg/structures/configuration"
	"strings"

	"github.com/atomix/atomix-go-client/pkg/atomix"
	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/errors"

	// "github.com/onosproject/helmit/pkg/helm"
	// "github.com/onosproject/helmit/pkg/kubernetes"
	// v1 "github.com/onosproject/helmit/pkg/kubernetes/core/v1"

	// "github.com/onosproject/onos-lib-go/pkg/certs"
	// _map "github.com/atomix/atomix-go-client/pkg/atomix/map"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"
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

func getFromStore(urn string) (*configuration.ConfigResponse, error) {
	ctx := context.Background()

	// Create a slice of maximum two URN elements
	urnElems := strings.SplitN(urn, ".", 2)

	// log.Info("Getting map...")

	// Get the store
	store, err := atomix.GetMap(ctx, urnElems[0])
	if err != nil {
		log.Errorf("Failed getting store \"%s\": %v", urnElems[0], err)
		return &configuration.ConfigResponse{}, err
	}

	// log.Info("Getting obj from store...")

	// TODO: Check if the URN contains more complex path and do something special then

	// Get the object from store
	obj, err := store.Get(ctx, urnElems[1])
	if err != nil {
		log.Errorf("Failed getting resource \"%s\": %v", urnElems[1], err)
		return &configuration.ConfigResponse{}, err
	}

	// log.Info("Unmarshaling object...")

	// Unmarshal the byte slice from the store into request data
	var req = configuration.ConfigResponse{}
	err = proto.Unmarshal(obj.Value, &req)
	if err != nil {
		log.Errorf("Failed to unmarshal request data from store: %v", err)
		return &configuration.ConfigResponse{}, nil
	}

	return &req, nil
}

func getNetworkChangeFromStore(urn string) ([]byte, error) {
	ctx := context.Background()

	// Create a slice of maximum two URN elements
	urnElems := strings.SplitN(urn, ".", 2)

	// Get the store
	store, err := atomix.GetMap(ctx, urnElems[0])
	if err != nil {
		log.Errorf("Failed getting store \"%s\": %v", urnElems[0], err)
		return nil, err
	}

	// TODO: Check if the URN contains more complex path and do something special then???

	// Get the object from store
	obj, err := store.Get(ctx, urnElems[1])
	if err != nil {
		log.Errorf("Failed getting resource \"%s\": %v", urnElems[1], err)
		return nil, err
	}

	// // Unmarshal the byte slice from the store into gNMI set request
	// var req = &pb.SetRequest{}
	// err = proto.Unmarshal(obj.Value, req)
	// if err != nil {
	// 	log.Errorf("Failed to unmarshal request data from store: %v", err)
	// 	return nil, err
	// }

	return obj.Value, nil
}

// ONLY FOR TESTING
func getTopoFromStore() {
	ctx := context.Background()

	cert, err := tls.X509KeyPair([]byte(certs.DefaultClientCrt), []byte(certs.DefaultClientKey))
	if err != nil {
		return
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	conn, err := grpc.Dial("onos-topo:5150", grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		log.Fatalf("Failed dialing onos-topo: %v", err)
		return
	}

	// defer conn.Close()

	client := topo.CreateTopoClient(conn)
	// resp, err := client.Create(ctx, &topo.CreateRequest{Object: &topo.Object{UUID: "123-456", Type: topo.Object_ENTITY, Obj: &topo.Object_Entity{Entity: &topo.Entity{}}}})
	resp, err := client.List(ctx, &topo.ListRequest{})
	if err != nil {
		log.Fatalf("Failed listing topo object: %v", errors.FromGRPC(err))
		return
	}

	log.Infof("Topo objects: %v", resp.Objects)

	watchClient, err := client.Watch(ctx, &topo.WatchRequest{Noreplay: false})
	if err != nil {
		log.Fatalf("Failed to watch topo for updates: %v", errors.FromGRPC(err))
		return
	}

	go func() {
		for {
			resp, err := watchClient.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Warn(err)
				break
			}
			log.Infof("Event: %v", resp.Event)
		}
	}()

	// log.Infof("Response: %v", resp.Objects)

	// // Get the store
	// store, err := atomix.GetMap(ctx, "onos-topo-objects")
	// if err != nil {
	// 	log.Errorf("Failed getting store: %v", err)
	// 	return
	// }

	// // Get all objects in store
	// eventChannel := make(chan _map.Event)
	// err = store.Watch(ctx, eventChannel)
	// // err = store.Entries(ctx, entryChannel)
	// if err != nil {
	// 	log.Errorf("Failed getting entries: %v", err)
	// 	return
	// }

	// go func() {
	// 	for {
	// 		select {
	// 		case event := <-eventChannel:
	// 			log.Infof("Event: %v", event)
	// 		}
	// 	}
	// }()
	//

	//

	//

	//

	// conn := getConn(ctx)
	// if conn == nil {
	// 	log.Error("Failed getting conn!")
	// 	return
	// }

	// defer conn.Close()

	// client := topo.NewTopoClient(conn)

	// // filters := &topo.Filters{
	// // 	KindFilter: &topo.Filter{
	// // 		Filter: &topo.Filter_In{In: &topo.InFilter{Values: []string{topo.E2NODE, topo.E2CELL}}},
	// // 	},
	// // }

	// resp, err := client.List(ctx, &topo.ListRequest{SortOrder: topo.SortOrder_ASCENDING})
	// if err != nil {
	// 	log.Fatalf("Failed getting topo list: %v", err)
	// 	return
	// }

	// log.Infof("Response: %v", resp)
}

// func getConn(ctx context.Context) *grpc.ClientConn {
// 	release := helm.Chart("open-cnc").Release("open-cnc")
// 	kubeClient, err := kubernetes.NewForRelease(release)
// 	if err != nil {
// 		log.Errorf("Failed getting new kubernetes client: %v", err)
// 		return nil
// 	}
// 	pods := kubeClient.CoreV1().Pods()
// 	if err != nil {
// 		log.Errorf("Failed getting new v1 client: %v", err)
// 		return nil
// 	}

// 	// podList, err := pods.List(ctx)
// 	// if err != nil {
// 	// 	log.Errorf("Failed getting pod-list: %v", err)
// 	// 	return nil
// 	// }

// 	log.Infof("v1 podList: %v", pods)

// 	// pods, err := kubeClient.CoreV1().Pods().List(context.Background())
// 	// if err != nil {
// 	// 	log.Errorf("Failed getting pod list: %v", err)
// 	// 	return nil
// 	// }

// 	// log.Infof("KubeClient: %v", pods)

// 	// releases := helm.Client().Namespace("open-cnc")

// 	// log.Infof("Releases: %v", releases)

// 	return nil

// 	// service := getService(release, ctx)
// 	// if service == nil {
// 	// 	log.Error("Failed getting service!")
// 	// 	return nil
// 	// }

// 	// cert, err := tls.X509KeyPair([]byte(certs.DefaultClientCrt), []byte(certs.DefaultClientKey))
// 	// if err != nil {
// 	// 	log.Errorf("Failed generating certs: %v", err)
// 	// 	return nil
// 	// }

// 	// // conn, err := grpc.Dial("onos-topo:5150", grpc.WithInsecure())
// 	// conn, err := grpc.Dial(service.Ports()[0].Address(true), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
// 	// 	Certificates:       []tls.Certificate{cert},
// 	// 	InsecureSkipVerify: true,
// 	// })))
// 	// if err != nil {
// 	// 	log.Fatalf("Failed dialing onos-topo: %v", err)
// 	// 	return nil
// 	// }

// 	// return conn
// }

// func getService(release *helm.HelmRelease, ctx context.Context) *v1.Service {
// 	releaseClient := kubernetes.NewForReleaseOrDie(release)
// 	nsList, err := releaseClient.CoreV1().Namespaces().List(ctx)
// 	if err != nil {
// 		log.Errorf("Failed to get namespaces: %v", err)
// 		return nil
// 	}

// 	log.Infof("Namespaces: %v", nsList)

// 	service, err := releaseClient.CoreV1().Services().Get(ctx, "onos-topo")
// 	if err != nil {
// 		log.Errorf("Failed to get service onos-topo: %v", err)
// 		return nil
// 	}

// 	return service
// }
