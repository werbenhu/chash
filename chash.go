// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"encoding/json"
	"sync"
)

// CHash a warpper of Consistent hashing
type CHash struct {
	sync.RWMutex
	groups map[string]*Group
}

func New() *CHash {
	return &CHash{
		groups: make(map[string]*Group),
	}
}

// GetGroup retrieves a group by name
func (c *CHash) GetGroup(groupName string) (*Group, error) {
	c.RLock()
	defer c.RUnlock()
	group, ok := c.groups[groupName]
	if !ok {
		return nil, ErrGroupNotFound
	}
	return group, nil
}

// CreateGroup creates a new group with the given name and the number of replicas
func (c *CHash) CreateGroup(groupName string, replicas int) (*Group, error) {
	c.Lock()
	defer c.Unlock()
	if existing, ok := c.groups[groupName]; ok {
		return existing, ErrGroupExisted
	}

	group := NewGroup(groupName, replicas)
	c.groups[groupName] = group
	return group, nil
}

// RemoveGroup removes a group by name
func (c *CHash) RemoveGroup(groupName string) {
	c.Lock()
	defer c.Unlock()
	delete(c.groups, groupName)
}

// RemoveAllGroup removes all groups
func (c *CHash) RemoveAllGroup() {
	c.Lock()
	defer c.Unlock()

	for k := range c.groups {
		delete(c.groups, k)
	}
}

// Insert inserts a new key-value pair into a group
func (c *CHash) Insert(groupName string, key string, payload []byte) error {
	c.Lock()
	group, ok := c.groups[groupName]
	c.Unlock()
	if !ok {
		return ErrGroupNotFound
	}
	return group.Insert(key, payload)
}

// Delete removes a key from a group
func (c *CHash) Delete(groupName string, key string) error {
	c.Lock()
	group, ok := c.groups[groupName]
	c.Unlock()
	if !ok {
		return ErrGroupNotFound
	}
	group.Delete(key)
	return nil
}

// Match returns the key-value pair closest to the given key in a group
func (c *CHash) Match(groupName string, key string) (string, []byte, error) {
	c.RLock()
	group, ok := c.groups[groupName]
	c.RUnlock()
	if !ok {
		return "", nil, ErrGroupNotFound
	}
	return group.Match(key)
}

// Serialize serializes the CHash structure to JSON
func (c *CHash) Serialize() ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	return json.Marshal(c.groups)
}

// Restore deserializes a JSON representation of the CHash structure
func (c *CHash) Restore(data []byte) error {
	c.Lock()
	defer c.Unlock()
	if err := json.Unmarshal(data, &c.groups); err != nil {
		return err
	}
	for _, group := range c.groups {
		group.Init()
		for _, node := range group.Elements {
			group.hashElement(node)
		}
	}
	return nil
}
