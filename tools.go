package servicenameallocator

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

// ResetNameForce force reset some nodes to initial state
func (ns *NameService) ResetNameForce(names []string) error {
	for _, name := range names {
		path := fmt.Sprintf("%s/%s", ns.root, name)

		exist, stat, err := ns.zkConn.Exists(path)
		if err != nil {
			return err
		}
		if exist {
			info := new(NameInfo)
			info.Name = ""
			info.Addr = ""
			info.ExpiredAt = time.Now()
			data, _ := json.Marshal(info)

			logger.Infof("will reset node:%v", path)
			_, err := ns.zkConn.Set(path, []byte(data), stat.Version)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// InitNameNodes creates the data for names under the root, if it doesn't exist.
// The existing data for name won't be affected.
func (ns *NameService) InitNameNodes(names []string) error {
	// Create parent node.
	_, err := ns.createRecursive(ns.root, []byte{}, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	for _, name := range names {
		path := fmt.Sprintf("%s/%s", ns.root, name)

		exist, _, err := ns.zkConn.Exists(path)
		if err != nil {
			return err
		}
		if !exist {
			info := new(NameInfo)
			info.Name = ""
			info.Addr = ""
			info.ExpiredAt = time.Now()
			data, _ := json.Marshal(info)
			logger.Infof("will create node:%v", path)
			_, err := ns.zkConn.Create(path, []byte(data), 0, zk.WorldACL(zk.PermAll))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ns *NameService) createRecursive(ppath string, data []byte, flags int32, acl []zk.ACL) (string, error) {
	// If the path already exists, return it directly
	exists, _, err := ns.zkConn.Exists(ppath)
	if err != nil {
		return "", fmt.Errorf("check path exists error: %v", err)
	}
	if exists {
		return ppath, nil
	}

	// If the parent path does not exist, recursively create the parent node
	parentPath := path.Dir(ppath)
	_, err = ns.createRecursive(parentPath, data, flags, acl)
	if err != nil {
		return "", err
	}

	// If the path does not exist, create it
	logger.Infof("will create path:%v", ppath)
	return ns.zkConn.Create(ppath, data, flags, acl)
}
