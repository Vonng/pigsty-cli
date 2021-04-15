package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"regexp"
	"strings"
)

/**************************************************************\
*                          Const                               *
\**************************************************************/
const (
	ROLE_PRIMARY = "primary"
	ROLE_REPLICA = "replica"
	ROLE_STANDBY = "standby"
	ROLE_OFFLINE = "offline"
	ROLE_DELAYED = "delayed"

	GROUP_META = "meta"
)

var AvailableRoles = map[string]bool{
	ROLE_PRIMARY: true,
	ROLE_REPLICA: false,
	ROLE_STANDBY: false,
	ROLE_OFFLINE: false,
	ROLE_DELAYED: false,
}

var (
	NameInvalid  = "null"
	NameCluster  = "cls"
	NameInstance = "ins"
	NameIP       = "ip"
)

/**************************************************************\
*                         Cluster                              *
\**************************************************************/
// Cluster define a postgres cluster
type Cluster struct {
	// Identity Parameters
	Name   string `json:"pg_cluster"`          // cluster name
	Shard  string `json:"pg_shard,omitempty"`  // shard name (optional)
	SIndex int    `json:"pg_sindex,omitempty"` // shard index (optional)

	// Variable
	Vars Vars `json:"vars"` // original vars

	// Parsed Info
	Instances   []Instance           `json:"hosts"`
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
	if c.Name != GROUP_META {
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

/**************************************************************\
*                        PgDB/User                             *
\**************************************************************/
type PgDatabase struct {
	Name       string `yaml:"name"`       // name is the only required field for a database
	Owner      string `yaml:"owner"`      // optional, database owner
	Template   string `yaml:"template"`   // optional, template1 by default
	Encoding   string `yaml:"encoding"`   // optional, UTF8 by default , must same as template database, leave blank to set to db default
	Locale     string `yaml:"locale"`     // optional, C by default , must same as template database, leave blank to set to db default
	LcCollate  string `yaml:"lc_collate"` // optional, C by default , must same as template database, leave blank to set to db default
	LcCtype    string `yaml:"lc_ctype"`   // optional, C by default , must same as template database, leave blank to set to db default
	AllowConn  bool   `yaml:"allowconn"`  // optional, true by default, false disable connect at all
	RevokeConn bool   `yaml:"revokeconn"` // optional, false by default, true revoke connect from public # (only default user and owner have connect privilege on database)
	Tablespace string `yaml:"tablespace"` // optional, 'pg_default' is the default tablespace
	ConnLimit  int    `yaml:"connlimit"`  // optional, connection limit, -1 or none disable limit (default)
	Pgbouncer  bool   `yaml:"pgbouncer"`  // optional, add this database to pgbouncer list? true by default
	Comment    string `yaml:"comment"`    // optional, comment string for database
	Extensions []struct {
		Name   string `yaml:"name"`
		Schema string `yaml:"schema"`
	} `yaml:"extensions,flow"`
	Parameters map[string]string `yaml:"parameters,flow"`
}

type PgUser struct {
	Name        string            `yaml:"name"`        // user name
	Password    string            `yaml:"password"`    // password, can be md5 encrypted
	Login       bool              `yaml:"login"`       // can login, true by default (should be false for role)
	Superuser   bool              `yaml:"superuser"`   // is superuser? false by default
	CreateDB    bool              `yaml:"createdb"`    // can create database? false by default
	CreateRole  bool              `yaml:"createrole"`  // can create role? false by default
	Inherit     bool              `yaml:"inherit"`     // can this role use inherited privileges?
	Replication bool              `yaml:"replication"` // can this role do replication? false by default
	BypassRLS   bool              `yaml:"bypassrls"`   // can this role bypass row level security? false by default
	ConnLimit   int               `yaml:"connlimit"`   // connection limit, -1 disable limit
	ExpireAt    string            `yaml:"expire_at"`   // 'timestamp' when this role is expired
	ExpireIn    int               `yaml:"expire_in"`   // now + n days when this role is expired (OVERWRITE expire_at)
	Roles       []string          `yaml:"roles,flow"`  // dborole_admin|dbrole_readwrite|dbrole_readonly
	Pgbouncer   bool              `yaml:"pgbouncer"`   // optional, add this database to pgbouncer list? true by default
	Comment     string            `yaml:"comment"`     // optional, comment string for database
	Parameters  map[string]string `yaml:"parameters,flow"`
}

type PgHba struct {
	Title string   `yaml:"title"`
	Role  string   `yaml:"role"`
	Rules []string `yaml:"rules"`
}

type PgService struct {
	Name           string            `yaml:"name"`
	SrcIP          string            `yaml:"src_ip"`
	SrcPort        int               `yaml:"src_port"`
	DstPort        int               `yaml:"dst_port"`
	CheckURL       string            `yaml:"check_url"`
	Selector       string            `yaml:"selector"`
	SelectorBackup string            `yaml:"selector_backup"`
	HAProxy        map[string]string `yaml:"haproxy"`
}

/**************************************************************\
*                           Vars                               *
\**************************************************************/
// Vars holds ordered config entries
type Vars struct {
	Keys []string
	Data map[string]interface{}
}

// Put will add new entry to vars
func (v Vars) Put(key string, value interface{}) {
	v.Keys = append(v.Keys, key)
	v.Data[key] = value
}

// Has check whether key exists
func (v Vars) Has(key string) bool {
	_, exists := v.Data[key]
	return exists
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

func (v Vars) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Data)
}

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

/**************************************************************\
*                           Config                             *
\**************************************************************/
type Config struct {
	Clusters []Cluster `yaml:"children" json:"children"`
	Vars     Vars      `yaml:"vars" json:"vars"`

	MetaCluster *Cluster             `yaml:"-" json:"-"`
	ClusterMap  map[string]*Cluster  `yaml:"-" json:"-"`
	InstanceMap map[string]*Instance `yaml:"-" json:"-"`
	IpMap       map[string]*Instance `yaml:"-" json:"-"`

	path string `yaml:"-"`
	raw  []byte `yaml:"-"`
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

// Translate limit selector into instance list
func Translate(limit string) ([]*Instance, error) {
	if len(limit) == 0 {
		return []*Instance{}, nil
	}
	//entries := strings.Split(limit, ",")
	return nil, nil
}

// ToInventory will generate ansible inventory file from config structure
func ToInventory() string {
	return ""
}

// MarshalYAML will parse yaml.Node into Vars structure and preserve order
func (c *Config) MarshalYAML() (interface{}, error) {
	fmt.Println("====================")
	var tmp = struct {
		Children map[string]*Cluster `yaml:"children"`
		Vars     Vars                `yaml:"vars"`
	}{c.ClusterMap, c.Vars}
	all := map[string]interface{}{
		"all": tmp,
	}

	fmt.Println("====================")
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

// IsValidIP tells if a string is a valid ip address
func IsValidIP(s string) bool {
	if address := net.ParseIP(s); address != nil {
		return true
	} else {
		return false
	}
}
