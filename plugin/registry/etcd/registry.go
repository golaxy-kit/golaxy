package etcd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	hash "github.com/mitchellh/hashstructure"
	"github.com/pangdogs/galaxy/plugin/registry"
	"github.com/pangdogs/galaxy/service"
	"go.etcd.io/etcd/client/v3"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	prefix = "/galaxy/registry/"
)

func newRegistry(options ...Option) registry.Registry {
	opts := Options{}
	Default()(&opts)

	for i := range options {
		options[i](&opts)
	}

	return &etcdRegistry{
		options:  opts,
		register: make(map[string]uint64),
		leases:   make(map[string]clientv3.LeaseID),
	}
}

type etcdRegistry struct {
	options  Options
	client   *clientv3.Client
	register map[string]uint64
	leases   map[string]clientv3.LeaseID
	sync.RWMutex
}

func (e *etcdRegistry) Init(ctx service.Context) {
	client, err := clientv3.New(e.configure())
	if err != nil {
		panic(err)
	}
	e.client = client
}

func (e *etcdRegistry) Shut() {
	if e.client != nil {
		e.client.Close()
	}
}

func (e *etcdRegistry) Register(ctx context.Context, service registry.Service, ttl time.Duration) error {
	if len(service.Nodes) <= 0 {
		return errors.New("require at least one node")
	}

	var anyErr error

	for _, node := range service.Nodes {
		if err := e.registerNode(ctx, service, node, ttl); err != nil {
			anyErr = err
		}
	}

	return anyErr
}

func (e *etcdRegistry) Deregister(ctx context.Context, service registry.Service) error {
	return nil
}

func (e *etcdRegistry) GetService(ctx context.Context, serviceName string) ([]registry.Service, error) {
	return nil, nil
}

func (e *etcdRegistry) ListServices(ctx context.Context) ([]registry.Service, error) {
	return nil, nil
}

func (e *etcdRegistry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return nil, nil
}

func (e *etcdRegistry) configure() clientv3.Config {
	if e.options.EtcdConfig != nil {
		return *e.options.EtcdConfig
	}

	config := clientv3.Config{
		Endpoints:   e.options.Endpoints,
		DialTimeout: e.options.Timeout,
		Username:    e.options.Username,
		Password:    e.options.Password,
		LogConfig:   e.options.ZapConfig,
	}

	if e.options.Secure || e.options.TLSConfig != nil {
		tlsConfig := e.options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		config.TLS = tlsConfig
	}

	return config
}

func (e *etcdRegistry) registerNode(ctx context.Context, s registry.Service, node registry.Node, ttl time.Duration) error {
	if len(s.Nodes) == 0 {
		return errors.New("require at least one node")
	}

	// check existing lease cache
	e.RLock()
	leaseID, ok := e.leases[s.Name+node.Id]
	e.RUnlock()

	if !ok {
		// missing lease, check if the key exists
		ctx, cancel := context.WithTimeout(ctx, e.options.Timeout)
		defer cancel()

		// look for the existing key
		rsp, err := e.client.Get(ctx, nodePath(s.Name, node.Id), clientv3.WithSerializable())
		if err != nil {
			return err
		}

		// get the existing lease
		for _, kv := range rsp.Kvs {
			if kv.Lease > 0 {
				leaseID = clientv3.LeaseID(kv.Lease)

				// decode the existing node
				srv := decode(kv.Value)
				if srv == nil || len(srv.Nodes) == 0 {
					continue
				}

				// create hash of service; uint64
				h, err := hash.Hash(srv.Nodes[0], nil)
				if err != nil {
					continue
				}

				// save the info
				e.Lock()
				e.leases[s.Name+node.Id] = leaseID
				e.register[s.Name+node.Id] = h
				e.Unlock()

				break
			}
		}
	}

	var leaseNotFound bool

	//// renew the lease if it exists
	//if leaseID > 0 {
	//	if logger.V(logger.TraceLevel, logger.DefaultLogger) {
	//		logger.Tracef("Renewing existing lease for %s %d", s.Name, leaseID)
	//	}
	//	if _, err := e.client.KeepAliveOnce(context.TODO(), leaseID); err != nil {
	//		if err != rpctypes.ErrLeaseNotFound {
	//			return err
	//		}
	//
	//		if logger.V(logger.TraceLevel, logger.DefaultLogger) {
	//			logger.Tracef("Lease not found for %s %d", s.Name, leaseID)
	//		}
	//		// lease not found do register
	//		leaseNotFound = true
	//	}
	//}

	// create hash of service; uint64
	h, err := hash.Hash(node, nil)
	if err != nil {
		return err
	}

	// get existing hash for the service node
	e.Lock()
	v, ok := e.register[s.Name+node.Id]
	e.Unlock()

	// the service is unchanged, skip registering
	if ok && v == h && !leaseNotFound {
		//if logger.V(logger.TraceLevel, logger.DefaultLogger) {
		//	logger.Tracef("Service %s node %s unchanged skipping registration", s.Name, node.Id)
		//}
		return nil
	}

	//service := &registry.Service{
	//	Name:      s.Name,
	//	Version:   s.Version,
	//	Metadata:  s.Metadata,
	//	Endpoints: s.Endpoints,
	//	Nodes:     []registry.Node{node},
	//}

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	var lgr *clientv3.LeaseGrantResponse
	if ttl.Seconds() > 0 {
		// get a lease used to expire keys since we have a ttl
		lgr, err = e.client.Grant(ctx, int64(ttl.Seconds()))
		if err != nil {
			return err
		}
	}

	//if logger.V(logger.TraceLevel, logger.DefaultLogger) {
	//	logger.Tracef("Registering %s id %s with lease %v and leaseID %v and ttl %v", service.Name, node.Id, lgr, lgr.ID, options.TTL)
	//}
	// create an entry for the node
	//if lgr != nil {
	//	_, err = e.client.Put(ctx, nodePath(service.Name, node.Id), encode(service), clientv3.WithLease(lgr.ID))
	//} else {
	//	_, err = e.client.Put(ctx, nodePath(service.Name, node.Id), encode(service))
	//}
	//if err != nil {
	//	return err
	//}

	e.Lock()
	// save our hash of the service
	e.register[s.Name+node.Id] = h
	// save our leaseID of the service
	if lgr != nil {
		e.leases[s.Name+node.Id] = lgr.ID
	}
	e.Unlock()

	return nil
}

func encode(s registry.Service) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func decode(ds []byte) *registry.Service {
	var s *registry.Service
	json.Unmarshal(ds, &s)
	return s
}

func nodePath(s, id string) string {
	service := strings.ReplaceAll(s, "/", "-")
	node := strings.ReplaceAll(id, "/", "-")
	return path.Join(prefix, service, node)
}

func servicePath(s string) string {
	return path.Join(prefix, strings.ReplaceAll(s, "/", "-"))
}