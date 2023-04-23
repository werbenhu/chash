package chash

import (
	"sync"
)

var (
	mu sync.Mutex

	// singleton is a pointer to a CHash instance, which will be created when necessary.
	singleton *CHash
)

// CreateGroup creates a new group in the CHash instance and returns a pointer to
// the Group object. If the group already exists, it returns an error.
func CreateGroup(groupName string, replicas int) (*Group, error) {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {

		// If singleton is nil, we create a new instance of CHash using the New()
		// function and assign it to the singleton variable.
		singleton = New()
	}
	return singleton.CreateGroup(groupName, replicas)
}

// RemoveGroup removes the specified group from the CHash instance.
func RemoveGroup(groupName string) {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {
		return
	}
	singleton.RemoveGroup(groupName)
}

// RemoveAllGroup removes all groups from the CHash instance.
func RemoveAllGroup() {
	mu.Lock()
	defer mu.Unlock()

	if singleton == nil {
		return
	}
	singleton.RemoveAllGroup()
}

// GetGroup retrieves the specified group from the CHash instance.
func GetGroup(groupName string) (*Group, error) {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {
		singleton = New()
	}
	return singleton.GetGroup(groupName)
}

// Serialize serializes the entire CHash object to a byte slice.
func Serialize() ([]byte, error) {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {
		singleton = New()
	}
	return singleton.Serialize()
}

// Restore restores the CHash object from serialized data.
func Restore(data []byte) error {
	mu.Lock()
	defer mu.Unlock()
	if singleton == nil {
		singleton = New()
	}
	return singleton.Restore(data)
}
