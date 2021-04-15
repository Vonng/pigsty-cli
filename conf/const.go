package conf

/**************************************************************\
*                          Const                               *
\**************************************************************/
// meta
const GROUP_META = "meta"

// roles
const (
	ROLE_PRIMARY = "primary"
	ROLE_REPLICA = "replica"
	ROLE_STANDBY = "standby"
	ROLE_OFFLINE = "offline"
	ROLE_DELAYED = "delayed"
)

// AvailableRoles contains valid role names
var AvailableRoles = map[string]bool{
	ROLE_PRIMARY: true,
	ROLE_REPLICA: false,
	ROLE_STANDBY: false,
	ROLE_OFFLINE: false,
	ROLE_DELAYED: false,
}

// name type
var (
	NameInvalid  = "null"
	NameCluster  = "cls"
	NameInstance = "ins"
	NameIP       = "ip"
)
