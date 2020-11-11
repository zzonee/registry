package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zzonee/registry"
	"go.etcd.io/etcd/clientv3"
	"log"
)

type etcdDiscovery struct {
	client *clientv3.Client
}

func NewDiscovery(conf clientv3.Config) (registry.Discovery, error) {
	d := &etcdDiscovery{}
	c, err := clientv3.New(conf)
	if err != nil {
		return nil, err
	}
	d.client = c
	return d, nil
}

// target: serviceName
func (d *etcdDiscovery) Discover(target string) (<-chan []registry.Instance, error) {
	ch := make(chan []registry.Instance)
	go d.watch(ch, target)
	return ch, nil
}

func (d *etcdDiscovery) watch(ch chan<- []registry.Instance, serviceName string) {
	prefix := fmt.Sprintf("/%s/%s/", etcdPrefix, serviceName)

	get := func() []registry.Instance {
		resp, err := d.client.Get(context.Background(), prefix, clientv3.WithPrefix())
		if err != nil {
			log.Printf("etcd discovery watch err:%v, servicename:%s", err, serviceName)
			return nil
		}
		var insss []registry.Instance
		for _, kv := range resp.Kvs {
			ins := registry.Instance{}
			if err = json.Unmarshal(kv.Value, &ins); err == nil {
				insss = append(insss, ins)
			} else {
				log.Printf("etcd discovery watch unmarshal err:%v, servicename:%s", err, serviceName)
			}
		}
		return insss
	}
	if inss := get(); len(inss) > 0 {
		ch <- inss
	}

	eventch := d.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for range eventch {
		ch <- get()
	}
	return
}

func (d *etcdDiscovery) Close() {
	_ = d.client.Close()
}
