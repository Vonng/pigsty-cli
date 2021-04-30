package server

import "testing"

func TestNewPigstyServer(t *testing.T) {
	InitDefaultServer(`:9633`, "/Users/vonng/pigsty/pigsty.yml", "/tmp/pd", "embed")
}
