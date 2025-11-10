package upgrade

import (
	"fmt"
	"runtime"

	"go-kit/tools/goctl/rpc/execx"
	"github.com/spf13/cobra"
)

// Upgrade gets the latest goctl by
// go install go-kit/tools/goctl@latest
func upgrade(_ *cobra.Command, _ []string) error {
	cmd := `GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go install go-kit/tools/goctl@latest`
	if runtime.GOOS == "windows" {
		cmd = `set GOPROXY=https://goproxy.cn,direct && go install go-kit/tools/goctl@latest`
	}
	info, err := execx.Run(cmd, "")
	if err != nil {
		return err
	}

	fmt.Print(info)
	return nil
}
