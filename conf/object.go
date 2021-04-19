package conf

/**************************************************************\
*                        PgDatabase                            *
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

/**************************************************************\
*                         PgUser                               *
\**************************************************************/
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
	Pgbouncer   bool              `yaml:"pgbouncer"`   // optional, add this database to pgbouncer list? true by default
	ConnLimit   int               `yaml:"connlimit"`   // connection limit, -1 disable limit
	ExpireIn    int               `yaml:"expire_in"`   // now + n days when this role is expired (OVERWRITE expire_at)
	ExpireAt    string            `yaml:"expire_at"`   // 'timestamp' when this role is expired
	Comment     string            `yaml:"comment"`     // optional, comment string for database
	Roles       []string          `yaml:"roles,flow"`  // dborole_admin|dbrole_readwrite|dbrole_readonly
	Parameters  map[string]string `yaml:"parameters,flow"`
}

/**************************************************************\
*                         PgHba                                *
\**************************************************************/
// PgHba hold one single hba rules
type PgHba struct {
	Title string   `yaml:"title" json:"title"`
	Role  string   `yaml:"role" json:"role"`
	Rules []string `yaml:"rules" json:"rules"`
}

/**************************************************************\
*                       PgService                              *
\**************************************************************/
// PgService hold one service definition
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
