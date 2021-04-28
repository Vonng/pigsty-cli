package server

import "testing"

func TestNewPigstyServer(t *testing.T) {
	InitDefaultServer(`/Users/vonng/pigsty/pigsty.yml`, "/Users/vonng/pigsty/public", ":9633")
	PS.Run()
}
