// copy from go-plugins/registry/etcdv3

package etcdv3

import (
	"context"
	"time"

	"github.com/gowins/dionysus/grpc/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Watcher struct {
	revision int64
	stop     chan struct{}
	w        clientv3.WatchChan
	client   *clientv3.Client
	timeout  time.Duration
}

func newWatcher(r *Registry, timeout time.Duration, opts ...registry.WatchOption) (registry.Watcher, error) {
	var wo registry.WatchOptions
	for _, o := range opts {
		o(&wo)
	}

	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan struct{})

	go func() {
		<-stop
		cancel()
	}()

	watchPath := prefix
	if len(wo.Service) > 0 {
		watchPath = servicePath(wo.Service) + "/"
	}

	resp, err := r.client.Get(ctx, watchPath, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	return &Watcher{
		revision: resp.Header.Revision,
		stop:     stop,
		w:        r.client.Watch(ctx, watchPath, clientv3.WithPrefix(), clientv3.WithPrevKV(), clientv3.WithRev(resp.Header.Revision+1)),
		client:   r.client,
		timeout:  timeout,
	}, nil
}

func (w *Watcher) Next() (*registry.Result, error) {
	for resp := range w.w {
		if resp.Err() != nil {
			return nil, resp.Err()
		}

		if resp.CompactRevision > w.revision {
			w.revision = resp.CompactRevision
		}
		if resp.Header.GetRevision() > w.revision {
			w.revision = resp.Header.GetRevision()
		}

		for _, ev := range resp.Events {
			service := decode(ev.Kv.Value)
			var action string

			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					action = "create"
				} else if ev.IsModify() {
					action = "update"
				}
			case clientv3.EventTypeDelete:
				action = "delete"

				// get service from prevKv
				service = decode(ev.PrevKv.Value)
			}

			if service == nil {
				continue
			}

			return &registry.Result{
				Action:  action,
				Service: service,
			}, nil
		}
	}

	return nil, registry.ErrWatcherStopped
}

func (w *Watcher) Stop() {
	select {
	case <-w.stop:
		return
	default:
		close(w.stop)
	}
}
