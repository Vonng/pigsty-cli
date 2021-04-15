package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestParseConfig(t *testing.T) {
	testCase := `
all:
  children:
    meta:
      vars:
        meta_node: true                     # mark node as meta controller
        ansible_group_priority: 99          # meta group is top priority
      hosts: {10.10.10.10: {ansible_host: meta}}

    pg-meta:
      hosts:
        10.10.10.10: {pg_seq: 1, pg_role: primary}
      vars:
        pg_cluster: pg-meta                 # define actual cluster name
        pg_version: 13                      # define installed pgsql version
        node_tune: tiny                     # tune node into oltp|olap|crit|tiny mode
`

	cfg, _ := ParseConfig([]byte(testCase))
	fmt.Println(cfg)
}

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig(`/Users/vonng/pigsty/pigsty.yml`)
	if err != nil {
		panic(err)
	}

	if b, err := yaml.Marshal(cfg); err != nil {
		panic(err)
	} else {
		fmt.Println(string(b))
	}
}
