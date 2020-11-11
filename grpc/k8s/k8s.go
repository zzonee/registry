package k8s

import (
	"github.com/zzonee/registry"
	util "github.com/zzonee/registry/grpc"
	"github.com/zzonee/registry/k8s"
	"google.golang.org/grpc/resolver"
	"log"
	"sync"
	"time"
)

// grpc.Dial("k8s://namespace/servicename:portname")
// grpc.Dial("k8s://namespace/servicename:port")
// grpc.Dial("k8s:///servicename:portname") // namespace=default
// grpc.Dial("k8s:///servicename:port")
type k8sBuilder struct {
	sync.Mutex
	rs map[string]registry.Discovery
}

func RegisterBuilder() error {
	b := &k8sBuilder{
		rs: map[string]registry.Discovery{},
	}
	resolver.Register(b)
	return nil
}

func (b *k8sBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var (
		err       error
		discovery registry.Discovery
		namesapce = target.Authority
	)
	if namesapce == "" {
		namesapce = "default"
	}
	if discovery, err = b.getDiscovery(namesapce); err != nil {
		return nil, err
	}
	ch, err := discovery.Discover(target.Endpoint)
	if err != nil {
		return nil, err
	}

	select {
	case inss := <-ch:
		util.UpdateAddress(inss, cc)
	case <-time.After(time.Minute):
		log.Printf("not resolve succuss in one minute, target:%v", target)
	}

	go func() {
		for inss := range ch {
			util.UpdateAddress(inss, cc)
		}
	}()
	return &util.NoopResolver{}, nil
}

func (b *k8sBuilder) getDiscovery(namespace string) (r registry.Discovery, err error) {
	b.Lock()
	defer b.Unlock()
	if r = b.rs[namespace]; r != nil {
		return
	}
	if r, err = k8s.NewDiscovery(namespace); err != nil {
		return
	}
	b.rs[namespace] = r
	return
}

func (b *k8sBuilder) Scheme() string {
	return "k8s"
}
