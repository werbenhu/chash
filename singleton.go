package chash

import (
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
	return singleton.CreateBucket(bucketName, replicas)
}

func RemoveBucket(bucketName string) {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {
		return
	}
	singleton.RemoveBucket(bucketName)
}

func GetBucket(bucketName string) (*Bucket, error) {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {
		singleton = New()
	}
	return singleton.GetBucket(bucketName)
}

func Serialize() ([]byte, error) {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {
		singleton = New()
	}
	return singleton.Serialize()
}

func Restore(data []byte) error {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {
		singleton = New()
	}
	return singleton.Restore(data)
}
