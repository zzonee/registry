package direct

import (
	"fmt"
	"github.com/zzonee/registry"
	"strings"
)

type directDiscovery struct{}

func NewDiscovery() registry.Discovery {
	return &directDiscovery{}
}

// target 格式： "127.0.0.1:8000,127.0.0.1:8001"
func (d *directDiscovery) Discover(target string) (<-chan []registry.Instance, error) {
	endpoints := strings.Split(target, ",")
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no endpoint")
	}
	var inss []registry.Instance
	for _, addr := range endpoints {
		ins := registry.Instance{Address: addr}
		inss = append(inss, ins)
	}

	ch := make(chan []registry.Instance)
	go func() {
		ch <- inss
	}()
	return ch, nil
}

func (d *directDiscovery) Close() {}
