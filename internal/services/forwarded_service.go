package services

import "sync"

type ForwardedService struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	LocalPort   uint32 `json:"localport"`
	ServicePort uint32 `json:"serviceport"`
}

type ForwardedServices struct {
	Services []ForwardedService
	sync.RWMutex
}
