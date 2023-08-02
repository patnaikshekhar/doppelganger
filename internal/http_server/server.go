package http_server

import (
	"doppelganger/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type HttpServer struct {
	ForwardedServices *services.ForwardedServices
}

func New(fwdServices *services.ForwardedServices) *HttpServer {
	return &HttpServer{
		ForwardedServices: fwdServices,
	}
}

func (s *HttpServer) Start() error {
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		s.ForwardedServices.RLock()
		defer s.ForwardedServices.RUnlock()

		servicesBytes, err := json.Marshal(s.ForwardedServices)
		if err != nil {
			log.Printf("Unable to marshal forwarded services: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(servicesBytes))
	})

	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		go os.Exit(0)
	})

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		return err
	}

	return nil
}
