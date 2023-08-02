package proxy

import (
	"doppelganger/internal/services"
	"doppelganger/internal/warden"
	"fmt"
	"log"
)

type MultiPortProxy struct {
	IncomingChannel chan services.ForwardedService
	Proxies         map[uint32]*warden.Warden
}

func NewMultiPortProxy(ic chan services.ForwardedService) *MultiPortProxy {
	return &MultiPortProxy{
		IncomingChannel: ic,
		Proxies:         make(map[uint32]*warden.Warden),
	}
}

func (p *MultiPortProxy) Start() {
	for event := range p.IncomingChannel {
		_, ok := p.Proxies[event.ServicePort]
		if !ok {
			log.Printf("Starting warden on port %d", event.ServicePort)
			w := warden.New(event.ServicePort)
			w.Add(fmt.Sprintf("%s.%s", event.Name, event.Namespace), event.LocalPort)
			p.Proxies[event.ServicePort] = w
			go func(event services.ForwardedService) {
				err := w.Start()
				if err != nil {
					log.Printf("Error starting warden on port %d", event.ServicePort)
				}
			}(event)
		} else {
			w := p.Proxies[event.ServicePort]
			w.Add(fmt.Sprintf("%s.%s", event.Name, event.Namespace), event.LocalPort)
		}
	}
}

/*
nginx.ingress 80 -> 30001
postgres.default 6789 -> 30002
myservice.default 80 -> 30003

80 -> Warden(80) -> nginx.ingress (30001) / myservice.default (30003)
6789 -> Warden(6789) -> postgres.default (30002)
*/
