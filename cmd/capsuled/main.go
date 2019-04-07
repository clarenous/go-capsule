package main

import (
	"github.com/clarenous/go-capsule/cmd/capsuled/cmd"
	_ "github.com/clarenous/go-capsule/consensus/algorithm/pow"
)

func main() {
	cmd.Execute()
}
