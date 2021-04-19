package conf

import (
	"fmt"
	"testing"
)

func TestParseDatabase(t *testing.T) {
	cfg, err := LoadConfig(`/Users/vonng/pigsty/pigsty.yml`)
	if err != nil {
		t.Fail()
	}

	users := cfg.ClusterMap["pg-meta"].Vars.Data["pg_users"].([]interface{})
	for _, u := range users {
		fmt.Printf("%v\n", u.(map[string]interface{}))
	}

}
