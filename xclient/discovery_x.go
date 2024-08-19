package xclient

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type XRegistryDiscovery struct {
	*MultiServersDiscovery
	registry   string
	timeout    time.Duration
	lastUpdate time.Time
}

const defaultUpdateTimeout = time.Second * 10

func NewXRegistryDiscovery(registryAddr string, timeout time.Duration) *XRegistryDiscovery {
	if timeout == 0 {
		timeout = defaultUpdateTimeout
	}
	return &XRegistryDiscovery{
		MultiServersDiscovery: NewMultiServersDiscovery(make([]string, 0)),
		registry:              registryAddr,
		timeout:               timeout,
	}
}

func (d *XRegistryDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	d.lastUpdate = time.Now()
	return nil
}

func (d *XRegistryDiscovery) Refresh() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.lastUpdate.Add(d.timeout).After(time.Now()) {
		return nil
	}
	log.Println("rpc registry: refresh servers from registry", d.registry)
	resp, err := http.Get(d.registry)
	if err != nil {
		log.Println("rpc registry refresh error", err)
		return err
	}
	servers := strings.Split(resp.Header.Get("X-Xrpc-Servers"), ",")
	d.servers = make([]string, 0, len(servers))
	for _, server := range servers {
		t := strings.TrimSpace(server)
		if t != "" {
			d.servers = append(d.servers, t)
		}
	}
	d.lastUpdate = time.Now()
	return nil
}

func (d *XRegistryDiscovery) Get(mode SelectMode) (string, error) {
	if err := d.Refresh(); err != nil {
		return "", err
	}
	return d.MultiServersDiscovery.Get(mode)
}

func (d *XRegistryDiscovery) GetAll() ([]string, error) {
	if err := d.Refresh(); err != nil {
		return nil, err
	}
	return d.MultiServersDiscovery.GetAll()
}
