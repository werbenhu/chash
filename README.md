
<div align='center'>
<a href="https://github.com/werbenhu/chash/actions"><img src="https://github.com/werbenhu/chash/workflows/Go/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/werbenhu/chash"><img src="https://goreportcard.com/badge/github.com/werbenhu/chash"></a>
<a href="https://coveralls.io/github/werbenhu/chash?branch=main"><img src="https://coveralls.io/repos/github/werbenhu/chash/badge.svg?branch=main"></a>  
<a href="https://github.com/werbenhu/chash"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
<a href="https://pkg.go.dev/github.com/werbenhu/chash"><img src="https://pkg.go.dev/badge/github.com/werbenhu/chash.svg"></a>
</div>

[English](README.md) | [简体中文](README-CN.md)
# chash
**Consistent hashing written by Go**

## What is consistent hashing

> Consistent hashing is a hashing technique that performs really well when operated in a dynamic environment where the distributed system scales up and scales down frequently. 

Consistent Hashing allows distributing data in such a way that minimize reorganization when nodes are added or removed, hence making the system easier to scale up or down.

The key idea is that it's a distribution scheme that DOES NOT depend directly on the number of servers.

In Consistent Hashing, when the hash table is resized, in general only k / n keys need to be remapped, where k is the total number of keys and n is the total number of servers.

When a new node is added, it takes shares from a few hosts without touching other's shares
When a node is removed, its shares are shared by other hosts.

## Getting started

With Go module support, simply add the following import

`import "github.com/werbenhu/chash"`


## Simple Usage

### Create a group
```go
// Create a group named "db" with 10000 virtual elements for each key
// which internally manages the group using a global singleton chash object.
dbGroup, _ := chash.CreateGroup("db", 10000)

// Insert elements (multiple MySQL servers).
dbGroup.Insert("192.168.1.100:3306", []byte("mysql0-info"))
dbGroup.Insert("192.168.1.101:3306", []byte("mysql1-info"))

// Create a second group.
webGroup, _ := chash.CreateGroup("web", 10000)

// Insert elements (multiple HTTP servers).
webGroup.Insert("192.168.2.100:80", []byte("web0-info"))
webGroup.Insert("192.168.2.101:80", []byte("web1-info"))
```

```go
// Use an existing group.
dbGroup, err := chash.GetGroup("db")
if err != nil {
    log.Printf("ERROR get group failed, err:%s\n", err.Error())
}

dbGroup.Insert("192.168.1.103:3306", []byte("mysql3-info"))
host, info, err := dbGroup.Match("user-id")
```

### Match a MySQL server for a user ID
```go
// match an element close to where key hashes to in the circle.
host, info, err := dbGroup.Match("user-id")
```

### Delete element from a group
```go
// delete element
dbGroup.Delete("192.168.1.102:3306")
```

### Get all elements of a group
```go
elements := dbGroup.GetElemens()
```

### Using an independent group
```go
// A group created by chash.NewGrou() is not managed by the global chash object
// You need to manage groups yourself if there are more than one group.
group := chash.NewGroup("db", 10000)

group.Insert("192.168.1.100:3306", []byte("mysql0-info"))
group.Insert("192.168.1.101:3306", []byte("mysql1-info"))
host, info, err := group.Match("user-id")
```

## Examples
See the [example](example/main.go) .

## Contributions
Contributions and feedback are both welcomed and encouraged! Open an [issue](https://github.com/werbenhu/chash/issues) to report a bug, ask a question, or make a feature request.
