package grpc

import (
	"github.com/zzonee/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

func UpdateAddress(inss []registry.Instance, conn resolver.ClientConn) {
	var address []resolver.Address
	for _, ins := range inss {
		address = append(address, resolver.Address{Addr: ins.Address})
	}
	conn.UpdateState(resolver.State{
		Addresses: address,
	})
}

type NoopResolver struct{}

func (r *NoopResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *NoopResolver) Close()                                {}
