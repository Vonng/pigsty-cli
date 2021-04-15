package conf

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

/**************************************************************\
*                           Config                             *
\**************************************************************/
// Config hold parsed pigsty configuration files
type Config struct {
	Clusters []Cluster `yaml:"children" json:"children"`
	Vars     Vars      `yaml:"vars" json:"vars"`

	// parsed fields
	MetaCluster *Cluster             `yaml:"-" json:"-"`
	ClusterMap  map[string]*Cluster  `yaml:"-" json:"-"`
	InstanceMap map[string]*Instance `yaml:"-" json:"-"`
	IpMap       map[string]*Instance `yaml:"-" json:"-"`
	path        string               `yaml:"-"`
	raw         []byte               `yaml:"-"`
}

// GetCluster will return cluster according to name
func (c *Config) GetCluster(name string) *Cluster {
	if name == GROUP_META {
		return c.MetaCluster
	} else {
		return c.ClusterMap[name]
	}
}

// GetInstance will return instance according to name or IP
func (c *Config) GetInstance(name string) *Instance {
	if c.IsMetaNode(name) {
		return c.MetaCluster.GetInstance(name)
	}
	if ins, exists := c.InstanceMap[name]; exists {
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

// IsMetaNode check whether given name is a meta node name or ip address
func (c *Config) IsMetaNode(name string) bool {
	for _, ins := range c.MetaCluster.Instances {
		if name == ins.Name || name == ins.IP {
			return true
		}
	}
	return false
}

// NameType tells type of a given name: cluster|instance|ip|invalid
func (c *Config) NameType(n string) string {
	// ip address is not likely to be used as instance name and cluster name
	if IsValidIP(n) {
		if _, ipFound := c.IpMap[n]; ipFound {
			return NameIP
		} else {
			return NameInvalid
		}
	}
	if _, insFound := c.InstanceMap[n]; insFound {
		return NameInstance
	}
	if _, clsFound := c.ClusterMap[n]; clsFound {
		return NameCluster
	}
	return NameInvalid
}

// GetInstancesByName will translate name into instance list
func (c *Config) GetInstancesByName(name string) []*Instance {
	switch c.NameType(name) {
	case NameInvalid:
		return []*Instance{}
	case NameIP:
		return []*Instance{c.IpMap[name]}
	case NameInstance:
		return []*Instance{c.InstanceMap[name]}
	case NameCluster:
		var res []*Instance
		for i, _ := range c.ClusterMap[name].Instances {
			res = append(res, &(c.ClusterMap[name].Instances[i]))
		}
		return res
	default:
		return []*Instance{}
	}
}

// MarshalYAML will parse yaml.Node into Vars structure and preserve order
func (c *Config) MarshalYAML() (interface{}, error) {
	var tmp = struct {
		Children map[string]*Cluster `yaml:"children"`
		Vars     Vars                `yaml:"vars"`
	}{c.ClusterMap, c.Vars}
	all := map[string]interface{}{
		"all": tmp,
	}
	return all, nil
}

// UnmarshalYAML will parse yaml.Node into Vars structure and preserve order
func (c *Config) UnmarshalYAML(v *yaml.Node) (err error) {
	var raw struct { // top layer of configuration file
		All struct {
			Clusters yaml.Node `yaml:"children"`
			Vars     Vars      `yaml:"vars"`
		} `yaml:"all"`
	}
	var cls struct {
		Hosts yaml.Node `yaml:"hosts"`
		Vars  Vars      `yaml:"vars"`
	}
	if err = v.Decode(&raw); err != nil {
		return
	}

	// parse pg cluster nodes
	clsNodes := raw.All.Clusters.Content
	clsCount := len(clsNodes) / 2
	clusters := make([]Cluster, clsCount)
	for i := 0; i < clsCount; i += 1 {
		clsname, clsnode := clsNodes[2*i].Value, clsNodes[2*i+1]
		if err = clsnode.Decode(&cls); err != nil {
			return
		}
		cluster := NewCluster(clsname, cls.Vars)

		// parse instances of cluster in order
		insnodes := cls.Hosts.Content
		inscount := len(insnodes) / 2
		for j := 0; j < inscount; j += 1 {
			var insvars Vars
			insip := insnodes[2*j].Value
			if err = insnodes[2*j+1].Decode(&insvars); err != nil {
				return
			}
			if err = cluster.AddInstance(insip, insvars); err != nil {
				return err
			}
		}
		clusters[i] = *cluster
	}
	*c = Config{
		Clusters: clusters,
		Vars:     raw.All.Vars,
	}
	c.BuildIndex()
	return nil
}

// BuildIndex will fill auxiliary fields in config struct
func (c *Config) BuildIndex() {
	clsMap := make(map[string]*Cluster)
	insMap := make(map[string]*Instance)
	ipMap := make(map[string]*Instance)

	// if meta node occurs on other pgsql group, it's vars will be overwritten
	for i, cls := range c.Clusters {
		if cls.Name == GROUP_META {
			c.MetaCluster = &(c.Clusters[i])
			continue
		}
		clsMap[cls.Name] = &(c.Clusters[i])
		for j, ins := range cls.Instances {
			insMap[ins.Name] = &(cls.Instances[j])
			ipMap[ins.IP] = &(cls.Instances[j])
		}
	}

	c.ClusterMap = clsMap
	c.InstanceMap = insMap
	c.IpMap = ipMap
}

// ParseConfig will unmarshal data into config
func ParseConfig(data []byte) (cfg *Config, err error) {
	err = yaml.Unmarshal(data, &cfg)
	return
}

// LoadConfig will read config file from disk
func LoadConfig(path string) (cfg *Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if cfg, err = ParseConfig(data); err != nil {
		return
	}
	cfg.path = path
	cfg.raw = data
	return cfg, nil
}

// InfraInfo print digest about infrastructure
func (cfg *Config) InfraInfo() string {
	var buf bytes.Buffer
	primaryIP := cfg.MetaCluster.Instances[0].IP

	// write meta node info
	buf.WriteString(fmt.Sprintf("Meta (%d): \n", len(cfg.MetaCluster.Instances)))
	for _, ins := range cfg.MetaCluster.Instances {
		if ins.IP == primaryIP {
			buf.WriteString(fmt.Sprintf("    - %s [primary]\n", ins.IP))
		} else {
			buf.WriteString(fmt.Sprintf("    - %s\n", ins.IP))
		}
	}

	// write dcs info
	dcsType, _ := cfg.Vars.GetString("dcs_type")
	dcsMap, _ := cfg.Vars.GetMap("dcs_servers")
	buf.WriteString(fmt.Sprintf("\nDCS (%s):\n", dcsType))
	for k, v := range dcsMap {
		buf.WriteString(fmt.Sprintf("    %s: %s\n", k, v))
	}

	// write nginx upstream info
	buf.WriteString("\nNginx: \n")
	routes, _ := cfg.Vars.GetArray("nginx_upstream")
	for _, r := range routes {
		entry := r.(map[string]interface{})
		buf.WriteString(fmt.Sprintf("    - %-12s (%s)\thttp://%s\t ->  http://%-16s\n",
			entry["name"], entry["url"],
			entry["host"],
			strings.Replace(entry["url"].(string), "127.0.0.1", primaryIP, -1)))
	}

	// write repo info
	buf.WriteString("\nRepo: \n")
	repoAddress, _ := cfg.Vars.GetString("repo_address")
	repoName, _ := cfg.Vars.GetString("repo_name")
	repoHome, _ := cfg.Vars.GetString("repo_home")
	buf.WriteString(fmt.Sprintf("    - http://%s -> %s:%s/%s\n", repoAddress, primaryIP, repoHome, repoName))

	// write NTP info
	if ntpConfig, _ := cfg.Vars.GetBool("node_ntp_config"); ntpConfig {
		buf.WriteString("\nNTP: \n")
		ntpServers, _ := cfg.Vars.GetArray("node_ntp_servers")
		for _, s := range ntpServers {
			buf.WriteString(fmt.Sprintf("    - %s\n", s))
		}
	}

	// write DNS info
	dnsEnabled, _ := cfg.Vars.GetString("node_dns_server")
	if dnsEnabled == "add" || dnsEnabled == "overwrite" {
		buf.WriteString("\nDNS: \n")
		dnsServers, _ := cfg.Vars.GetArray("node_dns_servers")
		for _, s := range dnsServers {
			buf.WriteString(fmt.Sprintf("    - %s\n", s))
		}
	}
	return buf.String()
}