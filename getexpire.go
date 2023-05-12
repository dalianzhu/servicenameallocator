package servicenameallocator

import (
	"encoding/json"
	"fmt"
	"time"
)

// GetNamesInfo retrieves the status of all nodes under the root path.
// freeNames are the names that are currently available.
// expiredNames are the names that are currently timed out,
// possibly due to the service holding the name disconnecting or crashing.
func (s *NameService) GetNamesInfo() (
	freeNames map[string]NameInfo, expiredNames map[string]NameInfo, err error) {
	children, _, err := s.zkConn.Children(s.root)
	if err != nil {
		return nil, nil, err
	}
	expiredNames = make(map[string]NameInfo)
	freeNames = make(map[string]NameInfo)
	for _, name := range children {
		path := fmt.Sprintf("%s/%s", s.root, name)
		data, _, err := s.zkConn.Get(path)
		if err != nil {
			// The node does not exist or fails to retrieve node data.
			continue
		}
		var info NameInfo
		if err := json.Unmarshal(data, &info); err != nil {
			// Failed to parse node data.
			continue
		}
		if info.Addr == "" {
			freeNames[name] = info
			continue
		}
		if info.ExpiredAt.Before(time.Now()) {
			expiredNames[name] = info
		}
	}
	return
}
