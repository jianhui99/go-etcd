package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
)

// ---- server site ----

func getKey(serviceName string) string {
	return serviceName
}

func CustomServiceRegister(serviceName string, addr string) error {
	cli, err := GetEtcdClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	key := getKey(serviceName)

	// 创建租约
	leaseRes, err := cli.Grant(ctx, 5)
	if err != nil {
		return err
	}

	_, err = cli.Put(ctx, key, addr, clientv3.WithLease(leaseRes.ID))
	if err != nil {
		return err
	}

	keepAliveChan, err := cli.KeepAlive(context.TODO(), leaseRes.ID)
	if err != nil {
		return err
	}

	go func() {
		for item := range keepAliveChan {
			fmt.Printf("keep alive, leaseId:%x, ttl:%v \n ", item.ID, item.TTL)
		}
	}()

	return nil
}

// ---- server site ----

// ---- client site ----

type ServiceCache struct {
	data map[string]string
	sync.RWMutex
}

var cache *ServiceCache

func init() {
	cache = &ServiceCache{
		data: make(map[string]string),
	}
}

func CustomServiceDiscovery(serviceName string) string {
	cache.RLock()
	defer cache.RUnlock()

	return cache.data[serviceName]
}

func CustomLoadService(serviceName string) {
	cli, err := GetEtcdClient()
	if err != nil {
		log.Fatal(err)
	}
	key := getKey(serviceName)
	ctx := context.Background()
	getResp, err := cli.Get(ctx, key)
	if err != nil {
		log.Fatal("etcd get err:", err)
	}

	if getResp.Count > 0 {
		log.Printf("etcd get service:%s, count:%d", serviceName, getResp.Count)
		for _, item := range getResp.Kvs {
			cache.Lock()
			cache.data[string(item.Key)] = string(item.Value)
			cache.Unlock()
		}
	}

}

func CustomWatchService(serviceName string) {
	cli, err := GetEtcdClient()
	if err != nil {
		log.Fatal(err)
	}
	key := getKey(serviceName)
	ctx := context.Background()
	rch := cli.Watch(ctx, key)
	for response := range rch {
		for _, ev := range response.Events {
			if ev.Type == clientv3.EventTypePut {
				cache.Lock()
				cache.data[string(ev.Kv.Key)] = string(ev.Kv.Value)
				cache.Unlock()
				continue
			}
			if ev.Type == clientv3.EventTypeDelete {
				cache.Lock()
				delete(cache.data, string(ev.Kv.Key))
				cache.Unlock()
				continue
			}
		}
	}
}

// ---- client site ----
