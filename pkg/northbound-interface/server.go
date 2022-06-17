package northboundinterface

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"main-service/pkg/logger"
	southbound "main-service/pkg/southbound-interface"
	"main-service/pkg/structures"

	"github.com/go-openapi/runtime/middleware/header"
)

const PORT uint16 = 8080

var log = logger.GetLogger()

func StartServer() {
	http.HandleFunc("/get_config", getConfig)

	log.Infof("API endpoint -> http://localhost:%d/get_config", PORT)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil); err != nil {
		log.Errorf("Failed to listen and server on %d, with error: %+v", PORT, err)
	}
}

func getConfig(writer http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(req.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(writer, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	var configRequest structures.ConfigRequest
	err := dec.Decode(&configRequest)

	log.Infof("%+v", reflect.TypeOf(req.Body))

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

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		http.Error(writer, msg, http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(configRequest)
	if err != nil {
		panic(err)
	}

	// Call southbound to store the "request"
	southbound.StoreRequestInStorage(data)

	// Call the southbound to get a "response"
	if data = southbound.GetConfigFromStorage(); data == nil {
		log.Error("Received no data from storage")
		return
	}

	//Write configRequest back to client
	// fmt.Fprintf(writer, "request: %+v", configRequest)

	writer.Header().Add("Content-Type", "application/json; charset=utf-8")
	writer.Write(data)
}
