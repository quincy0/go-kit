package main

import (
	"github.com/quincy0/go-kit/core/load"
	"github.com/quincy0/go-kit/core/logx"
	"github.com/quincy0/go-kit/tools/goctl/cmd"
)

func main() {
	logx.Disable()
	load.Disable()
	cmd.Execute()
}
