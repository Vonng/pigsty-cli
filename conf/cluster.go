package conf

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"
)

/**************************************************************\
*                         Cluster                              *
\**************************************************************/
// Cluster define a postgres cluster
type Cluster struct {
	// Identity
	Name   string `json:"pg_cluster"`          // cluster name
	Shard  string `json:"pg_shard,omitempty"`  // shard name (optional)
	SIndex int    `json:"pg_sindex,omitempty"` // shard index (optional)

	// Payload
	Instances []Instance `json:"hosts"` // original instances
	Vars      Vars       `json:"vars"`  // original vars

	// Parsed Fields
	PgUsers     []PgUser             `json:"-"`
	PgDatabases []PgDatabase         `json:"-"`
	NameMap     map[string]*Instance `json:"-"`
	SeqMap      map[int]*Instance    `json:"-"`
	IpMap       map[string]*Instance `json:"-"`
	Primary     *Instance            `json:"-"`
}

// NewCluster will create new cluster structure from name and vars
func NewCluster(name string, vars Vars) *Cluster {
	cluster := &Cluster{Name: name, Vars: vars}
	cluster.Shard, _ = vars.GetString("pg_shard")
	cluster.SIndex, _ = vars.GetInteger("pg_sindex")
	cluster.PgUsers = vars.ParseUsers()
	cluster.PgDatabases = vars.ParseDatabases()
	cluster.NameMap = make(map[string]*Instance)
	cluster.SeqMap = make(map[int]*Instance)
	cluster.IpMap = make(map[string]*Instance)
	return cluster
}

// AddInstance will add new instance into cluster according to ip(key) and vars(value)
func (c *Cluster) AddInstance(ip string, vars Vars) error {
	if ip == "" {
		return fmt.Errorf("invalid instance ip: %v", vars)
	}
	if !IsValidIP(ip) {
		return fmt.Errorf("invalid instance ip: %s", ip)
	}
	seq, seqExists := vars.GetInteger("pg_seq")
	role, roleExists := vars.GetString("pg_role")
	if c.Name != GROUP_META { // meta group does not require identity fields
		if !seqExists {
			return fmt.Errorf("instance pg_seq is required: %v", vars)
		}
		if !roleExists {
			return fmt.Errorf("instance pg_role is required: %s", ip)
		}
		if _, exists := AvailableRoles[role]; !exists {
			return fmt.Errorf("invalid pg_role value %s for %s", role, ip)
		}
	}
	c.Instances = append(c.Instances, Instance{
		IP:      ip,
		Name:    fmt.Sprintf("%s-%d", c.Name, seq),
		Seq:     seq,
		Role:    role,
		Vars:    vars,
		Cluster: c,
	})
	var ins *Instance
	ins = &(c.Instances[len(c.Instances)-1])
	if ins.Role == ROLE_PRIMARY {
		c.Primary = ins
	}
	c.IpMap[ins.IP] = ins
	c.SeqMap[ins.Seq] = ins
	c.NameMap[ins.Name] = ins
	return nil
}

// GetInstance will return cluster's instance according to instance name or ip
func (c *Cluster) GetInstance(name string) *Instance {
	if ins, exists := c.NameMap[name]; exists {
		return ins
	}
	if !IsValidIP(name) {
		return nil
	}
	if ins, exists := c.IpMap[name]; exists {
		return ins
	}
	return nil
}

// IPList return ordered instances IP list
func (c *Cluster) IPList() (res []string) {
	for _, ins := range c.Instances {
		res = append(res, ins.IP)
	}
	return
}

// NameList return ordered instances ins-name list
func (c *Cluster) NameList() (res []string) {
	for _, ins := range c.Instances {
		res = append(res, ins.Name)
	}
	return
}

// MarshalYAML will turn Vars into yaml dict according to v.Keys order
func (c Cluster) MarshalYAML() (interface{}, error) {
	insMap := make(map[string]Vars)
	for _, ins := range c.Instances {
		insMap[ins.IP] = ins.Vars
	}
	var hosts yaml.Node
	if err := hosts.Encode(insMap); err != nil {
		return nil, err
	}
	sortMapNode(&hosts, c.IPList())
	return struct {
		Hosts *yaml.Node `yaml:"hosts"`
		Vars  Vars       `yaml:"vars"`
	}{
		&hosts,
		c.Vars,
	}, nil
}

// String will print one-line representation of cluster
func (c *Cluster) String() string {
	var insSum []string
	for _, ins := range c.Instances {
		insSum = append(insSum, fmt.Sprintf(`%d-%s: %-15s`, ins.Seq, ins.Role, ins.IP))
	}
	return fmt.Sprintf("%-32s %s", c.Name, strings.Join(insSum, " "))
}

// JSON return json repr of this cluster
func (c *Cluster) JSON() string {
	b, err := json.MarshalIndent(*c, "", "    ")
	if err != nil {
		return "{}"
	}
	return string(b)
}

// YAML returns yaml repr of this cluster
func (c *Cluster) YAML() string {
	tmp := make(map[string]Cluster, 1)
	tmp[c.Name] = *c
	b, _ := yaml.Marshal(tmp)
	return string(b)
}

// Summary will print cluster summary
func (c *Cluster) Summary() string {
	var users, dbs, instances []string
	for _, user := range c.PgUsers {
		users = append(users, user.Name)
	}
	for _, db := range c.PgDatabases {
		dbs = append(dbs, db.Name)
	}
	for _, ins := range c.Instances {
		instances = append(instances, "      - "+ins.Summary())
	}
	s := fmt.Sprintf("---------------------------------------\n- Cluster  :  %s\n  Usernames:  %s\n  Databases:  %s\n  Instances:\n%s",
		c.Name,
		strings.Join(users, ", "),
		strings.Join(dbs, ", "),
		strings.Join(instances, "\n"),
	)
	return s
}

// Repr return cluster string representation according to format
func (c *Cluster) Repr(format string) (s string) {
	switch format {
	case "yaml", "y":
		return c.YAML()
	case "json", "j":
		return c.JSON()
	case "detail", "d", "summary", "s":
		return c.Summary()
	default:
		return c.String()
	}
}

// MatchName test whether a name or glob regexp matches cluster's name
func (c *Cluster) MatchName(name string) bool {
	if c.Name == name {
		return true
	}
	re, err := regexp.CompilePOSIX(name)
	if err != nil {
		return false
	}
	return re.MatchString(c.Name)
}

// MatchNames perform multiple match
func (c *Cluster) MatchNames(names []string) bool {
	for _, n := range names {
		if c.MatchName(n) {
			return true
		}
	}
	return false
}
