package servicenameallocator

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	heartbeatInterval = 3 * time.Second
)

// NameInfo ...
type NameInfo struct {
	Name      string
	Addr      string
	ExpiredAt time.Time
}

// KeepAlive is a method of NameService that maintains the heartbeat of a given name.
// It takes a string representing the name to be kept alive.
// It continuously updates the name's expiration time in ZooKeeper and waits for the next heartbeat interval.
// If the service's ctx is canceled, it stops the heartbeat and returns the name to the pool.
func (s *NameService) KeepAlive(name string) {
	path := fmt.Sprintf("%s/%s", s.root, name)
	for {
		data, stat, err := s.zkConn.Get(path)
		if err != nil {
			logger.Errorf("keepAlive get path error error:%v", err)
			continue
		}
		var info = new(NameInfo)
		if err := json.Unmarshal(data, &info); err != nil {
			// Log the error if failed to unmarshal the data, and continue to the next iteration
			logger.Errorf("keepAlive unmarshal error:%v", err)
			continue
		}
		if info.Name != name {
			logger.Errorf("keepAlive get errorName:%v", info.Name)
			continue
		}
		// Update the node data and set a new expiration time
		info.ExpiredAt = info.ExpiredAt.Add(heartbeatInterval * 3)
		data, err = json.Marshal(info)
		if err != nil {
			continue
		}
		_, err = s.zkConn.Set(path, data, stat.Version)
		if err != nil {
			if err == zk.ErrBadVersion {
				// If the version is not consistent, another client has modified the data, need to get the node data again
				continue
			}
			logger.Errorf("keepalive set error:%v", err)
		}
		select {
		case <-time.After(heartbeatInterval):
		case <-s.ctx.Done():
			// If the service is stopped, stop the heartbeat and return the name
			s.ReturnName(name)
			s.zkConn.Close()
			return
		}
	}
}
