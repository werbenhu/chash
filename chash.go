// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"encoding/json"
	"sync"
)

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
	if err := singleton.CreateBucket(bucketName, replicas); err != nil {
		return nil, err
	}
	return singleton.GetBucket(bucketName)
}

type CHash struct {
	sync.RWMutex
	buckets map[string]*Bucket
}

func New() *CHash {
	return &CHash{
		buckets: make(map[string]*Bucket),
	}
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
	c.Lock()
	defer c.Unlock()
	if _, ok := c.buckets[bucketName]; ok {
		return ErrBucketExisted
	}

	bucket := NewBucket(bucketName, replicas)
	c.buckets[bucketName] = bucket
	return nil
}

func (c *CHash) RemoveBucket(bucketName string) {
	c.Lock()
	defer c.Unlock()
	delete(c.buckets, bucketName)
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
			bucket.hashAgent(node)
		}
	}
	return nil
}
