package servicenameallocator

import (
	"context"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

/*
	NameService is a distributed name pool.

if a name pool has been initialized: [svc_1, svc_2, svc_3].
Using GetName, you can retrieve a name from it.
Using KeepAlive, you can maintain a heartbeat.
When ctx ends, KeepAlive will exit, and the name will be returned to the pool.

Usage:
```
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
*/
type NameService struct {
	ctx    context.Context
	addr   string
	zkConn *zk.Conn
	root   string
}

// NewNameService creates and returns a new instance of NameService, which represents a distributed name pool.
// It takes a context.Context, a string representing the IP address of this service,
// a slice of strings representing ZooKeeper servers, and a string representing the root path of this cluster.
// It returns a pointer to NameService and an error, if any. The returned NameService has a connection to ZooKeeper.
func NewNameService(ctx context.Context, thisSvcAddr string,
	zkServers []string, root string) (*NameService, error) {

	conn, _, err := zk.Connect(zkServers, time.Second*5)
	if err != nil {
		return nil, err
	}
	ns := &NameService{ctx: ctx,
		addr:   thisSvcAddr,
		zkConn: conn, root: root}
	return ns, nil
}
