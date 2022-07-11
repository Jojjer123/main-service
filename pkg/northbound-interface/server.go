package northboundinterface

import (
	// "bytes"
	// "encoding/json"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	// "io"
	// "io/ioutil"
	"net/http"
	"reflect"

	// "strings"
	"time"

	handler "main-service/pkg/event-handler"
	"main-service/pkg/logger"

	// store "main-service/pkg/store-wrapper"
	"main-service/pkg/structures/configuration"

	"github.com/go-openapi/runtime/middleware/header"
	"github.com/gogo/protobuf/jsonpb"
	// "google.golang.org/protobuf/encoding/protojson"
)

const PORT uint16 = 8080

var log = logger.GetLogger()

func StartServer() {
	http.HandleFunc("/get_config", getConfig)
	http.HandleFunc("/update_stream", updateStream)
	http.HandleFunc("/remove_stream", removeStream)
	http.HandleFunc("/join_stream", joinStream)
	http.HandleFunc("/leave_stream", leaveStream)

	log.Infof("API endpoint -> http://localhost:%d/get_config", PORT)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil); err != nil {
		log.Errorf("Failed to listen and server on %d, with error: %+v", PORT, err)
	}
}

func getConfig(writer http.ResponseWriter, req *http.Request) {
	timeOfReq := time.Now()

	if err := checkHeader(req); err != nil {
		http.Error(writer, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	// body, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	log.Errorf("Failed reading body: %v", err)
	// }

	var configRequest configuration.ConfigRequest

	// ioReader := ioutil.NopCloser(bytes.NewBuffer(body))

	err := jsonpb.Unmarshal(req.Body, &configRequest)

	// ioReader = ioutil.NopCloser(bytes.NewBuffer(body))

	// dec := json.NewDecoder(ioReader)
	// dec.DisallowUnknownFields()

	// var testReq configuration.ConfigRequest

	// err = dec.Decode(&testReq)

	// // NEED A REMAKE TO SUIT PROTO UMARSHSALING ERRORS
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnsupportedTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(writer, msg, http.StatusBadRequest)
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			http.Error(writer, msg, http.StatusBadRequest)
		case errors.As(err, &unmarshalTypeError):
			// msg := fmt.Sprintf("Request body contains invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			msg := "Request body contains invalid structure"
			http.Error(writer, msg, http.StatusBadRequest)
		case strings.HasPrefix(err.Error(), "json:unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json:unknown field")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(writer, msg, http.StatusBadRequest)
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(writer, msg, http.StatusBadRequest)
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(writer, msg, http.StatusRequestEntityTooLarge)
		default:
			log.Error(err.Error())
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if err != nil {
		log.Errorf("Failed to read req.Body: %v", err)
		http.Error(writer, "Failed to read body", http.StatusBadRequest)
		return
	}

	log.Infof("%+v", reflect.TypeOf(req.Body))

	// err = dec.Decode(&struct{}{})
	// if err != io.EOF {
	// 	msg := "Request body must only contain a single JSON object"
	// 	http.Error(writer, msg, http.StatusBadRequest)
	// 	return
	// }

	// Call handler to deal with addStream request
	confId, err := handler.HandleAddStreamEvent(&configRequest, timeOfReq)
	// _, err = handler.HandleAddStreamEvent(&configRequest, timeOfReq)
	if err != nil {
		log.Errorf("Failed handling event: %v", err)
		http.Error(writer, "Error in request???", http.StatusBadRequest)
		return
	}

	// log.Info("Handled event!")

	// Write configRequest back to client
	// fmt.Fprintf(writer, "request: %+v", configRequest)

	// TODO: BUILD RESPONSE (SIMULATE BUILD FOR NOW, NEED TO MOVE FROM TSN-SERVICE TO HERE LATER ON)
	// log.Infof("confId.GetValue(): %v", confId.GetValue())
	// confIdString := fmt.Sprintf("%v", confId.GetValue())
	// confData, err := store.GetResponseData(confIdString)
	// if err != nil {
	// 	log.Errorf("Failed getting response data: %v", err)
	// 	return
	// }

	// log.Infof("confData: %v", confData)

	// data, err := protojson.Marshal(confData)
	// if err != nil {
	// 	log.Errorf("Failed marshaling confData: %v", err)
	// 	return
	// }

	resp, err := createResponse(confId)
	if err != nil {
		log.Errorf("Failed to create UNI response!")
		return
	}

	// log.Infof("Data: %v", data)

	writer.Header().Add("Content-Type", "application/json; charset=utf-8")
	writer.Write(resp)
}

func updateStream(writer http.ResponseWriter, req *http.Request) {
	if err := checkHeader(req); err != nil {
		http.Error(writer, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	// var updateRequest stream.updateRequest

	// err := jsonpb.Unmarshal(req.Body, &updateRequest)
	// if err != nil {
	// 	log.Errorf("Failed to read req.Body: %v", err)
	// 	http.Error(writer, "Failed to read body", http.StatusBadRequest)
	// 	return
	// }

	// log.Infof("%+v", reflect.TypeOf(req.Body))

	writer.Write([]byte("Done!"))
}

func removeStream(writer http.ResponseWriter, req *http.Request) {
	if err := checkHeader(req); err != nil {
		http.Error(writer, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
}

func joinStream(writer http.ResponseWriter, req *http.Request) {
	if err := checkHeader(req); err != nil {
		http.Error(writer, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
}

func leaveStream(writer http.ResponseWriter, req *http.Request) {
	if err := checkHeader(req); err != nil {
		http.Error(writer, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
}

func checkHeader(req *http.Request) error {
	if req.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(req.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return errors.New(msg)
		}
	}

	return nil
}
