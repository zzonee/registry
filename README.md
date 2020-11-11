##服务注册与服务发现

接口：
```go
type Registry interface {
	Register(ops ...Option)
	Close()
}

type Discovery interface {
	Discover(target string) (<-chan []Instance, error)
	Close()
}
```

服务配置数据结构：
```go
type Instance struct {
	// 服务名
	ServiceName string            `json:"servicename"`
	// 地址，"127.0.0.1:8000"
	Address     string            `json:"address"`
	// 元数据，可以带上自定义的附加数据
	Metadata    map[string]string `json:"metadata"`
}
```

### 使用etcd做grpc服务发现
```go
package main
import (
	registry "github.com/zzonee/registry/grpc/etcd"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
)

func main() {
	...
	conf := clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
	}
	err := registry.RegisterBuilder(conf)
	if err != nil {
		panic(err)
	}
	
	...
	conn, err := grpc.Dial("etcd://default/servicename")
	...
}

```

### 在k8s中做grpc服务发现
```go
package main
import (
	registry "github.com/zzonee/registry/grpc/k8s"
	"google.golang.org/grpc"
)

func main() {
	...
	
	err := registry.RegisterBuilder()
	if err != nil {
		panic(err)
	}
	
	...
	conn, err := grpc.Dial("k8s://default/servicename:portname")
	...
}

```