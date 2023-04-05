# chash
Consistent hashing written by Go

### Getting chash

With Go module support, simply add the following import

`import "github.com/werbenhu/chash"`


### Simple Usage
```
group, _ := chash.CreateGroup("db-consistent-hash", 10000)

group.Insert("192.168.1.100:3306", []byte("mysql0-info"))
group.Insert("192.168.1.101:3306", []byte("mysql1-info"))

host, info, err := group.Match("user-id")
```

#### Single Group

group := chash.NewGroup("db-consistent-hash", 10000)

group.Insert("192.168.1.100:3306", []byte("mysql0-info"))
group.Insert("192.168.1.101:3306", []byte("mysql1-info"))

host, info, err := group.Match("user-id")

### Examples
See the [example](example/main.go) .

### Contributions
Contributions and feedback are both welcomed and encouraged! Open an [issue](https://github.com/werbenhu/chash/issues) to report a bug, ask a question, or make a feature request.