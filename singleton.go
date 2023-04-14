package chash

import (
	"sync"
)

var (
	mu        sync.Mutex
	singleton *CHash
)

func CreateGroup(groupName string, replicas int) (*Group, error) {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {
		singleton = New()
	}
	return singleton.CreateGroup(groupName, replicas)
}

func RemoveGroup(groupName string) {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {
		return
	}
	singleton.RemoveGroup(groupName)
}

func RemoveAllGroup() {
	mu.Lock()
	defer mu.Unlock()

	if singleton == nil {
		return
	}
	singleton.RemoveAllGroup()
}

func GetGroup(groupName string) (*Group, error) {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {
		singleton = New()
	}
	return singleton.GetGroup(groupName)
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
