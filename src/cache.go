package main

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type Cache interface {
	Get(key string) (value string, err error)
	Set(key string, value string) (err error)
}

type memcached struct {
	cache *memcache.Client
	ttl   int32
}

func (m *memcached) Get(key string) (value string, err error) {
	item, err := m.cache.Get(key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

func (m *memcached) Set(key string, value string) (err error) {
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: m.ttl,
	}
	err = m.cache.Set(item)
	return err
}

func NewMemcached(host string, port string, ttl int32) Cache {
	client := memcache.New(host + ":" + port)
	return &memcached{
		cache: client,
		ttl:   ttl,
	}
}
