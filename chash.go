// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"encoding/json"
	"sync"
)

type CHash struct {
	sync.RWMutex
	groups map[string]*Group
}

func New() *CHash {
	return &CHash{
		groups: make(map[string]*Group),
	}
}

func (c *CHash) GetGroup(groupName string) (*Group, error) {
	c.RLock()
	defer c.RUnlock()
	group, ok := c.groups[groupName]
	if !ok {
		return nil, ErrGroupNotFound
	}
	return group, nil
}

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

func (c *CHash) RemoveGroup(groupName string) {
	c.Lock()
	defer c.Unlock()
	delete(c.groups, groupName)
}

func (c *CHash) Insert(groupName string, key string, payload []byte) error {
	c.Lock()
	group, ok := c.groups[groupName]
	c.Unlock()
	if !ok {
		return ErrGroupNotFound
	}
	return group.Insert(key, payload)
}

func (c *CHash) Delete(groupName string, key string) error {
	c.Lock()
	group, ok := c.groups[groupName]
	c.Unlock()
	if !ok {
		return ErrGroupNotFound
	}
	return group.Delete(key)
}

func (c *CHash) Match(groupName string, key string) (string, []byte, error) {
	c.RLock()
	group, ok := c.groups[groupName]
	c.RUnlock()
	if !ok {
		return "", nil, ErrGroupNotFound
	}
	return group.Match(key)
}

func (c *CHash) Serialize() ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	return json.Marshal(c.groups)
}

func (c *CHash) Restore(data []byte) error {
	c.Lock()
	defer c.Unlock()
	if err := json.Unmarshal(data, &c.groups); err != nil {
		return err
	}
	for _, group := range c.groups {
		group.Init()
		for _, node := range group.Agents {
			group.hashAgent(node)
		}
	}
	return nil
}
