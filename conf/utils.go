package conf

import (
	"gopkg.in/yaml.v3"
	"net"
)

// IsValidIP tells if a string is a valid ip address
func IsValidIP(s string) bool {
	if address := net.ParseIP(s); address != nil {
		return true
	} else {
		return false
	}
}

// getMapKeys return map node's key in original order
func getMapKeys(value yaml.Node) []string {
	if value.Kind != yaml.MappingNode {
		return nil
	}
	var key string
	var keys []string
	for i := 0; i < len(value.Content); i += 2 {
		if err := value.Content[i].Decode(&key); err != nil {
			return nil
		} else {
			keys = append(keys, key)
		}
	}
	return keys
}

// sortMapNode will adjust map entry order according to keys
func sortMapNode(mapNode *yaml.Node, keys []string) {
	if mapNode.Kind != yaml.MappingNode {
		panic("invalid node type")
	}
	if len(mapNode.Content)&1 != 0 || len(mapNode.Content) != len(keys)*2 {
		panic("invalid map or keys")
	}
	nEntry := len(keys)
	sortedNodes := make([]*yaml.Node, 2*len(keys))
	kNode := make(map[string]*yaml.Node, nEntry)
	vNode := make(map[string]*yaml.Node, nEntry)
	for i := 0; i < nEntry; i++ {
		key := mapNode.Content[2*i].Value
		kNode[key] = mapNode.Content[2*i]   // key node
		vNode[key] = mapNode.Content[2*i+1] // value node
	}
	for i, key := range keys {
		sortedNodes[i*2] = kNode[key]
		sortedNodes[i*2+1] = vNode[key]
	}
	mapNode.Content = sortedNodes
}
