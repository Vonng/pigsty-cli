package conf

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

/**************************************************************\
*                           Vars                               *
\**************************************************************/
// Vars holds ordered config entries
type Vars struct {
	Keys []string               // Keys hold config entry's order
	Data map[string]interface{} // Data hold actual entries
}

/**************************************************************\
*                         Access                               *
\**************************************************************/

// Has check whether key exists
func (v Vars) Has(key string) bool {
	_, exists := v.Data[key]
	return exists
}

// Put will append new entry to vars
func (v Vars) Put(key string, value interface{}) {
	if !v.Has(key) {
		v.Keys = append(v.Keys, key)
	}
	v.Data[key] = value
}

// Get value by key from Vars
func (v Vars) Get(key string) interface{} {
	return v.Data[key]
}

// GetString value and asset it as string
func (v Vars) GetString(key string) (string, bool) {
	if i, exists := v.Data[key]; exists {
		s, ok := i.(string)
		return s, ok
	}
	return "", false
}

// GetInteger value and asset it as integer
func (v Vars) GetInteger(key string) (int, bool) {
	if i, ok := v.Data[key]; ok {
		res, ok := i.(int)
		return res, ok
	}
	return 0, false
}

// GetBool value and asset it as bool
func (v Vars) GetBool(key string) (bool, bool) {
	if i, ok := v.Data[key]; ok {
		res, ok := i.(bool)
		return res, ok
	}
	return false, false
}

// GetArray value and asset it as array
func (v Vars) GetArray(key string) ([]interface{}, bool) {
	if i, ok := v.Data[key]; ok {
		res, ok := i.([]interface{})
		return res, ok
	}
	return nil, false
}

// GetMap value and asset it as array
func (v Vars) GetMap(key string) (map[string]interface{}, bool) {
	if i, ok := v.Data[key]; ok {
		res, ok := i.(map[string]interface{})
		return res, ok
	}
	return nil, false
}

/**************************************************************\
*                        Serialize                             *
\**************************************************************/
// String will print config entries in origin order
func (v *Vars) String() string {
	b, _ := yaml.Marshal(v)
	return string(b)
}

// MarshalYAML will turn Vars into yaml dict according to v.Keys order
func (v Vars) MarshalYAML() (interface{}, error) {
	var n yaml.Node
	if err := n.Encode(v.Data); err != nil {
		return nil, err
	}
	if v.Keys != nil {
		sortMapNode(&n, v.Keys)
	}
	return n, nil
}

// MarshalJSON will turn vars into json
func (v Vars) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Data)
}

/**************************************************************\
*                       Deserialize                            *
\**************************************************************/
// UnmarshalYAML will parse yaml.Node into Vars structure and preserve order
func (v *Vars) UnmarshalYAML(node *yaml.Node) (err error) {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("vars must contain YAML mapping, has %v", node.Kind)
	}
	var keys []string
	var data map[string]interface{}
	if err = node.Decode(&data); err != nil {
		return err
	}
	for i := 0; i < len(node.Content); i += 2 {
		keys = append(keys, node.Content[i].Value)
	} // decode key in order
	*v = Vars{keys, data}
	return nil
}

// ParseDatabases will parse pg_databases field into structure
func (v *Vars) ParseDatabases() (dbs []PgDatabase) {
	i, exists := v.Data["pg_databases"]
	if !exists {
		return
	}
	b, err := yaml.Marshal(i)
	if err != nil {
		return
	}
	if err = yaml.Unmarshal(b, &dbs); err != nil {
		return
	}
	return
}

// ParseUsers will parse pg_databases field into structure
func (v *Vars) ParseUsers() (dbs []PgUser) {
	i, exists := v.Data["pg_users"]
	if !exists {
		return
	}
	b, err := yaml.Marshal(i)
	if err != nil {
		return
	}
	if err = yaml.Unmarshal(b, &dbs); err != nil {
		return
	}
	return
}
