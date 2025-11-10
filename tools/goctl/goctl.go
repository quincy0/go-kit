package main

import (
	"go-kit/core/load"
	"go-kit/core/logx"
	"go-kit/tools/goctl/cmd"
)

func main() {
	logx.Disable()
	load.Disable()
	cmd.Execute()
}
