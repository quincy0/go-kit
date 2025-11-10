package conf

import (
	"fmt"
	"os"
	"testing"
)

type Config struct {
	Port uint64
	Host string
}

func TestConfigLoad(t *testing.T) {
	os.Setenv("ENVIRONMENT", "dev")

	remoteConf, err := NewRemoteConfig("cloud-run", "hubbuy-scm-api")
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	c := Config{}

	// remoteConf.WithConfig("cloud-run", "app-test2")

	if err := remoteConf.Load(&c); err != nil {
		fmt.Println("err:", err)
	}

	select {}
}
