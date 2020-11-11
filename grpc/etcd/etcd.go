package etcd

import (
	"github.com/zzonee/registry"
	"github.com/zzonee/registry/etcd"
	util "github.com/zzonee/registry/grpc"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

// grpc.Dial("etcd://default/servicename")
type etcdBuilder struct {
	discovery registry.Discovery
}

func RegisterBuilder(conf clientv3.Config) error {
	d, err := etcd.NewDiscovery(conf)
	if err != nil {
		return err
	}
	b := &etcdBuilder{
		discovery: d,
	}

	resolver.Register(b)
	return nil
}

func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ch, err := b.discovery.Discover(target.Endpoint)
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

func (b *etcdBuilder) Scheme() string {
	return "etcd"
}
