package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
)

/**************************************************************\
*                        Instance                              *
\**************************************************************/
type Instance struct {
	IP      string   `yaml:"ip" json:"ip"`
	Name    string   `yaml:"pg_instance" json:"name"`
	Seq     int      `yaml:"pg_seq"  json:"seq"`
	Role    string   `yaml:"pg_role"  json:"role"`
	Vars    Vars     `yaml:"-,flow"  json:"vars"`
	Cluster *Cluster `yaml:"-" json:"-"`
}

// MarshalYAML will turn Vars into yaml dict according to v.Keys order
func (i Instance) MarshalYAML() (interface{}, error) {
	return i.Vars, nil
}

// String will turn instance into string representation
func (i Instance) String() string {
	b, _ := yaml.Marshal(i)
	return string(b)
}

// Summary will print instance digest
func (i *Instance) Summary() string {
	return fmt.Sprintf(`[%d](%s) %-15s %s`, i.Seq, i.Role, i.IP, i.Name)
}

// IsValid tells whether an instance content is valid
func (i *Instance) IsValid() bool {
	if i.IP == "" || i.Name == "" || i.Role == "" {
		return false // empty fields
	}
	matched, err := regexp.MatchString(`[a-zA-Z][a-zA-Z0-9_-]*-\d+`, i.Name)
	if err != nil || !matched {
		return false // invalid instance name
	}

	if _, exists := AvailableRoles[i.Role]; !exists {
		return false // invalid role
	}

	// TODO: validate ip
	return true
}

// MatchName test whether a name or glob regexp matches cluster's name
func (i *Instance) MatchName(name string) bool {
	if i.Name == name || i.Cluster.Name == name {
		return true
	}
	if IsValidIP(name) && i.IP == name {
		return true
	}

	re, err := regexp.CompilePOSIX(name)
	if err != nil {
		return false
	}
	if re.MatchString(i.Name) || re.MatchString(i.Cluster.Name) {
		return true
	}
	return false
}

func (i *Instance) MatchNames(names []string) bool {
	for _, n := range names {
		if i.MatchName(n) {
			return true
		}
	}
	return false
}
