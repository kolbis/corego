package servicediscovery

import (
	"time"

	"github.com/go-kit/kit/sd/dnssrv"
	tlelogger "github.com/kolbis/corego/logger"
	"github.com/kolbis/corego/utils"
)

const (
	// DefaultTTL is 30 seconds
	DefaultTTL time.Duration = time.Second * 30
)

// DNSServiceDiscovery provide service discovery capabilities for consul and for DNS (k8s)
type DNSServiceDiscovery struct {
	Logger      tlelogger.Logger
	DNSIntances map[string]*dnssrv.Instancer
}

// NewDNSServiceDiscovery creates a new instance of the service directory
func NewDNSServiceDiscovery(logger tlelogger.Logger) DNSServiceDiscovery {
	sd := DNSServiceDiscovery{
		Logger:      logger,
		DNSIntances: map[string]*dnssrv.Instancer{},
	}
	return sd
}

// DNSInstance will return DNS instancer which will be used to lookup a DNS service
// It will cache the instances
func (sd *DNSServiceDiscovery) DNSInstance(serviceName string) *dnssrv.Instancer {
	key := utils.NewKeys().Build("dns", serviceName)

	var instancer *dnssrv.Instancer = sd.DNSIntances[key]
	if instancer != nil {
		return instancer
	}

	instancer = dnssrv.NewInstancer(serviceName, DefaultTTL, sd.Logger)
	sd.DNSIntances[key] = instancer

	return instancer
}
