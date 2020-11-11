package etcd

import (
	"github.com/zzonee/registry"
	"go.etcd.io/etcd/clientv3"
	"os"
	"testing"
	"time"
)

var conf clientv3.Config
var stop chan bool

func TestMain(m *testing.M) {
	conf.Endpoints = []string{
		"127.0.0.1:2379",
	}
	stop = make(chan bool, 2)
	os.Exit(m.Run())
}

func Test_Register(t *testing.T) {
	reg1, err := NewRegistry(conf)
	if err != nil {
		t.Failed()
		return
	}
	reg2, err := NewRegistry(conf)
	if err != nil {
		t.Failed()
		return
	}
	go discovery(t, 1)
	reg1.Register(
		registry.ServiceName("test"),
		registry.Address("127.0.0.1:8000"),
		registry.Metadata(map[string]string{"version": "1"}),
		registry.RegisterTTL(time.Second*30),
		registry.RegisterInterval(time.Second*15)) // 15s 注册一次

	reg2.Register(
		registry.ServiceName("test"),
		registry.Address("127.0.0.1:8001"),
		registry.Metadata(map[string]string{"version": "1"}),
		registry.RegisterTTL(time.Second*30)) // 10s 注册一次

	go discovery(t, 2)
	<-time.After(time.Second * 25)
	reg1.Close()
	<-time.After(time.Second * 10)
	reg2.Close()
	<-stop
}

func discovery(t *testing.T, id int) {
	d, err := NewDiscovery(conf)
	if err != nil {
		stop <- true
		t.Failed()
		return
	}
	ch, _ := d.Discover("test")
	exit := time.After(time.Minute)
	i := 0
	for {
		select {
		case inss := <-ch:
			for _, ins := range inss {
				t.Log("d", id, "turn", i, "discov===>", ins)
			}
			if len(inss) == 0 {
				t.Log("d", id, "turn", i, "no instance")
			}
		case <-exit:
			t.Log("d", id, "discov===>exit")
			stop <- true
			return
		}
		i += 1
	}
}
