package dns

import "doppelganger/internal/services"

type DNS interface {
	Add(services.ForwardedService) error
}
