package k8s

import (
	"context"
	"fmt"
	"github.com/zzonee/registry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
	"sync"
)

type k8sDiscovery struct {
	clientset *kubernetes.Clientset
	namespace string

	closeOnce sync.Once
	closeCh   chan struct{}
}

func NewDiscovery(namespace string) (registry.Discovery, error) {
	if namespace == "" {
		namespace = "default"
	}
	conf, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, err
	}
	return &k8sDiscovery{
		clientset: clientset,
		namespace: namespace,
		closeCh:   make(chan struct{}),
	}, nil
}

// target: service-name:port or service-name:port-name
func (d *k8sDiscovery) Discover(target string) (<-chan []registry.Instance, error) {
	service, port := parse(target)
	if service == "" || port == "" {
		return nil, fmt.Errorf("target not valid: %s", target)
	}
	ch := make(chan []registry.Instance)
	return ch, d.watch(ch, service, port)
}

func parse(target string) (service, port string) {
	ss := strings.Split(target, ":")
	if len(ss) == 2 {
		service, port = ss[0], ss[1]
	}
	return
}

func (d *k8sDiscovery) watch(ch chan<- []registry.Instance, service, port string) error {
	watcher, err := d.clientset.CoreV1().Endpoints(d.namespace).
		Watch(context.Background(),
			metav1.ListOptions{FieldSelector: fmt.Sprintf("%s=%s", "metadata.name", service)})
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-d.closeCh:
				return
			case <-watcher.ResultChan():
			}

			endpoints, err := d.clientset.CoreV1().Endpoints(d.namespace).
				List(context.TODO(), metav1.ListOptions{FieldSelector: fmt.Sprintf("%s=%s", "metadata.name", service)})
			if err != nil {
				continue
			}

			var inss []registry.Instance
			for _, endpoint := range endpoints.Items {
				for _, subset := range endpoint.Subsets {
					realPort := port
					for _, p := range subset.Ports {
						if p.Name == port {
							realPort = fmt.Sprint(p.Port)
							break
						}
					}
					for _, addr := range subset.Addresses {
						ins := registry.Instance{Address: fmt.Sprintf("%s:%s", addr.IP, realPort)}
						inss = append(inss, ins)
					}
				}
			}
			ch <- inss
		}
	}()
	return nil
}

func (d *k8sDiscovery) Close() {
	d.closeOnce.Do(func() {
		close(d.closeCh)
	})
}
