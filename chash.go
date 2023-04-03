// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"encoding/json"
	"sync"
)

type Handler interface {
	OnAgentInsert(name string, agent *Agent) error
	OnAgentDelete(name string, agent *Agent) error
	OnBucketCreate(bucket *Bucket) error
	OnBucketRemove(bucket *Bucket) error
}

var (
	mu        sync.Mutex
	singleton *CHash
)

func CreateBucket(bucketName string, replicas int) (*Bucket, error) {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {
		singleton = New()
	}
	singleton.CreateBucket(bucketName, replicas)
	return singleton.GetBucket(bucketName)
}

type CHash struct {
	sync.RWMutex
	handler Handler
	buckets map[string]*Bucket
}

func New() *CHash {
	return &CHash{
		buckets: make(map[string]*Bucket),
	}
}

func (c *CHash) SetHandler(handler Handler) {
	c.handler = handler
}

func (c *CHash) GetBucket(bucketName string) (*Bucket, error) {
	c.RLock()
	defer c.RUnlock()
	bucket, ok := c.buckets[bucketName]
	if !ok {
		return nil, ErrBucketNotFound
	}
	return bucket, nil
}

func (c *CHash) CreateBucket(bucketName string, replicas int) error {
	c.RLock()
	if _, ok := c.buckets[bucketName]; ok {
		c.RUnlock()
		return ErrBucketExisted
	}
	c.RUnlock()

	bucket := NewBucket(bucketName, replicas, c.handler)
	if c.handler != nil {
		if err := c.handler.OnBucketCreate(bucket); err != nil {
			return err
		}
	}

	c.Lock()
	defer c.Unlock()
	c.buckets[bucketName] = bucket
	return nil
}

func (c *CHash) RemoveBucket(bucketName string) error {
	c.RLock()
	existing, ok := c.buckets[bucketName]
	if !ok {
		c.RUnlock()
		return ErrBucketExisted
	}
	c.RUnlock()

	if c.handler != nil {
		if err := c.handler.OnBucketRemove(existing); err != nil {
			return err
		}
	}
	c.Lock()
	defer c.Unlock()
	delete(c.buckets, bucketName)
	return nil
}

func (c *CHash) InsertAgent(bucketName string, key string, payload []byte) error {
	c.Lock()
	bucket, ok := c.buckets[bucketName]
	c.Unlock()
	if !ok {
		return ErrBucketNotFound
	}
	return bucket.Insert(key, payload)
}

func (c *CHash) DeleteAgent(bucketName string, key string) error {
	c.Lock()
	bucket, ok := c.buckets[bucketName]
	c.Unlock()
	if !ok {
		return ErrBucketNotFound
	}
	return bucket.Delete(key)
}

func (c *CHash) Match(bucketName string, key string) (string, []byte, error) {
	c.RLock()
	bucket, ok := c.buckets[bucketName]
	c.RUnlock()
	if !ok {
		return "", nil, ErrBucketNotFound
	}
	return bucket.Match(key)
}

func (c *CHash) Serialize() ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	return json.Marshal(c.buckets)
}

func (c *CHash) Restore(data []byte) error {
	c.Lock()
	defer c.Unlock()
	if err := json.Unmarshal(data, &c.buckets); err != nil {
		return err
	}
	for _, bucket := range c.buckets {
		bucket.Init()
		for _, node := range bucket.Agents {
			bucket.Insert(node.Key, node.Payload)
		}
	}
	return nil
}
