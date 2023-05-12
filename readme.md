NameService is a distributed name pool.

if a name pool has been initialized(using InitNameNodes): [svc_1, svc_2, svc_3].

Using GetName, you can retrieve a name from it.

Using KeepAlive, you can maintain a heartbeat.

When ctx ends, KeepAlive will exit, and the name will be returned to the pool.

Usage:
```go
ns, err := NewNameService(ctx,
	"127.0.0.123", // Configure the IP of this service.
	[]string{ // zk configuration
		"127.0.0.1",
	},
	"/name_allocator/test_tasks") // Select the root path for this cluster.
if err != nil {
	return err
}

name, err := ns.GetName() // Get a name.
if err != nil { // If the name pool is exhausted, an error will be returned.
	log.Errorf("get name err:%v", err)
	return err
}

log.Infof("get name:%v", name)
go ns.KeepAlive(name) // If you get a name, use keepAlive to hold onto it.
// dosth with name
```