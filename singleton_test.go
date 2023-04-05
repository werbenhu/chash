// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCHashSingletonCreateGroup(t *testing.T) {
	singleton = nil
	group, err := CreateGroup("werbenhu1", 2000)
	assert.NotNil(t, singleton)
	assert.Nil(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "werbenhu1", group.Name)

	existing, err := CreateGroup("werbenhu1", 3000)
	assert.Equal(t, ErrGroupExisted, err)
	assert.NotNil(t, existing)
	assert.Equal(t, group, existing)
}

func TestCHashSingletonRemoveGroup(t *testing.T) {
	singleton = nil

	group1, err := CreateGroup("werbenhu1", 2000)
	assert.Nil(t, err)
	assert.NotNil(t, group1)

	group2, err := CreateGroup("werbenhu2", 1000)
	assert.Nil(t, err)
	assert.NotNil(t, group2)

	assert.Equal(t, 2, len(singleton.groups))
	RemoveGroup("werbenhu1")
	assert.Equal(t, 1, len(singleton.groups))

	_, err = GetGroup("werbenhu1")
	assert.Equal(t, ErrGroupNotFound, err)
}

func TestCHashSingletonGetGroup(t *testing.T) {
	singleton = nil
	group1, err := CreateGroup("werbenhu2", 2000)
	assert.Nil(t, err)
	assert.NotNil(t, group1)

	group1, err = GetGroup("werbenhu2")
	assert.Nil(t, err)
	assert.NotNil(t, group1)

	group2, err := GetGroup("werbenhu3")
	assert.Nil(t, group2)
	assert.Equal(t, ErrGroupNotFound, err)
}

func TestCHashSingletonSerialize(t *testing.T) {
	singleton = nil
	CreateGroup("werbenhu1", 2000)
	CreateGroup("werbenhu2", 1000)

	singleton.Insert("werbenhu1", "192.168.1.101:8080", []byte("werbenhu101"))
	singleton.Insert("werbenhu1", "192.168.1.102:8080", []byte("werbenhu102"))

	singleton.Insert("werbenhu2", "192.168.2.101:8080", []byte("werbenhu201"))
	singleton.Insert("werbenhu2", "192.168.2.102:8080", []byte("werbenhu202"))

	bs, err := Serialize()
	expert := `{"werbenhu1":{"name":"werbenhu1","numberOfReplicas":2000,"elements":{"192.168.1.101:8080":{"key":"192.168.1.101:8080","payload":"d2VyYmVuaHUxMDE="},"192.168.1.102:8080":{"key":"192.168.1.102:8080","payload":"d2VyYmVuaHUxMDI="}}},"werbenhu2":{"name":"werbenhu2","numberOfReplicas":1000,"elements":{"192.168.2.101:8080":{"key":"192.168.2.101:8080","payload":"d2VyYmVuaHUyMDE="},"192.168.2.102:8080":{"key":"192.168.2.102:8080","payload":"d2VyYmVuaHUyMDI="}}}}`

	assert.Nil(t, err)
	assert.Equal(t, expert, string(bs))
}

func TestCHashSingletonRestore(t *testing.T) {
	singleton = nil
	data := []byte(`{"werbenhu1":{"name":"werbenhu1","numberOfReplicas":2000,"elements":{"192.168.1.101:8080":{"key":"192.168.1.101:8080","payload":"d2VyYmVuaHUxMDE="},"192.168.1.102:8080":{"key":"192.168.1.102:8080","payload":"d2VyYmVuaHUxMDI="}}},"werbenhu2":{"name":"werbenhu2","numberOfReplicas":1000,"elements":{"192.168.2.101:8080":{"key":"192.168.2.101:8080","payload":"d2VyYmVuaHUyMDE="},"192.168.2.102:8080":{"key":"192.168.2.102:8080","payload":"d2VyYmVuaHUyMDI="}}}}`)
	err := Restore(data)
	assert.Nil(t, err)

	group1, err := singleton.GetGroup("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, group1)
	assert.Equal(t, 4000, len(group1.rows))
	assert.Equal(t, group1.Elements, map[string]*Element{
		"192.168.1.101:8080": {
			Key:     "192.168.1.101:8080",
			Payload: []byte("werbenhu101"),
		},
		"192.168.1.102:8080": {
			Key:     "192.168.1.102:8080",
			Payload: []byte("werbenhu102"),
		},
	})

	group2, err := singleton.GetGroup("werbenhu2")
	assert.Nil(t, err)
	assert.NotNil(t, group2)
	assert.Equal(t, 2000, len(group2.rows))
	assert.Equal(t, group2.Elements, map[string]*Element{
		"192.168.2.101:8080": {
			Key:     "192.168.2.101:8080",
			Payload: []byte("werbenhu201"),
		},
		"192.168.2.102:8080": {
			Key:     "192.168.2.102:8080",
			Payload: []byte("werbenhu202"),
		},
	})

	wrongData := append(data, []byte("--werbenhu--")...)
	err = Restore(wrongData)
	assert.NotNil(t, err)
}
