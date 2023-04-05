# chash
Consistent hashing written by Go

### Getting chash

With Go module support, simply add the following import

`import "github.com/werbenhu/chash"`


### Simple Usage

#### Create a group
```
// create db group 
dbGroup, _ := chash.CreateGroup("db", 10000)

// insert elements
dbGroup.Insert("192.168.1.100:3306", []byte("mysql0-info"))
dbGroup.Insert("192.168.1.101:3306", []byte("mysql1-info"))

// create second group
webGroup, _ := chash.CreateGroup("web", 10000)

// insert elements
webGroup.Insert("192.168.2.100:80", []byte("web0-info"))
webGroup.Insert("192.168.2.101:80", []byte("web1-info"))
```

```
// use existed group
dbGroup, err := chash.GetGroup("db")
if err != nil {
    log.Printf("ERROR get group failed, err:%s\n", err.Error())
}

dbGroup.Insert("192.168.1.103:3306", []byte("mysql3-info"))
host, info, err := dbGroup.Match("user-id")
```

#### Match the element for a key
```
// match an element close to where key hashes to in the circle.
host, info, err := dbGroup.Match("user-id")
```

#### Delete element from a group
```
// delete element
dbGroup.Delete("192.168.1.102:3306")
```

#### Single Group
```
// you need to manager groups if there is more than one group.
group := chash.NewGroup("db", 10000)

group.Insert("192.168.1.100:3306", []byte("mysql0-info"))
group.Insert("192.168.1.101:3306", []byte("mysql1-info"))
host, info, err := group.Match("user-id")
```

### Examples
See the [example](example/main.go) .

### Contributions
Contributions and feedback are both welcomed and encouraged! Open an [issue](https://github.com/werbenhu/chash/issues) to report a bug, ask a question, or make a feature request.