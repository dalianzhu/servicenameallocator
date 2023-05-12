NameService is a distributed name pool that depends on ZooKeeper.

If a name pool has been initialized(using `InitNameNodes`): [svc_1, svc_2, svc_3].

Using `GetName` can obtain a name from it, for example "svc_2".

Using `KeepAlive` can maintain a heartbeat.

When ctx ends, `KeepAlive` will exit, and the name will be returned to the pool.

You can also use `InitNameNodes` to add new names to the pool.

Usage:
```go
ns, err := NewNameService(ctx,
	"127.0.0.123", // sets the address for this service, will serve as its identifier.
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

// You can use this name, for example, to load configurations.
loadConfigs(name)
```