package resolver

import (
	"github.com/quincy0/go-kit/zrpc/resolver/internal"
)

// Register registers schemes defined zrpc.
// Keep it in a separated package to let third party register manually.
func Register() {
	internal.RegisterResolver()
}
