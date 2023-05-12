package servicenameallocator

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

// GetName retrieves a name from the pool. If the pool is empty or there is a network error, an error is returned.
func (s *NameService) GetName() (string, error) {
	children, _, err := s.zkConn.Children(s.root)
	if err != nil {
		return "", err
	}
	if len(children) == 0 {
		return "", fmt.Errorf("name pool is empty")
	}
	rand.Shuffle(len(children), func(i, j int) { children[i], children[j] = children[j], children[i] })
	for _, name := range children {
		path := fmt.Sprintf("%s/%s", s.root, name)
		data, stat, err := s.zkConn.Get(path)
		if err != nil {
			// Node does not exist or failed to get node data
			logger.Errorf("getName get path error:%v", err)
			continue
		}
		var info NameInfo
		if err := json.Unmarshal(data, &info); err != nil {
			// Failed to parse node data
			logger.Errorf("getName unmarshal error:%v", err)
			continue
		}
		if info.Addr != "" {
			// Node is already in use, skip
			continue
		}
		info.Name = name
		info.Addr = s.addr
		info.ExpiredAt = info.ExpiredAt.Add(heartbeatInterval * 3)
		data, _ = json.Marshal(info)
		_, err = s.zkConn.Set(path, []byte(data), stat.Version)
		if err != nil {
			// Failed to update node data
			logger.Errorf("getName set path error:%v", err)
			continue
		}
		return name, nil
	}
	return "", fmt.Errorf("cannot find a free name")
}

// ReturnName returns a name to the pool by resetting its information.
func (ns *NameService) ReturnName(name string) error {
	nodePath := fmt.Sprintf("%s/%s", ns.root, name)

	// Get the current data of the name node
	data, stat, err := ns.zkConn.Get(nodePath)
	if err != nil {
		// Failed to get the data of the name node
		return err
	}

	// Update the data of the name node
	var info NameInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return err
	}

	// Clean up all information to indicate that the node has been returned
	info.Name = ""
	info.Addr = ""
	data, _ = json.Marshal(info)
	_, err = ns.zkConn.Set(nodePath, []byte(data), stat.Version)
	if err != nil {
		return err
	}
	return err
}
