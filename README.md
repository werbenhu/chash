# chash
Consistent Hashing written by Go

### Getting chash

With Go module support, simply add the following import

`import "github.com/werbenhu/chash"`


### Simple Usage
```
bucket, _ := chash.CreateBucket("db-consistent-hash", 10000)

bucket.Insert("192.168.1.100:3306", []byte("mysql0-info"))
bucket.Insert("192.168.1.101:3306", []byte("mysql1-info"))

host, info, err := bucket.Match("user-id")
```

### Examples
See the [example](example/main.go) .

### Contributions
Contributions and feedback are both welcomed and encouraged! Open an [issue](https://github.com/werbenhu/chash/issues) to report a bug, ask a question, or make a feature request.