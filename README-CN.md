
<div align='center'>
<a href="https://github.com/werbenhu/chash/actions"><img src="https://github.com/werbenhu/chash/workflows/Go/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/werbenhu/chash"><img src="https://goreportcard.com/badge/github.com/werbenhu/chash"></a>
<a href="https://coveralls.io/github/werbenhu/chash?branch=main"><img src="https://coveralls.io/repos/github/werbenhu/chash/badge.svg?branch=main"></a>   
<a href="https://github.com/werbenhu/chash"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
<a href="https://pkg.go.dev/github.com/werbenhu/chash"><img src="https://pkg.go.dev/badge/github.com/werbenhu/chash.svg"></a>
</div>

[English](README.md) | [简体中文](README-CN.md)

# chash
**用Go写的一致性哈希算法**

## 什么是一致性哈希

一致性哈希`Consistent Hash`是为了解决由于分布式系统中节点的增加或减少而带来的大量失效问题，它可以有效地降低这种失效影响，从而提高分布式系统的性能和可用性。

### 普通哈希的问题

普通哈希函数是 key % n，其中 n 是服务器数量。它有两个主要缺点：
1. 不能水平扩展，或者换句话说，不具备分区容错性。当添加新服务器时，所有现有的映射都会被破坏。这可能会引入痛苦的维护工作和系统停机时间。
2. 可能不能实现负载均衡。如果数据不是均匀分布的，这可能会导致一些服务器过热饱和，而其他服务器则处于空闲状态并几乎为空。

问题2可以通过先对键进行哈希，然后哈希(key) % n，以便哈希键更有可能被均匀分布来解决。但是，这不能解决问题1。我们需要找到一个可以分配key并且不依赖于n的解决方案。

### 一致性哈希的简单认识

关于环，添加删除节点，还有虚拟节点等，参考：[一致性哈希的简单认识](https://baijiahao.baidu.com/s?id=1735480432495470467&wfr=spider&for=pc)

## 入门指南

使用Go模块支持，只需添加以下导入即可

`import "github.com/werbenhu/chash"`


## 简单用法

### 创建一个组
```
// 创建db组，通过chash.CreateGroup()创建组，并指定每个节点的虚拟节点数量为10000
// 内部会有一个全局的单例chash对象管理组
dbGroup, _ := chash.CreateGroup("db", 10000)

// 插入元素(多个myqsl服务器)
dbGroup.Insert("192.168.1.100:3306", []byte("mysql0-info"))
dbGroup.Insert("192.168.1.101:3306", []byte("mysql1-info"))

// 创建第二个组
webGroup, _ := chash.CreateGroup("web", 10000)

// 插入元素(多个http服务器)
webGroup.Insert("192.168.2.100:80", []byte("web0-info"))
webGroup.Insert("192.168.2.101:80", []byte("web1-info"))
```

```
// 使用已经存在的组
dbGroup, err := chash.GetGroup("db")
if err != nil {
    log.Printf("ERROR get group failed, err:%s\n", err.Error())
}

dbGroup.Insert("192.168.1.103:3306", []byte("mysql3-info"))
host, info, err := dbGroup.Match("user-id")
```

### 给某个用户ID，匹配一个mysql服务器
```
// 匹配哈希到环上的键的附近元素。
host, info, err := dbGroup.Match("user-id")
```

### 从组中删除元素
```
// 删除元素
dbGroup.Delete("192.168.1.102:3306")
```

### 获取组的所有元素
```
elements := dbGroup.GetElemens()
```

### 也可以使用一个独立组
```
// chash.NewGrou()创建的组，不被全局的chash对象所管理
// 如果有多个组，您需要自己管理组。
group := chash.NewGroup("db", 10000)

group.Insert("192.168.1.100:3306", []byte("mysql0-info"))
group.Insert("192.168.1.101:3306", []byte("mysql1-info"))
host, info, err := group.Match("user-id")
```

## 示例
请参见 [example](example/main.go) .

## 贡献
欢迎贡献和反馈！提出问题或提出功能请求请提 [issue](https://github.com/werbenhu/chash/issues) .
