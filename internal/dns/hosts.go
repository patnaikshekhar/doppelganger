package dns

import (
	"doppelganger/internal/services"
	"fmt"
	"os"
	"strings"
)

type Host struct {
}

func NewHosts() *Host {
	return &Host{}
}

func (h *Host) Add(service services.ForwardedService) error {
	hostBytes, err := os.ReadFile("/etc/hosts")
	if err != nil {
		return err
	}

	serviceName := fmt.Sprintf("127.0.0.1 %s.%s #DG", service.Name, service.Namespace)

	if strings.Contains(string(hostBytes), serviceName) {
		return nil
	}

	contents := string(hostBytes) + serviceName + "\n"
	err = os.WriteFile("/etc/hosts", []byte(contents), 0655)
	if err != nil {
		return err
	}

	return nil

}
