package warden

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Warden struct {
	port  uint32
	hosts map[string]uint32
}

func New(port uint32) *Warden {
	return &Warden{
		port:  port,
		hosts: make(map[string]uint32),
	}
}

func (w *Warden) Add(host string, localport uint32) {
	w.hosts[host] = localport
}

func (w *Warden) Start() error {
	wardenServer := http.Server{
		Addr:    fmt.Sprintf(":%d", w.port),
		Handler: w,
	}

	err := wardenServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (w *Warden) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	h := strings.Split(req.Host, ":")[0]
	localport, ok := w.hosts[h]

	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	backend := fmt.Sprintf("http://localhost:%d", localport)
	backendUrl, err := url.Parse(backend)

	if err != nil {
		log.Printf("Error occured while parsing backend URL %s", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(backendUrl)
	proxy.ServeHTTP(res, req)
}
