package direct

import (
	"github.com/zzonee/registry"
	"github.com/zzonee/registry/direct"
	util "github.com/zzonee/registry/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

func init() {
	resolver.Register(&directBuilder{
		discovery: direct.NewDiscovery(),
	})
}

func RegisterBuilder() error {
	return nil
}

// grpc.Dial("direct://default/127.0.0.1:8000,127.0.0.1:8001")
type directBuilder struct {
	discovery registry.Discovery
}

func (b *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
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
	return &util.NoopResolver{}, nil
}

func (b *directBuilder) Scheme() string {
	return "direct"
}
