package registry

import "time"

// Registry Options
type Options struct {
	ServiceName string            `json:"servicename"`
	Address     string            `json:"address"`
	Metadata    map[string]string `json:"metadata"`
	// 服务有效时长
	RegisterTTL      time.Duration `json:"-"` // time to live, 服务失活一段时间后自动从注册中心删除
	RegisterInterval time.Duration `json:"-"` // 注册间隔时长，也可不要 默认为RegisterTTL/3
}

type Option func(*Options)

type Registry interface {
	Register(ops ...Option)
	Close()
}

func ServiceName(name string) Option {
	return func(o *Options) {
		o.ServiceName = name
	}
}

func Address(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

func Metadata(m map[string]string) Option {
	return func(o *Options) {
		o.Metadata = m
	}
}

func RegisterTTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = ttl
	}
}

func RegisterInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = interval
	}
}

type Instance struct {
	ServiceName string            `json:"servicename"`
	Address     string            `json:"address"`
	Metadata    map[string]string `json:"metadata"`
}

// 服务发现接口
// target的具体格式由其实现决定
type Discovery interface {
	Discover(target string) (<-chan []Instance, error)
	Close()
}
